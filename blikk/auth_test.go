package blikk

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

type fakeBlikk struct {
	tokenRequests  atomic.Int32
	apiRequests    atomic.Int32
	tokenExpiry    time.Time
	rejectedTokens map[string]bool
	tokenRequest   *http.Request
	tokenBody      []byte
}

func (f *fakeBlikk) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/v1/Auth/Token" {
		f.serveToken(w, r)
		return
	}
	f.apiRequests.Add(1)
	token := r.Header.Get("Authorization")[len("Bearer "):]
	if f.rejectedTokens[token] {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"Invalid token"}`))
		return
	}
	_, _ = w.Write([]byte(`{"id":1,"firstName":"` + token + `"}`))
}

func (f *fakeBlikk) serveToken(w http.ResponseWriter, r *http.Request) {
	n := f.tokenRequests.Add(1)
	f.tokenRequest = r
	f.tokenBody, _ = io.ReadAll(r.Body)
	expires := f.tokenExpiry.UTC().Format("2006-01-02 15:04:05 -07:00")
	fmt.Fprintf(w, `{"objectName":"auth.token","accessToken":"token-%d","expires":"%s"}`, n, expires)
}

func newAppCredentialClient(t *testing.T, fake *fakeBlikk) (*Client, *appTokenSource) {
	t.Helper()
	server := httptest.NewServer(fake)
	t.Cleanup(server.Close)
	client, err := NewClient("my-app", "my-secret", WithBaseURL(server.URL), WithRateLimit(rate.Inf))
	require.NoError(t, err)
	return client, client.tokenSource.(*appTokenSource)
}

func TestTokenIsFetchedWithBasicAuthAndCached(t *testing.T) {
	now := time.Date(2026, 10, 13, 9, 40, 57, 0, time.UTC)
	fake := &fakeBlikk{tokenExpiry: now.Add(time.Hour)}
	_, source := newAppCredentialClient(t, fake)
	source.now = func() time.Time { return now }

	first, err := source.Token(context.Background())
	require.NoError(t, err)
	second, err := source.Token(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "token-1", first)
	assert.Equal(t, "token-1", second)
	assert.EqualValues(t, 1, fake.tokenRequests.Load())
	assert.Equal(t, http.MethodPost, fake.tokenRequest.Method)
	wantAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("my-app:my-secret"))
	assert.Equal(t, wantAuth, fake.tokenRequest.Header.Get("Authorization"))
	assert.Empty(t, fake.tokenBody)
	assert.Equal(t, now.Add(time.Hour), source.expiry)
}

func TestTokenIsRefreshedWithinSixtySecondsOfExpiry(t *testing.T) {
	now := time.Date(2026, 10, 13, 9, 40, 57, 0, time.UTC)
	fake := &fakeBlikk{tokenExpiry: now.Add(time.Hour)}
	_, source := newAppCredentialClient(t, fake)
	source.now = func() time.Time { return now }

	_, err := source.Token(context.Background())
	require.NoError(t, err)

	now = now.Add(time.Hour - 61*time.Second)
	token, err := source.Token(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "token-1", token)

	now = now.Add(2 * time.Second)
	token, err = source.Token(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "token-2", token)
	assert.EqualValues(t, 2, fake.tokenRequests.Load())
}

func TestUnauthorizedResponseInvalidatesTokenAndRetriesOnce(t *testing.T) {
	fake := &fakeBlikk{tokenExpiry: time.Now().Add(time.Hour), rejectedTokens: map[string]bool{"token-1": true}}
	client, _ := newAppCredentialClient(t, fake)

	user, err := client.Users.Get(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, "token-2", user.FirstName)
	assert.EqualValues(t, 2, fake.tokenRequests.Load())
	assert.EqualValues(t, 2, fake.apiRequests.Load())
}

func TestPersistentUnauthorizedIsReturnedAfterOneRetry(t *testing.T) {
	fake := &fakeBlikk{tokenExpiry: time.Now().Add(time.Hour), rejectedTokens: map[string]bool{"token-1": true, "token-2": true, "token-3": true}}
	client, _ := newAppCredentialClient(t, fake)

	_, err := client.Users.Get(context.Background(), 1)

	require.Error(t, err)
	assert.True(t, IsUnauthorized(err))
	assert.EqualValues(t, 2, fake.apiRequests.Load())
}

func TestTokenEndpointErrorIsReturnedAsError(t *testing.T) {
	server := httptest.NewServer(respondWith(http.StatusUnauthorized, ""))
	t.Cleanup(server.Close)
	client, err := NewClient("bad", "creds", WithBaseURL(server.URL), WithRateLimit(rate.Inf))
	require.NoError(t, err)

	_, err = client.Users.Get(context.Background(), 1)

	assert.True(t, IsUnauthorized(err))
}

func TestTokenWithoutExpiryFallsBackToOneHour(t *testing.T) {
	now := time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC)
	server := httptest.NewServer(respondWith(http.StatusOK, `{"accessToken":"abc"}`))
	t.Cleanup(server.Close)
	source := &appTokenSource{tokenURL: server.URL, httpClient: server.Client(), now: func() time.Time { return now }, sleep: sleepContext}

	token, err := source.Token(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "abc", token)
	assert.Equal(t, now.Add(time.Hour), source.expiry)
}

func TestTokenResponseWithoutAccessTokenIsAnError(t *testing.T) {
	server := httptest.NewServer(respondWith(http.StatusOK, `{"expires":"2026-10-13 10:40:57 +00:00"}`))
	t.Cleanup(server.Close)
	source := &appTokenSource{tokenURL: server.URL, httpClient: server.Client(), now: time.Now, sleep: sleepContext}

	_, err := source.Token(context.Background())

	assert.ErrorContains(t, err, "accessToken")
}

func TestNewClientRequiresCredentialsWithoutTokenSource(t *testing.T) {
	_, err := NewClient("", "secret")
	assert.Error(t, err)

	_, err = NewClient("", "", WithTokenSource(StaticToken("abc")))
	assert.NoError(t, err)
}

func TestStaticToken(t *testing.T) {
	source := StaticToken("abc")
	source.Invalidate("abc")
	token, err := source.Token(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "abc", token)
}

func TestTokenEndpointRateLimitIsRetried(t *testing.T) {
	server := httptest.NewServer(rateLimitedThen(2, "1", respondWith(http.StatusOK, `{"accessToken":"abc","expires":"2026-10-13 10:40:57 +00:00"}`)))
	t.Cleanup(server.Close)
	sleeps := &recordedSleeps{}
	source := &appTokenSource{tokenURL: server.URL, httpClient: server.Client(), now: time.Now, sleep: sleeps.sleep}

	token, err := source.Token(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "abc", token)
	assert.Equal(t, []time.Duration{time.Second, time.Second}, sleeps.durations)
}

func TestTokenEndpointRateLimitGivesUpAfterFiveRetries(t *testing.T) {
	server := httptest.NewServer(rateLimitedThen(100, "1", respondWith(http.StatusOK, `{}`)))
	t.Cleanup(server.Close)
	sleeps := &recordedSleeps{}
	source := &appTokenSource{tokenURL: server.URL, httpClient: server.Client(), now: time.Now, sleep: sleeps.sleep}

	_, err := source.Token(context.Background())

	assert.True(t, IsRateLimited(err))
	assert.Len(t, sleeps.durations, 5)
}

func TestInvalidateIgnoresTokenThatIsNoLongerCached(t *testing.T) {
	fake := &fakeBlikk{tokenExpiry: time.Now().Add(time.Hour)}
	_, source := newAppCredentialClient(t, fake)
	_, err := source.Token(context.Background())
	require.NoError(t, err)
	source.Invalidate("token-1")
	_, err = source.Token(context.Background())
	require.NoError(t, err)

	source.Invalidate("token-1")
	token, err := source.Token(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "token-2", token)
	assert.EqualValues(t, 2, fake.tokenRequests.Load())
}

func TestConcurrentUnauthorizedResponsesCauseOneRefresh(t *testing.T) {
	fake := &fakeBlikk{tokenExpiry: time.Now().Add(time.Hour), rejectedTokens: map[string]bool{"token-1": true}}
	client, source := newAppCredentialClient(t, fake)
	_, err := source.Token(context.Background())
	require.NoError(t, err)

	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := client.Users.Get(context.Background(), 1)
			assert.NoError(t, err)
		}()
	}
	wg.Wait()

	assert.EqualValues(t, 2, fake.tokenRequests.Load())
}

func TestConcurrentTokenCallersShareOneFetch(t *testing.T) {
	var fetches atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetches.Add(1)
		time.Sleep(20 * time.Millisecond)
		_, _ = w.Write([]byte(`{"accessToken":"shared","expires":"2099-01-01 00:00:00 +00:00"}`))
	}))
	t.Cleanup(server.Close)
	source := &appTokenSource{tokenURL: server.URL, httpClient: server.Client(), now: time.Now, sleep: sleepContext}

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := source.Token(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, "shared", token)
		}()
	}
	wg.Wait()

	assert.EqualValues(t, 1, fetches.Load())
}

func TestTokenWaiterWithExpiredContextReturnsWhileRefreshIsInFlight(t *testing.T) {
	fetchStarted := make(chan struct{})
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(fetchStarted)
		<-release
		_, _ = w.Write([]byte(`{"accessToken":"slow","expires":"2099-01-01 00:00:00 +00:00"}`))
	}))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })
	source := &appTokenSource{tokenURL: server.URL, httpClient: server.Client(), now: time.Now, sleep: sleepContext}

	go func() { _, _ = source.Token(context.Background()) }()
	<-fetchStarted
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	_, err := source.Token(ctx)

	assert.ErrorIs(t, err, context.Canceled)
	assert.Less(t, time.Since(start), 100*time.Millisecond)
}
