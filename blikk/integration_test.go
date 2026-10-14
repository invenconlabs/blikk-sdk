//go:build integration

package blikk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sharedIntegrationClient = sync.OnceValues(func() (*Client, error) {
	return NewClient(os.Getenv("BLIKK_APP_ID"), os.Getenv("BLIKK_APP_SECRET"))
})

func integrationClient(t *testing.T) *Client {
	t.Helper()
	if os.Getenv("BLIKK_APP_ID") == "" || os.Getenv("BLIKK_APP_SECRET") == "" {
		t.Skip("BLIKK_APP_ID and BLIKK_APP_SECRET must be set for integration tests")
	}
	client, err := sharedIntegrationClient()
	require.NoError(t, err)
	return client
}

func integrationContext(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)
	return ctx
}

func TestIntegrationUsers(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.Users.List(ctx, UserListParams{ListOptions: ListOptions{PageSize: 10, SortBy: UserSortByName}})
	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	assert.Equal(t, 1, page.Page)
	assert.Positive(t, page.TotalItemCount)
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)
	assert.NotEmpty(t, listed.FirstName)
	assert.NotNil(t, listed.CreatedDate)

	user, err := client.Users.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, user.ID)
	assert.Equal(t, listed.FirstName, user.FirstName)
	assert.Equal(t, listed.LastName, user.LastName)
	assert.Equal(t, listed.SalaryType, user.SalaryType)
	assert.NotEmpty(t, user.License)
}

func TestIntegrationProjects(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.Projects.List(ctx, ProjectListParams{ListOptions: ListOptions{PageSize: 10, SortBy: ProjectSortByUpdatedDate, SortOrder: SortDescending}})
	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)
	assert.NotEmpty(t, listed.OrderNumber)
	assert.NotEmpty(t, listed.Title)
	require.NotNil(t, listed.Status)
	assert.NotZero(t, listed.Status.ID)
	assert.NotEmpty(t, listed.InvoiceType)

	project, err := client.Projects.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, project.ID)
	assert.Equal(t, listed.OrderNumber, project.OrderNumber)
	assert.Equal(t, listed.Title, project.Title)

	byOrderNumber, err := client.Projects.GetByOrderNumber(ctx, listed.OrderNumber)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, byOrderNumber.ID)
}

func TestIntegrationProjectNotFound(t *testing.T) {
	client := integrationClient(t)

	_, err := client.Projects.Get(integrationContext(t), 999999999)

	assert.True(t, IsNotFound(err), "expected 404, got %v", err)
}

func TestIntegrationUserNotFound(t *testing.T) {
	client := integrationClient(t)

	_, err := client.Users.Get(integrationContext(t), 999999999)

	assert.True(t, IsNotFound(err), "expected not found, got %v", err)
}

func TestIntegrationTimeReports(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)
	today := time.Now()
	params := TimeReportListParams{
		ListOptions: ListOptions{PageSize: 25},
		From:        DateOf(today.AddDate(0, 0, -30)),
		To:          DateOf(today),
	}

	page, err := client.TimeReports.List(ctx, params)
	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("no time reports in the last 30 days")
	}
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)
	assert.False(t, listed.Date.IsZero())
	assert.NotZero(t, listed.User.ID)
	assert.False(t, listed.CreatedDate.IsZero())
	assert.False(t, listed.Date.Time().Before(params.From.Time()), "report date %s before filter start %s", listed.Date, params.From)

	report, err := client.TimeReports.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, report.ID)
	assert.Equal(t, listed.Date, report.Date)
	assert.Equal(t, listed.User, report.User)
	assert.Equal(t, listed.Hours, report.Hours)
}

func TestIntegrationTimeReportsListAllCrossesPages(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)
	today := time.Now()
	params := TimeReportListParams{
		ListOptions: ListOptions{PageSize: 5},
		From:        DateOf(today.AddDate(0, 0, -30)),
		To:          DateOf(today),
	}

	seen := map[int]bool{}
	for report, err := range client.TimeReports.ListAll(ctx, params) {
		require.NoError(t, err)
		assert.False(t, seen[report.ID], "duplicate report %d", report.ID)
		seen[report.ID] = true
		if len(seen) == 12 {
			break
		}
	}
	if len(seen) < 12 {
		t.Skip("fewer than 12 time reports in the last 30 days")
	}
}

func TestIntegrationUserDayStatistics(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)
	today := time.Now()

	page, err := client.TimeReports.ListUserDayStatistics(ctx, UserDayStatisticsParams{
		From: DateOf(today.AddDate(0, 0, -6)),
		To:   DateOf(today),
	})

	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	statistics := page.Items[0]
	assert.NotZero(t, statistics.UserID)
	assert.NotEmpty(t, statistics.Name)
	require.Len(t, statistics.Dates, 7)
	assert.False(t, statistics.Dates[0].Date.IsZero())
}

