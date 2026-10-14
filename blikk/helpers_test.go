package blikk

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const testToken = "test-token"

func newTestClient(t *testing.T, handler http.HandlerFunc, opts ...Option) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	defaults := []Option{
		WithBaseURL(server.URL),
		WithTokenSource(StaticToken(testToken)),
		WithRequestsPerSecond(0),
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
	r.Body = io.NopCloser(bytes.NewReader(body))
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

func (h *recordingHandler) all() []recordedRequest {
	h.mu.Lock()
	defer h.mu.Unlock()
	return slices.Clone(h.requests)
}

func newRecordingClient(t *testing.T, respond http.HandlerFunc, opts ...Option) (*Client, *recordingHandler) {
	t.Helper()
	handler := &recordingHandler{respond: respond}
	return newTestClient(t, handler.ServeHTTP, opts...), handler
}

// routes maps "METHOD /escaped/path" (no query string) to a handler; any other request fails the test with a 404.
func routes(t *testing.T, handlers map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler, found := handlers[r.Method+" "+r.URL.EscapedPath()]
		if !found {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.EscapedPath())
			http.NotFound(w, r)
			return
		}
		handler(w, r)
	}
}

func serveFixtureAt(t *testing.T, method, path, name string) http.HandlerFunc {
	return routes(t, map[string]http.HandlerFunc{method + " " + path: serveFixture(t, name)})
}

func serveFile(contentType, contentDisposition string, content []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		if contentDisposition != "" {
			w.Header().Set("Content-Disposition", contentDisposition)
		}
		_, _ = w.Write(content)
	}
}

func (r recordedRequest) decodeBody(t *testing.T, out any) {
	t.Helper()
	require.NoError(t, json.Unmarshal(r.Body, out), "request body: %s", r.Body)
}

func (r recordedRequest) multipartForm(t *testing.T) *multipart.Form {
	t.Helper()
	request := httptest.NewRequest(r.Method, r.Path, bytes.NewReader(r.Body))
	request.Header = r.Header.Clone()
	require.NoError(t, request.ParseMultipartForm(1<<20), "request body: %s", r.Body)
	return request.MultipartForm
}

func formValues(form *multipart.Form) map[string]string {
	values := make(map[string]string, len(form.Value))
	for name, fieldValues := range form.Value {
		values[name] = strings.Join(fieldValues, ",")
	}
	return values
}

type uploadedPart struct {
	Filename    string
	ContentType string
	Content     string
}

func uploadedParts(t *testing.T, form *multipart.Form, fieldName string) []uploadedPart {
	t.Helper()
	parts := make([]uploadedPart, 0, len(form.File[fieldName]))
	for _, header := range form.File[fieldName] {
		file, err := header.Open()
		require.NoError(t, err)
		content, err := io.ReadAll(file)
		require.NoError(t, err)
		require.NoError(t, file.Close())
		parts = append(parts, uploadedPart{Filename: header.Filename, ContentType: header.Header.Get("Content-Type"), Content: string(content)})
	}
	return parts
}

func (r recordedRequest) jsonBody(t *testing.T) map[string]any {
	t.Helper()
	var decoded map[string]any
	r.decodeBody(t, &decoded)
	return decoded
}
