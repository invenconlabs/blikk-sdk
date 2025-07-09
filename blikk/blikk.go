package blikk

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/invenconlabs/blikk-sdk/dateutils"
)

type DateOnly = dateutils.DateOnly

var (
	PreviousWeek    = dateutils.PreviousWeek
	PreviousMonth   = dateutils.PreviousMonth
	FirstDayOfMonth = dateutils.FirstDayOfMonth
	LastDayOfMonth  = dateutils.LastDayOfMonth
)

// Client is the main client for interacting with the Blikk API.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithBaseURL allows overriding the default Blikk API base URL.
// This is useful for testing or pointing to a different environment.
func WithBaseURL(u string) ClientOption {
	return func(c *Client) {
		c.baseURL = u
	}
}

// WithHTTPClient allows providing a custom http.Client.
// Useful for custom timeouts or testing with a mock client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// NewClient creates a new Blikk API client.
func NewClient(token string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL:    "https://publicapi.blikk.com/",
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// GetAccessToken retrieves a new access token using the app ID and secret
// from environment variables (BLIKK_APP_ID, BLIKK_APP_SECRET).
func GetAccessToken() (string, error) {
	appId := os.Getenv("BLIKK_APP_ID")
	appSecret := os.Getenv("BLIKK_APP_SECRET")

	encoded := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", appId, appSecret)))

	req, err := http.NewRequest("POST", "https://publicapi.blikk.com/v1/Auth/Token", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+encoded)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("failed to get access token, status %d: %s", res.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var accessTokenResponse accessTokenResponse
	err = json.Unmarshal(body, &accessTokenResponse)
	if err != nil {
		return "", err
	}

	return accessTokenResponse.AccessToken, nil
}

// List retrieves a collection of resources.
// It handles pagination automatically, fetching all pages of results.
func List[T ListItem](c *Client, options ListOptions) ([]T, error) {
	var items []T
	var itemType T

	if !itemType.validFilter(&options) {
		return nil, fmt.Errorf("invalid filter options for %T", itemType)
	}
	u, err := url.Parse(c.baseURL + itemType.path())
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	q := u.Query()
	// Build query parameters from the options struct using reflection.
	// It adds fields with a "paramName" tag to the query, skipping zero values.
	// Special handling is provided for DateOnly fields and slices/arrays.
	for i := 0; i < reflect.TypeOf(options).NumField(); i++ {
		field := reflect.TypeOf(options).Field(i)
		paramName := field.Tag.Get("paramName")
		if paramName == "" {
			continue
		}

		fieldValue := reflect.ValueOf(options).Field(i)
		if fieldValue.IsZero() {
			continue
		}

		if fieldValue.Type() == reflect.TypeOf(&DateOnly{}) {
			date := fieldValue.Interface().(*DateOnly)
			q.Set(paramName, date.Format(time.DateOnly))
			continue
		}

		if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array {
			var values []string
			for j := 0; j < fieldValue.Len(); j++ {
				values = append(values, fmt.Sprintf("%v", fieldValue.Index(j).Interface()))
			}
			q.Set(paramName, strings.Join(values, ","))
			continue
		}
		q.Set(paramName, fmt.Sprintf("%v", fieldValue.Interface()))
	}

	u.RawQuery = q.Encode()

	for {
		body, err := c.doGetRequest(u)
		if err != nil {
			return nil, err
		}

		var response ListResponse[T]
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		items = append(items, response.Items...)

		// Blikk API is 1-indexed for pages
		if response.Page >= response.TotalPages {
			break
		}

		options.Page = response.Page + 1
		q.Set("page", fmt.Sprintf("%d", options.Page))
		u.RawQuery = q.Encode()
	}

	return items, nil
}

// Get retrieves a single resource by its identifier.
func Get[T GetItem](c *Client, query string) (T, error) {
	var item T

	u, err := url.Parse(c.baseURL + item.path(query))
	if err != nil {
		return item, fmt.Errorf("invalid base URL: %w", err)
	}

	body, err := c.doGetRequest(u)
	if err != nil {
		return item, err
	}

	err = json.Unmarshal(body, &item)
	if err != nil {
		return item, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return item, nil
}

func (c *Client) doGetRequest(u *url.URL) ([]byte, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.retryRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

func (c *Client) retryRequest(req *http.Request) (*http.Response, error) {
	for {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			retryAfter := resp.Header.Get("Retry-After")
			waitDuration := time.Second // Default wait time

			if retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					waitDuration = time.Duration(seconds) * time.Second
				} else if d, err := http.ParseTime(retryAfter); err == nil {
					waitDuration = time.Until(d)
				}
			}

			if waitDuration < 0 {
				waitDuration = time.Second
			}
			time.Sleep(waitDuration)
			continue
		}
		return resp, nil
	}
}
