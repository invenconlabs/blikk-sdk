package blikk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	tokenPath            = "v1/Auth/Token"
	tokenRefreshMargin   = 60 * time.Second
	defaultTokenLifetime = time.Hour
)

type TokenSource interface {
	Token(ctx context.Context) (string, error)
	// Invalidate discards the cached token only if it is still stale, so concurrent 401s cause one refresh.
	Invalidate(stale string)
}

type staticToken string

func StaticToken(token string) TokenSource { return staticToken(token) }

func (t staticToken) Token(context.Context) (string, error) { return string(t), nil }

func (staticToken) Invalidate(string) {}

type appTokenSource struct {
	appID      string
	appSecret  string
	tokenURL   string
	userAgent  string
	httpClient *http.Client
	now        func() time.Time
	sleep      func(ctx context.Context, d time.Duration) error

	mu      sync.Mutex
	token   string
	expiry  time.Time
	refresh *tokenRefresh
}

type tokenRefresh struct {
	done  chan struct{}
	token string
	err   error
}

type tokenResponse struct {
	AccessToken string   `json:"accessToken"`
	Expires     DateTime `json:"expires"`
}

func (s *appTokenSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	if s.hasFreshToken() {
		token := s.token
		s.mu.Unlock()
		return token, nil
	}
	refresh := s.refreshInFlight(ctx)
	s.mu.Unlock()

	select {
	case <-refresh.done:
		return refresh.token, refresh.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// refreshInFlight must be called with s.mu held.
func (s *appTokenSource) refreshInFlight(ctx context.Context) *tokenRefresh {
	if s.refresh == nil {
		s.refresh = &tokenRefresh{done: make(chan struct{})}
		go s.runRefresh(context.WithoutCancel(ctx), s.refresh)
	}
	return s.refresh
}

func (s *appTokenSource) runRefresh(ctx context.Context, refresh *tokenRefresh) {
	token, expiry, err := s.fetchToken(ctx)
	s.mu.Lock()
	if err == nil {
		s.token, s.expiry = token, expiry
	}
	s.refresh = nil
	s.mu.Unlock()
	refresh.token, refresh.err = token, err
	close(refresh.done)
}

func (s *appTokenSource) Invalidate(stale string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.token == stale {
		s.token = ""
		s.expiry = time.Time{}
	}
}

func (s *appTokenSource) hasFreshToken() bool {
	return s.token != "" && s.now().Add(tokenRefreshMargin).Before(s.expiry)
}

func (s *appTokenSource) fetchToken(ctx context.Context) (string, time.Time, error) {
	for rateLimitRetries := 0; ; rateLimitRetries++ {
		statusCode, header, body, err := s.requestToken(ctx)
		if err != nil {
			return "", time.Time{}, err
		}
		if statusCode == http.StatusTooManyRequests && rateLimitRetries < maxRateLimitRetries {
			if err := s.sleep(ctx, retryAfter(header, body)); err != nil {
				return "", time.Time{}, err
			}
			continue
		}
		if !isSuccess(statusCode) {
			return "", time.Time{}, newError(statusCode, body)
		}
		return s.decodeToken(body)
	}
}

func (s *appTokenSource) requestToken(ctx context.Context) (int, http.Header, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.tokenURL, nil)
	if err != nil {
		return 0, nil, nil, err
	}
	req.SetBasicAuth(s.appID, s.appSecret)
	req.Header.Set("Accept", jsonContentType)
	req.Header.Set("User-Agent", s.userAgent)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("blikk: requesting access token: %w", err)
	}
	statusCode, body, err := readResponse(resp)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("blikk: reading access token response: %w", err)
	}
	return statusCode, resp.Header, body, nil
}

func (s *appTokenSource) decodeToken(body []byte) (string, time.Time, error) {
	var decoded tokenResponse
	if err := json.Unmarshal(body, &decoded); err != nil {
		return "", time.Time{}, fmt.Errorf("blikk: decoding access token response: %w", err)
	}
	if decoded.AccessToken == "" {
		return "", time.Time{}, errors.New("blikk: access token response has no accessToken")
	}
	return decoded.AccessToken, s.expiryOf(decoded), nil
}

func (s *appTokenSource) expiryOf(response tokenResponse) time.Time {
	if response.Expires.IsZero() {
		return s.now().Add(defaultTokenLifetime)
	}
	return response.Expires.Time()
}
