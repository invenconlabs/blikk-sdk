package blikk

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestRequestCarriesAuthAndStandardHeaders(t *testing.T) {
	handler := &recordingHandler{respond: respondWith(http.StatusOK, `{"id":42}`)}
	client := newTestClient(t, handler.ServeHTTP, WithUserAgent("my-app/1.0"))

	_, err := client.Users.Get(context.Background(), 42)

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodGet, request.Method)
	assert.Equal(t, "/v1/Admin/Users/42", request.Path)
	assert.Equal(t, "Bearer "+testToken, request.Header.Get("Authorization"))
	assert.Equal(t, "application/json", request.Header.Get("Accept"))
	assert.Equal(t, "my-app/1.0", request.Header.Get("User-Agent"))
	assert.Empty(t, request.Header.Get("Content-Type"))
}

func TestBaseURLWithPathPrefixAndTrailingSlash(t *testing.T) {
	handler := &recordingHandler{respond: respondWith(http.StatusOK, `{}`)}
	client := newTestClient(t, handler.ServeHTTP)
	client.baseURL += "/proxy/"

	_, err := client.Projects.Get(context.Background(), 7)

	require.NoError(t, err)
	assert.Equal(t, "/proxy/v1/Core/Projects/7", handler.last().Path)
}

func TestPathParametersAreEscaped(t *testing.T) {
	handler := &recordingHandler{respond: respondWith(http.StatusOK, `{}`)}
	client := newTestClient(t, handler.ServeHTTP)

	_, err := client.Projects.GetByOrderNumber(context.Background(), "100/5 a?b")

	require.NoError(t, err)
	assert.Equal(t, "/v1/Core/Projects/order/100%2F5%20a%3Fb", handler.last().Path)
}

func TestRateLimiterIsConsultedBeforeEveryRequest(t *testing.T) {
	handler := &recordingHandler{respond: respondWith(http.StatusOK, `{}`)}
	client := newTestClient(t, handler.ServeHTTP, WithRequestsPerSecond(0.001))

	_, err := client.Users.Get(context.Background(), 1)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err = client.Users.Get(ctx, 1)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.True(t, strings.HasPrefix(err.Error(), "blikk: "), err.Error())
	assert.Len(t, handler.requests, 1)
}

func TestDefaultRateLimitIsFourPerSecondWithBurstOne(t *testing.T) {
	client, err := NewClient("id", "secret")
	require.NoError(t, err)

	assert.Equal(t, rate.Limit(4), client.limiter.Limit())
	assert.Equal(t, 1, client.limiter.Burst())
}

func TestRequestsPerSecondOption(t *testing.T) {
	limited, err := NewClient("id", "secret", WithRequestsPerSecond(2.5))
	require.NoError(t, err)
	unlimited, err := NewClient("id", "secret", WithRequestsPerSecond(0))
	require.NoError(t, err)

	assert.Equal(t, rate.Limit(2.5), limited.limiter.Limit())
	assert.Equal(t, rate.Inf, unlimited.limiter.Limit())
}

func TestNewClientWithTokenSource(t *testing.T) {
	handler := &recordingHandler{respond: respondWith(http.StatusOK, `{"id":3}`)}
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := NewClientWithTokenSource(StaticToken("given"), WithBaseURL(server.URL), WithRequestsPerSecond(0))
	require.NoError(t, err)
	_, err = client.Users.Get(context.Background(), 3)

	require.NoError(t, err)
	assert.Equal(t, "Bearer given", handler.last().Header.Get("Authorization"))

	_, err = NewClientWithTokenSource(nil)
	assert.Error(t, err)
}

type recordedSleeps struct {
	mu        sync.Mutex
	durations []time.Duration
}

func (s *recordedSleeps) sleep(_ context.Context, d time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.durations = append(s.durations, d)
	return nil
}

