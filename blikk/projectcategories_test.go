package blikk

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectCategoriesListEncodesSortingAndDecodesFixture(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Admin/ProjectCategories", "projectcategories_list.json"))

	page, err := client.ProjectCategories.List(context.Background(), ListOptions{SortBy: ProjectCategorySortByName, SortOrder: SortDescending})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"page":      {"1"},
		"sortBy":    {"name"},
		"sortOrder": {"descending"},
	}, handler.last().Query)
	assert.Equal(t, []ProjectCategoryListItem{{ID: 1, Name: "Lorem ipsum", Color: "#e3168a"}}, page.Items)
}

func TestProjectCategoriesListAllWalksPages(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Admin/ProjectCategories", "projectcategories_list.json"))

	var names []string
	for category, err := range client.ProjectCategories.ListAll(context.Background(), ListOptions{}) {
		require.NoError(t, err)
		names = append(names, category.Name)
	}

	assert.Equal(t, []string{"Lorem ipsum"}, names)
	assert.Equal(t, []string{"100"}, handler.last().Query["pageSize"])
}
