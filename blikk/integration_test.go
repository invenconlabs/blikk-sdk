//go:build integration

package blikk

import (
	"context"
	"os"
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

	page, err := client.Users.List(ctx, UserListParams{ListOptions: ListOptions{PageSize: 10, SortBy: "name"}})
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

	page, err := client.Projects.List(ctx, ProjectListParams{ListOptions: ListOptions{PageSize: 10, SortBy: "updatedDate", SortOrder: SortDescending}})
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
