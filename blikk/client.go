package blikk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	DefaultBaseURL   = "https://publicapi.blikk.com"
	DefaultRateLimit = rate.Limit(4)

	defaultUserAgent    = "blikk-sdk-go"
	defaultHTTPTimeout  = 30 * time.Second
	maxRateLimitRetries = 5
	maxRetryAfter       = 30 * time.Second
	defaultRetryAfter   = time.Second
	maxErrorBodyBytes   = 64 << 10
	jsonContentType     = "application/json"
)

type Client struct {
	Users       *UsersService
	Projects    *ProjectsService
	TimeReports *TimeReportsService

	baseURL     string
	httpClient  *http.Client
	tokenSource TokenSource
	limiter     *rate.Limiter
	userAgent   string
	sleep       func(ctx context.Context, d time.Duration) error
}

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = baseURL }
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) { c.httpClient = httpClient }
}

// WithTokenSource replaces the app ID and secret exchange; the credentials passed to NewClient are then ignored.
func WithTokenSource(tokenSource TokenSource) Option {
	return func(c *Client) { c.tokenSource = tokenSource }
}

func WithUserAgent(userAgent string) Option {
	return func(c *Client) { c.userAgent = userAgent }
}

// WithRateLimit sets the client-side request rate; the API allows 4 requests per second per token. Use rate.Inf to disable.
func WithRateLimit(limit rate.Limit) Option {
	return func(c *Client) { c.limiter = rate.NewLimiter(limit, 1) }
}

func NewClient(appID, appSecret string, opts ...Option) (*Client, error) {
	c := &Client{
		baseURL:    DefaultBaseURL,
		httpClient: &http.Client{Timeout: defaultHTTPTimeout},
		limiter:    rate.NewLimiter(DefaultRateLimit, 1),
		userAgent:  defaultUserAgent,
		sleep:      sleepContext,
	}
	for _, opt := range opts {
		opt(c)
	}
	if err := validateBaseURL(c.baseURL); err != nil {
		return nil, err
	}
	if c.tokenSource == nil {
		tokenSource, err := c.newAppTokenSource(appID, appSecret)
		if err != nil {
			return nil, err
		}
		c.tokenSource = tokenSource
	}
	c.Users = &UsersService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.TimeReports = &TimeReportsService{client: c}
	return c, nil
}

func validateBaseURL(baseURL string) error {
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("blikk: invalid base URL %q", baseURL)
	}
	return nil
}

func (c *Client) newAppTokenSource(appID, appSecret string) (*appTokenSource, error) {
	if appID == "" || appSecret == "" {
		return nil, errors.New("blikk: app ID and app secret are required unless WithTokenSource is used")
	}
	tokenURL, err := url.JoinPath(c.baseURL, tokenPath)
	if err != nil {
		return nil, fmt.Errorf("blikk: invalid base URL: %w", err)
	}
	return &appTokenSource{
		appID:      appID,
		appSecret:  appSecret,
		tokenURL:   tokenURL,
		userAgent:  c.userAgent,
		httpClient: c.httpClient,
		now:        time.Now,
		sleep:      c.sleep,
	}, nil
}

func (c *Client) do(ctx context.Context, method, resourcePath string, query url.Values, body any, out any) error {
	requestURL, err := c.buildURL(resourcePath, query)
	if err != nil {
		return err
	}
	payload, err := encodeBody(body)
	if err != nil {
		return err
	}

	rateLimitRetries := 0
	tokenInvalidated := false
	for {
		token, err := c.tokenSource.Token(ctx)
		if err != nil {
			return err
		}
		statusCode, header, responseBody, err := c.send(ctx, method, requestURL, payload, token)
		if err != nil {
			return err
		}
		switch {
		case statusCode == http.StatusTooManyRequests && rateLimitRetries < maxRateLimitRetries:
			rateLimitRetries++
			if err := c.sleep(ctx, retryAfter(header, responseBody)); err != nil {
				return err
			}
		case statusCode == http.StatusUnauthorized && !tokenInvalidated:
			tokenInvalidated = true
			c.tokenSource.Invalidate(token)
		case isSuccess(statusCode):
			return decodeBody(method, resourcePath, responseBody, out)
		default:
			return newError(statusCode, responseBody)
		}
	}
}

func (c *Client) buildURL(resourcePath string, query url.Values) (string, error) {
	requestURL, err := url.JoinPath(c.baseURL, resourcePath)
	if err != nil {
		return "", fmt.Errorf("blikk: building URL for %s: %w", resourcePath, err)
	}
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}
	return requestURL, nil
}

func (c *Client) send(ctx context.Context, method, requestURL string, payload []byte, token string) (int, http.Header, []byte, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return 0, nil, nil, err
	}
	req, err := c.newRequest(ctx, method, requestURL, payload, token)
	if err != nil {
		return 0, nil, nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("blikk: %s %s: %w", method, req.URL.Path, err)
	}
	statusCode, body, err := readResponse(resp)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("blikk: reading response of %s %s: %w", method, req.URL.Path, err)
	}
	return statusCode, resp.Header, body, nil
}

func (c *Client) newRequest(ctx context.Context, method, requestURL string, payload []byte, token string) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("blikk: creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", jsonContentType)
	req.Header.Set("User-Agent", c.userAgent)
	if payload != nil {
		req.Header.Set("Content-Type", jsonContentType)
	}
	return req, nil
}

func readResponse(resp *http.Response) (int, []byte, error) {
	defer resp.Body.Close()
	reader := io.Reader(resp.Body)
	if !isSuccess(resp.StatusCode) {
		reader = io.LimitReader(resp.Body, maxErrorBodyBytes)
	}
	body, err := io.ReadAll(reader)
	return resp.StatusCode, body, err
}

func encodeBody(body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("blikk: encoding request body: %w", err)
	}
	return payload, nil
}

func decodeBody(method, resourcePath string, body []byte, out any) error {
	if out == nil || len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("blikk: decoding response of %s %s: %w", method, resourcePath, err)
	}
	return nil
}

func isSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

type rateLimitBody struct {
	WaitInSeconds float64 `json:"waitInSeconds"`
}

func retryAfter(header http.Header, body []byte) time.Duration {
	wait, ok := retryAfterFromHeader(header)
	if !ok {
		wait, ok = retryAfterFromBody(body)
	}
	if !ok {
		wait = defaultRetryAfter
	}
	return min(max(wait, 0), maxRetryAfter)
}

func retryAfterFromHeader(header http.Header) (time.Duration, bool) {
	seconds, err := strconv.ParseFloat(header.Get("Retry-After"), 64)
	if err != nil {
		return 0, false
	}
	return secondsToDuration(seconds), true
}

func retryAfterFromBody(body []byte) (time.Duration, bool) {
	var decoded rateLimitBody
	if err := json.Unmarshal(body, &decoded); err != nil || decoded.WaitInSeconds <= 0 {
		return 0, false
	}
	return secondsToDuration(decoded.WaitInSeconds), true
}

func secondsToDuration(seconds float64) time.Duration {
	return time.Duration(seconds * float64(time.Second))
}

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func resourcePath(base string, pathParams ...string) (string, error) {
	segments := []string{base}
	for _, param := range pathParams {
		if err := validatePathParam(param); err != nil {
			return "", err
		}
		segments = append(segments, url.PathEscape(param))
	}
	return strings.Join(segments, "/"), nil
}

// url.JoinPath cleans the path, so these would silently address a different resource.
func validatePathParam(param string) error {
	if param == "" || param == "." || param == ".." {
		return fmt.Errorf("blikk: invalid path parameter %q", param)
	}
	return nil
}
