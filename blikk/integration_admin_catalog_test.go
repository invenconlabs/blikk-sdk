//go:build integration

package blikk

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationPlanning(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)
	from, to := planningWeek()

	users, err := client.Users.List(ctx, UserListParams{ListOptions: ListOptions{PageSize: 1, SortBy: UserSortByName}})
	require.NoError(t, err)
	require.NotEmpty(t, users.Items)
	userID := users.Items[0].ID

	projects, err := client.Planning.ListProjectsWithPlanning(ctx, PlanningProjectListParams{FromDate: from, ToDate: to})
	skipIfPlanningForbidden(t, err)
	require.NoError(t, err)
	for _, project := range projects.Items {
		assert.NotZero(t, project.ID)
	}

	planningUsers, err := client.Planning.ListUsersWithPlanning(ctx, PlanningUserListParams{FromDate: from, ToDate: to, ExcludeDeleted: Ptr(false)})
	require.NoError(t, err)
	for _, user := range planningUsers.Items {
		assert.NotZero(t, user.ID)
	}

	projectSummaries, err := client.Planning.ListProjectSummariesForUser(ctx, PlanningProjectSummaryListParams{UserID: userID, FromDate: from, ToDate: to, GroupBy: PlanningGroupByDay})
	require.NoError(t, err)
	for _, summary := range projectSummaries.Items {
		assert.NotZero(t, summary.ProjectID)
		require.NotNil(t, summary.FromDate)
		assert.False(t, summary.FromDate.Time().Before(from.Time()), "summary starts %s before %s", summary.FromDate, from)
	}

	if len(projects.Items) == 0 {
		t.Skip("no projects with resource planning in the coming week")
	}
	userSummaries, err := client.Planning.ListUserSummariesForProject(ctx, PlanningUserSummaryListParams{ProjectID: projects.Items[0].ID, FromDate: from, ToDate: to})
	require.NoError(t, err)
	require.NotEmpty(t, userSummaries.Items)
	assert.NotZero(t, userSummaries.Items[0].UserID)
	assert.NotEmpty(t, userSummaries.Items[0].UserName)
}

func planningWeek() (Date, Date) {
	today := time.Now()
	return DateOf(today), DateOf(today.AddDate(0, 0, 6))
}

func skipIfPlanningForbidden(t *testing.T, err error) {
	t.Helper()
	if hasStatus(err, http.StatusForbidden) {
		t.Skipf("the integration app lacks the Planning_Read permission: %v", err)
	}
}

func TestIntegrationVacationRequests(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.VacationRequests.List(ctx, VacationRequestListParams{
		ListOptions:    ListOptions{PageSize: 10, SortBy: VacationRequestSortByCreatedDate, SortOrder: SortDescending},
		VacationStatus: VacationStatusAttested,
	})

	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("no attested vacation requests")
	}
	request := page.Items[0]
	assert.NotZero(t, request.ID)
	assert.NotZero(t, request.User.ID)
	assert.NotZero(t, request.AbsenceProject.ID)
	assert.False(t, request.FromDate.IsZero())
	assert.False(t, request.CreatedDate.IsZero())
	assert.NotNil(t, request.AttestedDate)
}

func TestIntegrationArticles(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.Articles.List(ctx, ArticleListParams{
		ListOptions: ListOptions{PageSize: 10, SortBy: ArticleSortByName},
		Type:        ArticleTypeAll,
	})
	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)
	assert.NotEmpty(t, listed.Name)

	article, err := client.Articles.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, article.ID)
	assert.Equal(t, listed.Name, article.Name)
	assert.Equal(t, listed.ArticleNumber, article.ArticleNumber)
	assert.NotEmpty(t, article.Type)
	assert.NotZero(t, article.VatCode.ID)
}

func TestIntegrationArticlesFilterByType(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.Articles.List(ctx, ArticleListParams{ListOptions: ListOptions{PageSize: 1}, Type: ArticleTypeTime})
	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("no time articles")
	}

	article, err := client.Articles.Get(ctx, page.Items[0].ID)
	require.NoError(t, err)
	assert.Equal(t, ArticleTypeTime, article.Type)
}

func TestIntegrationBookingProjects(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.BookingProjects.List(ctx, BookingProjectListParams{
		ListOptions: ListOptions{PageSize: 10, SortBy: BookingProjectSortByCreatedDate, SortOrder: SortDescending},
		IsActive:    Ptr(true),
	})
	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("no active booking projects")
	}
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)
	assert.NotEmpty(t, listed.BookingNumber)
	assert.True(t, listed.IsActive)
	assert.False(t, listed.Created.IsZero())

	project, err := client.BookingProjects.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed, *project)

	byNumber, err := client.BookingProjects.GetByBookingNumber(ctx, listed.BookingNumber)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, byNumber.ID)
}

func TestIntegrationBookingProjectByUnknownNumberIsNotFound(t *testing.T) {
	client := integrationClient(t)

	_, err := client.BookingProjects.GetByBookingNumber(integrationContext(t), "no-such-booking-number-999999")

	assert.True(t, IsNotFound(err), "expected not found, got %v", err)
}