func TestIntegrationUserDayStatisticsRejectsLongRange(t *testing.T) {
	client := integrationClient(t)
	today := time.Now()

	_, err := client.TimeReports.ListUserDayStatistics(integrationContext(t), UserDayStatisticsParams{
		From: DateOf(today.AddDate(0, 0, -60)),
		To:   DateOf(today),
	})

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 400, apiErr.StatusCode)
	assert.NotEmpty(t, apiErr.ModelState)
}

const writeTestComment = "blikk-sdk integration test"

var writeTestWeekStart, writeTestWeekEnd = NewDate(2026, time.October, 19), NewDate(2026, time.October, 25)

func writeTestUserID(t *testing.T) int {
	t.Helper()
	raw := os.Getenv("BLIKK_WRITE_TEST_USER_ID")
	if raw == "" {
		t.Skip("BLIKK_WRITE_TEST_USER_ID not set; skipping live write test")
	}
	userID, err := strconv.Atoi(raw)
	require.NoError(t, err, "BLIKK_WRITE_TEST_USER_ID must be an integer")
	return userID
}

func countWriteTestWeekReports(t *testing.T, ctx context.Context, client *Client, userID int) int {
	t.Helper()
	page, err := client.TimeReports.List(ctx, TimeReportListParams{UserIDs: []int{userID}, From: writeTestWeekStart, To: writeTestWeekEnd})
	require.NoError(t, err)
	return page.TotalItemCount
}

func borrowTimeReportReferences(t *testing.T, ctx context.Context, client *Client, userID int) TimeReportInput {
	t.Helper()
	page, err := client.TimeReports.List(ctx, TimeReportListParams{
		ListOptions: ListOptions{PageSize: 1, SortBy: TimeReportSortByDate, SortOrder: SortDescending},
		UserIDs:     []int{userID},
	})
	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skipf("user %d has no time reports to borrow project, activity, time code and time article from", userID)
	}
	existing := page.Items[0]
	input := TimeReportInput{
		UserID:        userID,
		ActivityID:    existing.Activity.ID,
		TimeCodeID:    existing.TimeCode.ID,
		TimeArticleID: existing.TimeArticle.ID,
	}
	switch {
	case existing.Project != nil:
		input.ProjectID = existing.Project.ID
	case existing.InternalProject != nil:
		input.InternalProjectID = existing.InternalProject.ID
	case existing.AbsenceProject != nil:
		input.AbsenceProjectID = existing.AbsenceProject.ID
	}
	return input
}

const additionSearchLimit = 100

type recordedResponse struct {
	statusCode int
	body       []byte
}

type responseRecorder struct {
	mu        sync.Mutex
	responses map[string]recordedResponse
}

func (r *responseRecorder) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewReader(body))
	r.mu.Lock()
	r.responses[req.Method+" "+req.URL.Path] = recordedResponse{statusCode: resp.StatusCode, body: body}
	r.mu.Unlock()
	return resp, nil
}

func (r *responseRecorder) last(method, path string) recordedResponse {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.responses[method+" "+path]
}

func newRecordingWriteClient(t *testing.T) (*Client, *responseRecorder) {
	t.Helper()
	recorder := &responseRecorder{responses: map[string]recordedResponse{}}
	client, err := NewClient(os.Getenv("BLIKK_APP_ID"), os.Getenv("BLIKK_APP_SECRET"),
		WithHTTPClient(&http.Client{Transport: recorder, Timeout: 30 * time.Second}))
	require.NoError(t, err)
	return client, recorder
}

func createReportWithCleanup(t *testing.T, ctx context.Context, client *Client, input TimeReportInput) int {
	t.Helper()
	id, err := client.TimeReports.Create(ctx, input)
	require.NoError(t, err)
	require.NotZero(t, id)
	t.Cleanup(func() { deleteReportIfPresent(t, client, id) })
	return id
}

func createReportWithCleanupOrSkip(t *testing.T, client *Client, create func() (int, error)) int {
	t.Helper()
	id, err := create()
	var apiErr *Error
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusBadRequest {
		t.Logf("create rejected: %v; body %s", err, apiErr.Body)
		t.Skip("create returned 400; skipping")
	}
	require.NoError(t, err)
	require.NotZero(t, id)
	t.Cleanup(func() { deleteReportIfPresent(t, client, id) })
	return id
}

