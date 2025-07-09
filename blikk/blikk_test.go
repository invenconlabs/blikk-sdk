package blikk

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestServer creates a new mock server and a client pointing to it.
func setupTestServer(t *testing.T, handler http.Handler) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	client := NewClient("fake-token", WithBaseURL(server.URL+"/"), WithHTTPClient(server.Client()))
	return client, server
}

func TestList_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/Admin/Users", r.URL.Path)
		assert.Equal(t, "Bearer fake-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{
			"objectName": "list",
			"page": 1,
			"pageSize": 100,
			"itemCount": 1,
			"totalItemCount": 1,
			"totalPages": 1,
			"items": [{"objectName": "Users", "id": 1, "firstName": "Test"}]
		}`)
	})

	client, server := setupTestServer(t, handler)
	defer server.Close()

	users, err := List[Users](client, NewListOptions())
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, "Test", users[0].FirstName)
}

func TestList_Pagination(t *testing.T) {
	page := 1
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, fmt.Sprintf("%d", page), r.URL.Query().Get("page"))
		w.WriteHeader(http.StatusOK)
		if page == 1 {
			fmt.Fprintln(w, `{
				"totalPages": 2, "page": 1,
				"items": [{"id": 1}]
			}`)
			page++
		} else {
			fmt.Fprintln(w, `{
				"totalPages": 2, "page": 2,
				"items": [{"id": 2}]
			}`)
		}
	})

	client, server := setupTestServer(t, handler)
	defer server.Close()

	opts := NewListOptions()
	items, err := List[Users](client, opts)
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, 1, items[0].ID)
	assert.Equal(t, 2, items[1].ID)
}

func TestGet_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/Admin/Users/123", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"id": 123, "firstName": "Specific"}`)
	})

	client, server := setupTestServer(t, handler)
	defer server.Close()

	user, err := Get[User](client, "123")
	require.NoError(t, err)
	assert.Equal(t, 123, user.ID)
	assert.Equal(t, "Specific", user.FirstName)
}

func TestClient_RetryRequest(t *testing.T) {
	attempts := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.Header().Set("Retry-After", "1") // 1 second
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"id": 1}`)
	})

	// Override the HTTP client to have a very short timeout for the test
	httpClient := &http.Client{Timeout: 2 * time.Second}
	server := httptest.NewServer(handler)
	defer server.Close()

	client := NewClient("fake-token", WithBaseURL(server.URL+"/"), WithHTTPClient(httpClient))

	_, err := Get[User](client, "1")
	require.NoError(t, err)
	assert.Equal(t, 2, attempts, "Expected the client to make two attempts")
}

func TestClient_ServerError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "internal server error")
	})

	client, server := setupTestServer(t, handler)
	defer server.Close()

	_, err := List[Users](client, NewListOptions())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code 500")
	assert.Contains(t, err.Error(), "internal server error")
}

func TestClient_BadJSON(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"id": 1, "firstName": "Test"`) // Malformed JSON
	})

	client, server := setupTestServer(t, handler)
	defer server.Close()

	_, err := Get[User](client, "1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal response")
}
