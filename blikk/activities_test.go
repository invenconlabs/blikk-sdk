package blikk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivitiesListDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "activities_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	page, err := client.Activities.List(context.Background(), ActivityListParams{
		ListOptions: ListOptions{SortBy: ActivitySortByUpdatedDate, SortOrder: SortDescending},
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/Activities", request.Path)
	assert.Equal(t, map[string][]string{"page": {"1"}, "sortBy": {"updatedDate"}, "sortOrder": {"descending"}}, request.Query)

	require.Len(t, page.Items, 1)
	item := page.Items[0]
	assert.Equal(t, 1, item.ID)
	assert.Equal(t, "Lorem ipsum", item.Name)
	assert.True(t, item.IsActive)
	assert.Equal(t, &NameID{ID: 2, Name: "Dolar sit amet"}, item.PreferredTimeCode)
	assert.Equal(t, &NameID{ID: 1, Name: "Lorem ipsum"}, item.PreferredTimeArticle)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, item.CreatedBy)
	assert.Nil(t, item.UpdatedBy)
	assert.Equal(t, "2026-10-03 09:40:57", item.Created.String())
	assert.Nil(t, item.Updated)
}

func TestActivitiesListAllUsesListPath(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "activities_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	var ids []int
	for item, err := range client.Activities.ListAll(context.Background(), ActivityListParams{}) {
		require.NoError(t, err)
		ids = append(ids, item.ID)
	}

	assert.Equal(t, []int{1}, ids)
	assert.Equal(t, "/v1/Admin/Activities", handler.last().Path)
}
