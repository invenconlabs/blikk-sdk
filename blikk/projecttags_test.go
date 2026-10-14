package blikk

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectTagsListEncodesSortingAndDecodesFixture(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Admin/ProjectTags", "projecttags_list.json"))

	page, err := client.ProjectTags.List(context.Background(), ListOptions{Page: 2, SortBy: ProjectTagSortByCreatedDate})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"page":   {"2"},
		"sortBy": {"createdDate"},
	}, handler.last().Query)
	require.Len(t, page.Items, 1)
	tag := page.Items[0]
	assert.Equal(t, 1, tag.ID)
	assert.Equal(t, "Big project", tag.Title)
	assert.Equal(t, "#41b18b", tag.Color)
	assert.Nil(t, tag.UpdatedDate)
	assert.Nil(t, tag.UpdatedBy)
	assert.Equal(t, "2026-09-28 07:16:57", tag.CreatedDate.String())
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, tag.CreatedBy)
}