func rateLimitedThen(failures int32, retryAfter string, success http.HandlerFunc) http.HandlerFunc {
	var calls atomic.Int32
	return func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) <= failures {
			if retryAfter != "" {
				w.Header().Set("Retry-After", retryAfter)
			}
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"waitInSeconds":2}`))
			return
		}
		success(w, r)
	}
}

func TestTooManyRequestsIsRetriedAfterRetryAfter(t *testing.T) {
	handler := &recordingHandler{respond: rateLimitedThen(2, "0.01", respondWith(http.StatusOK, `{"id":5}`))}
	client := newTestClient(t, handler.ServeHTTP)
	sleeps := &recordedSleeps{}
	client.sleep = sleeps.sleep

	user, err := client.Users.Get(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, 5, user.ID)
	assert.Len(t, handler.requests, 3)
	assert.Equal(t, []time.Duration{10 * time.Millisecond, 10 * time.Millisecond}, sleeps.durations)
}

func TestTooManyRequestsGivesUpAfterFiveRetries(t *testing.T) {
	handler := &recordingHandler{respond: rateLimitedThen(100, "0", respondWith(http.StatusOK, `{}`))}
	client := newTestClient(t, handler.ServeHTTP)
	sleeps := &recordedSleeps{}
	client.sleep = sleeps.sleep

	_, err := client.Users.Get(context.Background(), 5)

	require.Error(t, err)
	assert.True(t, IsRateLimited(err))
	assert.Len(t, handler.requests, 6)
	assert.Len(t, sleeps.durations, 5)
}

func TestRetryAfterIsCappedAndFallsBackToBody(t *testing.T) {
	assert.Equal(t, 30*time.Second, retryAfter(http.Header{"Retry-After": {"120"}}, nil))
	assert.Equal(t, 2*time.Second, retryAfter(http.Header{}, []byte(`{"waitInSeconds":2}`)))
	assert.Equal(t, time.Second, retryAfter(http.Header{}, nil))
	assert.Equal(t, time.Second, retryAfter(http.Header{"Retry-After": {"-3"}}, nil))
}

func TestUnusableRetryAfterHeaderIsIgnored(t *testing.T) {
	for _, header := range []string{"1e10", "Inf", "NaN", "-1"} {
		got := retryAfter(http.Header{"Retry-After": {header}}, []byte(`{"waitInSeconds":2}`))
		want := 2 * time.Second
		if header == "1e10" {
			want = 30 * time.Second
		}
		assert.Equal(t, want, got, "Retry-After %q", header)
	}
}

func TestRetrySleepHonoursContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	err := sleepContext(ctx, time.Minute)

	assert.ErrorIs(t, err, context.Canceled)
	assert.Less(t, time.Since(start), time.Second)
}

func TestCancelledContextStopsRateLimitRetries(t *testing.T) {
	handler := &recordingHandler{respond: rateLimitedThen(100, "10", respondWith(http.StatusOK, `{}`))}
	client := newTestClient(t, handler.ServeHTTP)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.Users.Get(ctx, 5)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Len(t, handler.requests, 1)
}

func TestAnySuccessStatusIsAccepted(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusCreated, http.StatusNoContent} {
		client := newTestClient(t, respondWith(status, ""))
		assert.NoError(t, client.TimeReports.Delete(context.Background(), 1), "status %d", status)
	}
}

func TestMalformedSuccessBodyIsADecodeError(t *testing.T) {
	client := newTestClient(t, respondWith(http.StatusOK, `{"id":"not a number"}`))

	_, err := client.Users.Get(context.Background(), 1)

	assert.ErrorContains(t, err, "decoding response of GET v1/Admin/Users/1")
}

func TestNewClientRejectsInvalidBaseURL(t *testing.T) {
	for _, baseURL := range []string{"", "publicapi.blikk.com", "://bad", "http://publicapi.blikk.com", "http://10.0.0.1:8080", "ftp://publicapi.blikk.com"} {
		_, err := NewClient("id", "secret", WithBaseURL(baseURL))
		assert.Error(t, err, "base URL %q", baseURL)
	}
}

func TestNewClientAcceptsHTTPSAndLoopbackHTTP(t *testing.T) {
	for _, baseURL := range []string{"https://publicapi.blikk.com", "http://127.0.0.1:8080", "http://localhost", "http://[::1]:9000/proxy"} {
		_, err := NewClient("id", "secret", WithBaseURL(baseURL))
		assert.NoError(t, err, "base URL %q", baseURL)
	}
}

func TestUnsafePathParametersAreRejectedWithoutARequest(t *testing.T) {
	handler := &recordingHandler{respond: respondWith(http.StatusOK, `{"page":1,"totalPages":1,"items":[]}`)}
	client := newTestClient(t, handler.ServeHTTP)

	for _, orderNumber := range []string{"..", ".", ""} {
		project, err := client.Projects.GetByOrderNumber(context.Background(), orderNumber)
		assert.Error(t, err, "order number %q", orderNumber)
		assert.Nil(t, project)
	}
	assert.Empty(t, handler.requests)
}

func TestPathParametersWithSlashesAndDotsStayOneSegment(t *testing.T) {
	handler := &recordingHandler{respond: respondWith(http.StatusOK, `{}`)}
	client := newTestClient(t, handler.ServeHTTP)

	for orderNumber, wantPath := range map[string]string{
		"100/5":   "/v1/Core/Projects/order/100%2F5",
		"../x":    "/v1/Core/Projects/order/..%2Fx",
		"a/../..": "/v1/Core/Projects/order/a%2F..%2F..",
		"..5":     "/v1/Core/Projects/order/..5",
	} {
		_, err := client.Projects.GetByOrderNumber(context.Background(), orderNumber)
		require.NoError(t, err)
		assert.Equal(t, wantPath, handler.last().Path, "order number %q", orderNumber)
	}
}

func TestDownloadStreamsBodyWithContentTypeAndFilename(t *testing.T) {
	disposition := `attachment; filename=offer.pdf; filename*=UTF-8''off%C3%A9r%201.pdf`
	client, handler := newRecordingClient(t, routes(t, map[string]http.HandlerFunc{"GET /v1/Core/Offers/1/Download": serveFile("application/pdf", disposition, []byte("%PDF-1.7"))}))

	download, err := client.download(context.Background(), "v1/Core/Offers/1/Download", url.Values{"lang": {"sv"}})

	require.NoError(t, err)
	defer download.Content.Close()
	content, err := io.ReadAll(download.Content)
	require.NoError(t, err)
	assert.Equal(t, "%PDF-1.7", string(content))
	assert.Equal(t, "application/pdf", download.ContentType)
	assert.Equal(t, "offér 1.pdf", download.Filename)
	request := handler.last()
	assert.Equal(t, map[string][]string{"lang": {"sv"}}, request.Query)
	assert.Equal(t, "*/*", request.Header.Get("Accept"))
	assert.Equal(t, "Bearer "+testToken, request.Header.Get("Authorization"))
}

func TestDownloadWithoutContentDispositionHasNoFilename(t *testing.T) {
	client := newTestClient(t, serveFile("image/png", "", []byte{0x89, 'P', 'N', 'G'}))

	download, err := client.download(context.Background(), "v1/Core/Tasks/1/Attachments/2/Download", nil)

	require.NoError(t, err)
	defer download.Content.Close()
	assert.Equal(t, "image/png", download.ContentType)
	assert.Empty(t, download.Filename)
}

func TestContentDispositionFilenameDropsDirectories(t *testing.T) {
	for header, want := range map[string]string{
		`attachment; filename="report.pdf"`:       "report.pdf",
		`attachment; filename="../../etc/passwd"`: "passwd",
		`attachment; filename="C:\\temp\\a.pdf"`:  "a.pdf",
		`attachment; filename=".."`:               "",
		`attachment`:                              "",
		`not a disposition;;`:                     "",
		``:                                        "",
	} {
		assert.Equal(t, want, filenameFromContentDisposition(header), "header %q", header)
	}
}

func TestDownloadReturnsBeforeTheBodyIsComplete(t *testing.T) {
	release := make(chan struct{})
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("first"))
		w.(http.Flusher).Flush()
		<-release
		_, _ = w.Write([]byte("second"))
	})

	download, err := client.download(context.Background(), "v1/core/invoices/42/attachments/1/download", nil)
	require.NoError(t, err)
	defer download.Content.Close()
	first := make([]byte, len("first"))
	_, err = io.ReadFull(download.Content, first)
	require.NoError(t, err)
	close(release)
	rest, err := io.ReadAll(download.Content)

	require.NoError(t, err)
	assert.Equal(t, "first", string(first))
	assert.Equal(t, "second", string(rest))
}

func TestDownloadRetriesTooManyRequests(t *testing.T) {
	handler := &recordingHandler{respond: rateLimitedThen(2, "0.01", serveFile("application/pdf", "", []byte("pdf")))}
	client := newTestClient(t, handler.ServeHTTP)
	sleeps := &recordedSleeps{}
	client.sleep = sleeps.sleep

	download, err := client.download(context.Background(), "v1/Core/Offers/1/Download", nil)

	require.NoError(t, err)
	defer download.Content.Close()
	content, err := io.ReadAll(download.Content)
	require.NoError(t, err)
	assert.Equal(t, "pdf", string(content))
	assert.Len(t, handler.requests, 3)
	assert.Len(t, sleeps.durations, 2)
}

func TestDownloadRefreshesRejectedToken(t *testing.T) {
	fake := &fakeBlikk{tokenExpiry: time.Now().Add(time.Hour), rejectedTokens: map[string]bool{"token-1": true}}
	client, _ := newAppCredentialClient(t, fake)

	download, err := client.download(context.Background(), "v1/Core/Offers/1/Download", nil)

	require.NoError(t, err)
	defer download.Content.Close()
	content, err := io.ReadAll(download.Content)
	require.NoError(t, err)
	assert.Contains(t, string(content), "token-2")
	assert.EqualValues(t, 2, fake.apiRequests.Load())
}

func TestDownloadErrorIsAnAPIError(t *testing.T) {
	client := newTestClient(t, respondWith(http.StatusNotFound, `{"message":"Attachment could not be found"}`))

	download, err := client.download(context.Background(), "v1/core/invoices/42/attachments/9/download", nil)

	assert.Nil(t, download)
	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, "/v1/core/invoices/42/attachments/9/download", apiErr.Path)
	assert.Equal(t, "Attachment could not be found", apiErr.Message)
	assert.True(t, IsNotFound(err))
}

func TestEveryServiceIsWiredToTheClient(t *testing.T) {
	client := newTestClient(t, respondWith(http.StatusOK, `{}`))

	clientValue := reflect.ValueOf(client).Elem()
	for i := range clientValue.NumField() {
		field := clientValue.Type().Field(i)
		if !field.IsExported() || field.Type.Kind() != reflect.Pointer {
			continue
		}
		service := clientValue.Field(i)
		require.False(t, service.IsNil(), "service %s is nil", field.Name)
		backReference := service.Elem().FieldByName("client")
		require.True(t, backReference.IsValid(), "service %s has no client field", field.Name)
		assert.Equal(t, client, (*Client)(backReference.UnsafePointer()), "service %s", field.Name)
	}
}

func TestUploadResendsTheSameMultipartBodyAfterTooManyRequests(t *testing.T) {
	handler := &recordingHandler{respond: rateLimitedThen(1, "0", respondWith(http.StatusCreated, `{"id":7}`))}
	client := newTestClient(t, handler.ServeHTTP)
	client.sleep = (&recordedSleeps{}).sleep
	var created createdResource

	err := client.upload(context.Background(), "v1/uploads", map[string]string{"title": "Invoice"},
		[]UploadFile{{FieldName: "files", Filename: "a.pdf", Content: strings.NewReader("pdf")}}, &created)

	require.NoError(t, err)
	assert.Equal(t, 7, created.ID)
	requests := handler.all()
	require.Len(t, requests, 2)
	assert.Equal(t, requests[0].Body, requests[1].Body)
	form := requests[1].multipartForm(t)
	assert.Equal(t, map[string]string{"title": "Invoice"}, formValues(form))
	assert.Equal(t, []uploadedPart{{Filename: "a.pdf", ContentType: "application/octet-stream", Content: "pdf"}}, uploadedParts(t, form, "files"))
}

func TestUploadQuotesFieldAndFileNames(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	err := client.upload(context.Background(), "v1/uploads", nil,
		[]UploadFile{{FieldName: "file", Filename: `räkning "maj".pdf`, Content: strings.NewReader("pdf")}}, nil)

	require.NoError(t, err)
	parts := uploadedParts(t, handler.last().multipartForm(t), "file")
	require.Len(t, parts, 1)
	assert.Equal(t, `räkning "maj".pdf`, parts[0].Filename)
}

func TestUploadWithoutContentFailsBeforeSending(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	err := client.upload(context.Background(), "v1/uploads", nil, []UploadFile{{FieldName: "file", Filename: "a.png"}}, nil)

	require.ErrorContains(t, err, `upload file "a.png" has no content`)
	assert.Empty(t, handler.all())
}
