package blikk

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

const testToken = "test-token"

func newTestClient(t *testing.T, handler http.HandlerFunc, opts ...Option) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	defaults := []Option{
		WithBaseURL(server.URL),
		WithTokenSource(StaticToken(testToken)),
		WithRateLimit(rate.Inf),
	}
	client, err := NewClient("", "", append(defaults, opts...)...)
	require.NoError(t, err)
	return client
}

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return data
}

func serveFixture(t *testing.T, name string) http.HandlerFunc {
	fixture := loadFixture(t, name)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(fixture)
	}
}

type recordedRequest struct {
	Method string
	Path   string
	Query  map[string][]string
	Header http.Header
	Body   []byte
}

func respondWith(statusCode int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(body))
	}
}

type recordingHandler struct {
	mu       sync.Mutex
	requests []recordedRequest
	respond  http.HandlerFunc
}

func (h *recordingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	h.mu.Lock()
	h.requests = append(h.requests, recordedRequest{
		Method: r.Method,
		Path:   r.URL.EscapedPath(),
		Query:  r.URL.Query(),
		Header: r.Header.Clone(),
		Body:   body,
	})
	h.mu.Unlock()
	if h.respond != nil {
		h.respond(w, r)
	}
}

func (h *recordingHandler) last() recordedRequest {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.requests[len(h.requests)-1]
}