func deleteReportIfPresent(t *testing.T, client *Client, id int) {
	ctx := context.Background()
	_, err := client.TimeReports.Get(ctx, id)
	var apiErr *Error
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusBadRequest {
		return
	}
	if !assert.NoError(t, err, "cleanup could not read report %d; delete it by hand", id) {
		return
	}
	assert.NoError(t, client.TimeReports.Delete(ctx, id), "cleanup delete of report %d", id)
}

func assertReportIsGone(t *testing.T, ctx context.Context, client *Client, id int) {
	t.Helper()
	_, err := client.TimeReports.Get(ctx, id)
	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	t.Logf("Get after delete of %d returned status %d: %s", id, apiErr.StatusCode, apiErr.Body)
}

type additionArticle struct {
	AdditionType string `json:"additionType"`
}

func findAdditionToBorrow(t *testing.T, ctx context.Context, client *Client, userID int) (AdditionInput, bool) {
	t.Helper()
	for _, userIDs := range [][]int{{userID}, nil} {
		page, err := client.TimeReports.List(ctx, TimeReportListParams{
			ListOptions:  ListOptions{PageSize: additionSearchLimit, SortBy: TimeReportSortByDate, SortOrder: SortDescending},
			UserIDs:      userIDs,
			HasAdditions: Ptr(true),
		})
		require.NoError(t, err)
		for _, listed := range page.Items {
			report, err := client.TimeReports.Get(ctx, listed.ID)
			require.NoError(t, err)
			for _, addition := range report.Additions {
				if input, ok := additionInputFrom(t, ctx, client, addition); ok {
					t.Logf("borrowing addition article %d (%q) from report %d: %+v", addition.Article.ID, addition.Article.Name, report.ID, input)
					return input, true
				}
			}
		}
	}
	return AdditionInput{}, false
}

// The detail response has no yesNo/numberCount, so the input field is chosen by the article's additionType.
func additionInputFrom(t *testing.T, ctx context.Context, client *Client, addition AdditionReport) (AdditionInput, bool) {
	t.Helper()
	path, err := resourcePath("v1/admin/articles", strconv.Itoa(addition.Article.ID))
	require.NoError(t, err)
	var article additionArticle
	require.NoError(t, client.do(ctx, http.MethodGet, path, nil, nil, &article))
	switch {
	case strings.EqualFold(article.AdditionType, "Number"):
		return AdditionInput{ArticleID: addition.Article.ID, NumberCount: addition.SalaryCount}, true
	case strings.EqualFold(article.AdditionType, "Boolean"):
		return AdditionInput{ArticleID: addition.Article.ID, YesNo: true}, true
	default:
		t.Logf("skipping addition article %d with additionType %q", addition.Article.ID, article.AdditionType)
		return AdditionInput{}, false
	}
}

func additionCount(t *testing.T, ctx context.Context, client *Client, id int) int {
	t.Helper()
	report, err := client.TimeReports.Get(ctx, id)
	require.NoError(t, err)
	return len(report.Additions)
}

func describeCountChange(before, after int) string {
	switch {
	case after == before:
		return "kept"
	case after == 0:
		return "cleared"
	default:
		return "changed"
	}
}

func inputAsMap(input TimeReportInput) (map[string]any, error) {
	encoded, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	var body map[string]any
	err = json.Unmarshal(encoded, &body)
	return body, err
}

func putWithExplicitEmptyAdditions(ctx context.Context, client *Client, id int, input TimeReportInput) error {
	body, err := inputAsMap(input)
	if err != nil {
		return err
	}
	body["additions"] = []any{}
	path, err := resourcePath(timeReportsPath, strconv.Itoa(id))
	if err != nil {
		return err
	}
	return client.do(ctx, http.MethodPut, path, nil, body, nil)
}

func postWithExplicitZeroInvoiceableDistance(ctx context.Context, client *Client, input TimeReportInput) (int, error) {
	body, err := inputAsMap(input)
	if err != nil {
		return 0, err
	}
	travelReport, ok := body["travelReport"].(map[string]any)
	if !ok {
		return 0, errors.New("input has no travelReport")
	}
	travelReport["invoiceableDistance"] = 0
	var created createdResource
	err = client.do(ctx, http.MethodPost, timeReportsPath, nil, body, &created)
	return created.ID, err
}

func logTravelDistances(t *testing.T, ctx context.Context, client *Client, label string, id int) {
	t.Helper()
	report, err := client.TimeReports.Get(ctx, id)
	require.NoError(t, err)
	if report.TravelReport == nil {
		t.Logf("%s: report %d has no travelReport", label, id)
		return
	}
	t.Logf("%s: travel distance=%v invoiceableDistance=%v", label, report.TravelReport.Distance, report.TravelReport.InvoiceableDistance)
}

