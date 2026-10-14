package blikk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpportunityTagsListEncodesSortAndPaging(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "opportunitytags_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	_, err := client.OpportunityTags.List(context.Background(), OpportunityTagListParams{
		ListOptions: ListOptions{Page: 2, PageSize: 10, SortBy: OpportunityTagSortByName, SortOrder: SortDescending},
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/OpportunityTags", request.Path)
	assert.Equal(t, map[string][]string{
		"page":      {"2"},
		"pageSize":  {"10"},
		"sortBy":    {"name"},
		"sortOrder": {"descending"},
	}, request.Query)
}

func TestOpportunityTagsListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixture(t, "opportunitytags_list.json"))

	page, err := client.OpportunityTags.List(context.Background(), OpportunityTagListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	tag := page.Items[0]
	assert.Equal(t, 1, tag.ID)
	assert.Equal(t, "High value opportunity", tag.Title)
	assert.Equal(t, "#41b18b", tag.Color)
	assert.Equal(t, "2026-09-28 07:16:57", tag.CreatedDate.String())
	assert.Nil(t, tag.UpdatedDate)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, tag.CreatedBy)
}
