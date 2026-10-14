package blikk

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectCollectionsListEncodesSortingAndDecodesFixture(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Core/ProjectCollections", "projectcollections_list.json"))

	page, err := client.ProjectCollections.List(context.Background(), ListOptions{PageSize: 50, SortBy: ProjectCollectionSortByProjectNumber, SortOrder: SortAscending})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"page":      {"1"},
		"pageSize":  {"50"},
		"sortBy":    {"projectNumber"},
		"sortOrder": {"ascending"},
	}, handler.last().Query)
	require.Len(t, page.Items, 1)
	collection := page.Items[0]
	assert.Equal(t, 1, collection.ID)
	assert.Equal(t, "26 Lorem ipsum", collection.Name)
	assert.Equal(t, "Lorem ipsum", collection.Title)
	assert.Equal(t, "26", collection.ProjectNumber)
	assert.Equal(t, "Lorem ipsum dolor sit amit.", collection.Description)
	assert.Equal(t, &ProjectStatus{ID: 1, Name: "Ongoing"}, collection.Status)
	assert.Equal(t, ProjectCollectionTypeNormal, collection.Type)
	assert.Equal(t, &ProjectCollectionCategory{Name: "Category 1", Color: "Green"}, collection.Category)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, collection.Manager)
	assert.Equal(t, &NameID{ID: 1, Name: "Ipsum"}, collection.Customer)
	assert.Equal(t, "2021-02-01 00:00:00", collection.FinishedDate.String())
}

func TestProjectCollectionsGetDecodesFixture(t *testing.T) {
	client, _ := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Core/ProjectCollections/42", "projectcollection.json"))

	collection, err := client.ProjectCollections.Get(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, 42, collection.ID)
	assert.Equal(t, "42 Lorem ipsum", collection.Name)
	assert.Equal(t, &ProjectStatus{ID: 1, Name: "Active"}, collection.Status)
	assert.Equal(t, &ProjectCollectionCategory{Name: "Category 1", Color: "Red"}, collection.Category)
	assert.Equal(t, &NameID{ID: 1, Name: "ACME Inc."}, collection.Customer)
	assert.Equal(t, "2026-10-13 00:00:00", collection.FinishedDate.String())
}

func TestProjectCollectionTypeMatchesLowercaseValues(t *testing.T) {
	for raw, want := range map[string]ProjectCollectionType{
		`"normal"`:   ProjectCollectionTypeNormal,
		`"absence"`:  ProjectCollectionTypeAbsence,
		`"Internal"`: ProjectCollectionTypeInternal,
	} {
		var got ProjectCollectionType
		require.NoError(t, got.UnmarshalJSON([]byte(raw)))
		assert.Equal(t, want, got)
	}
}

func TestProjectCollectionsNextProjectNumberUnwrapsObject(t *testing.T) {
	client, _ := newRecordingClient(t, routes(t, map[string]http.HandlerFunc{
		"GET /v1/Core/ProjectCollections/123/nextProjectNumber": respondWith(http.StatusOK, `{"nextProjectNumberInCollection":"P001-001"}`),
	}))

	next, err := client.ProjectCollections.NextProjectNumber(context.Background(), 123)

	require.NoError(t, err)
	assert.Equal(t, "P001-001", next)
}
