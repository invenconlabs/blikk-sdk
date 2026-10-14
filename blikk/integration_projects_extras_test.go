//go:build integration

package blikk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationProjectCollections(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.ProjectCollections.List(ctx, ListOptions{PageSize: 10, SortBy: ProjectCollectionSortByName})
	require.NoError(t, err)
	assert.Equal(t, 1, page.Page)
	if len(page.Items) == 0 {
		t.Skip("account has no project collections")
	}
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)

	collection, err := client.ProjectCollections.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, collection.ID)
	assert.Equal(t, listed.ProjectNumber, collection.ProjectNumber)

	_, err = client.ProjectCollections.NextProjectNumber(ctx, listed.ID)
	require.NoError(t, err)
}

func TestIntegrationProjectCategories(t *testing.T) {
	client := integrationClient(t)

	page, err := client.ProjectCategories.List(integrationContext(t), ListOptions{SortBy: ProjectCategorySortByName})

	require.NoError(t, err)
	assert.Equal(t, 1, page.Page)
	for _, category := range page.Items {
		assert.NotZero(t, category.ID)
		assert.NotEmpty(t, category.Name)
	}
}

func TestIntegrationProjectStatuses(t *testing.T) {
	client := integrationClient(t)

	page, err := client.ProjectStatuses.List(integrationContext(t), ProjectStatusListParams{ListOptions: ListOptions{SortBy: ProjectStatusSortBySortOrder}})

	require.NoError(t, err)
	require.NotEmpty(t, page.Items, "every account has at least the standard project status")
	status := page.Items[0]
	assert.NotZero(t, status.ID)
	assert.NotEmpty(t, status.Name)
	assert.False(t, status.CreatedDate.IsZero())
}

func TestIntegrationProjectTags(t *testing.T) {
	client := integrationClient(t)

	page, err := client.ProjectTags.List(integrationContext(t), ListOptions{SortBy: ProjectTagSortByName})

	require.NoError(t, err)
	assert.Equal(t, 1, page.Page)
	for _, tag := range page.Items {
		assert.NotZero(t, tag.ID)
		assert.False(t, tag.CreatedDate.IsZero())
	}
}

func TestIntegrationProjectComments(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	projectID := findProjectWithComments(t, ctx, client)
	page, err := client.Projects.ListComments(ctx, projectID, ProjectCommentListParams{ListOptions: ListOptions{SortBy: ProjectCommentSortByCreatedDate}})

	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	comment := page.Items[0]
	assert.NotZero(t, comment.ID)
	assert.NotNil(t, comment.CreatedDate)
}

// Comments are rare, and the oldest projects are the likeliest to have some.
func findProjectWithComments(t *testing.T, ctx context.Context, client *Client) int {
	t.Helper()
	projects, err := client.Projects.List(ctx, ProjectListParams{ListOptions: ListOptions{PageSize: 30, SortBy: ProjectSortByCreatedDate}})
	require.NoError(t, err)
	for _, project := range projects.Items {
		comments, err := client.Projects.ListComments(ctx, project.ID, ProjectCommentListParams{ListOptions: ListOptions{PageSize: 1}})
		require.NoError(t, err)
		if comments.TotalItemCount > 0 {
			return project.ID
		}
	}
	t.Skip("none of the 30 oldest projects has comments")
	return 0
}

func TestIntegrationProjectTimeCalculation(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)
	projects, err := client.Projects.List(ctx, ProjectListParams{ListOptions: ListOptions{PageSize: 1}})
	require.NoError(t, err)
	require.NotEmpty(t, projects.Items)

	calculation, err := client.Projects.TimeCalculation(ctx, projects.Items[0].ID)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, calculation.Total, 0.0)
}

func TestIntegrationProjectsNextProjectNumber(t *testing.T) {
	client := integrationClient(t)

	next, err := client.Projects.NextProjectNumber(integrationContext(t))

	require.NoError(t, err)
	assert.NotEmpty(t, next)
}