func TestIntegrationTimeReportCreateUpdateDelete(t *testing.T) {
	userID := writeTestUserID(t)
	client, recorder := newRecordingWriteClient(t)
	ctx := integrationContext(t)
	countBefore := countWriteTestWeekReports(t, ctx, client, userID)
	t.Cleanup(func() {
		assert.Equal(t, countBefore, countWriteTestWeekReports(t, context.Background(), client, userID), "week 43 report count changed")
	})

	references := borrowTimeReportReferences(t, ctx, client, userID)
	input := references
	input.Date = NewDate(2026, time.October, 20)
	input.Hours = 0.5
	input.Comment = writeTestComment

	id := createReportWithCleanup(t, ctx, client, input)
	createResponse := recorder.last(http.MethodPost, "/"+timeReportsPath)
	t.Logf("Create response: status %d, body %s", createResponse.statusCode, createResponse.body)

	created, err := client.TimeReports.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, input.Date, created.Date)
	assert.Equal(t, 0.5, created.Hours)
	assert.Equal(t, writeTestComment, created.Comment)
	assert.Equal(t, userID, created.User.ID)

	input.Hours = 1
	input.Comment = writeTestComment + " (updated)"
	require.NoError(t, client.TimeReports.Update(ctx, id, input))

	updated, err := client.TimeReports.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, 1.0, updated.Hours)
	assert.Equal(t, writeTestComment+" (updated)", updated.Comment)

	require.NoError(t, client.TimeReports.Delete(ctx, id))
	assertReportIsGone(t, ctx, client, id)

	t.Run("PUT slice semantics for additions", func(t *testing.T) {
		addition, found := findAdditionToBorrow(t, ctx, client, userID)
		if !found {
			t.Skipf("no usable addition in the latest %d reports with additions (user %d first, then tenant-wide); skipping additions probe", additionSearchLimit, userID)
		}

		probes := []struct {
			description string
			date        Date
			put         func(id int, input TimeReportInput) error
		}{
			{"PUT with additions omitted (nil)", NewDate(2026, time.October, 21), func(id int, input TimeReportInput) error {
				input.Additions = nil
				return client.TimeReports.Update(ctx, id, input)
			}},
			{"PUT with additions as empty slice (omitempty drops it, same wire as nil)", NewDate(2026, time.October, 22), func(id int, input TimeReportInput) error {
				input.Additions = []AdditionInput{}
				return client.TimeReports.Update(ctx, id, input)
			}},
			{"PUT with explicit \"additions\":[] on the wire", NewDate(2026, time.October, 23), func(id int, input TimeReportInput) error {
				input.Additions = nil
				return putWithExplicitEmptyAdditions(ctx, client, id, input)
			}},
		}
		for _, probe := range probes {
			input := references
			input.Date = probe.date
			input.Hours = 0.5
			input.Comment = writeTestComment + " (additions)"
			input.Additions = []AdditionInput{addition}

			id := createReportWithCleanup(t, ctx, client, input)
			before := additionCount(t, ctx, client, id)
			require.Equal(t, 1, before, "report %d should start with one addition", id)

			require.NoError(t, probe.put(id, input), probe.description)
			after := additionCount(t, ctx, client, id)
			t.Logf("%s: additions after = %d (%s)", probe.description, after, describeCountChange(before, after))
			assert.Equal(t, 0, after, "%s should clear additions because PUT replaces the whole report", probe.description)

			require.NoError(t, client.TimeReports.Delete(ctx, id))
			assertReportIsGone(t, ctx, client, id)
		}
	})

	t.Run("zero on invoiceable override", func(t *testing.T) {
		input := references
		input.Date = NewDate(2026, time.October, 24)
		input.Hours = 0.5
		input.Comment = writeTestComment + " (travel)"
		input.TravelReport = &TravelReportInput{Place: "blikk-sdk test", Distance: 10, InvoiceableDistance: 0}

		probes := []struct {
			label  string
			create func() (int, error)
		}{
			{"invoiceableDistance omitted", func() (int, error) { return client.TimeReports.Create(ctx, input) }},
			{"explicit \"invoiceableDistance\":0", func() (int, error) { return postWithExplicitZeroInvoiceableDistance(ctx, client, input) }},
		}
		for _, probe := range probes {
			id := createReportWithCleanupOrSkip(t, client, probe.create)
			logTravelDistances(t, ctx, client, probe.label, id)
			require.NoError(t, client.TimeReports.Delete(ctx, id))
			assertReportIsGone(t, ctx, client, id)
		}
	})
}
