package blikk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeCodesListDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "timecodes_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	page, err := client.TimeCodes.List(context.Background(), TimeCodeListParams{
		ListOptions: ListOptions{SortBy: TimeCodeSortByUpdatedDate, SortOrder: SortAscending},
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/TimeCodes", request.Path)
	assert.Equal(t, map[string][]string{"page": {"1"}, "sortBy": {"updatedDate"}, "sortOrder": {"ascending"}}, request.Query)

	require.Len(t, page.Items, 1)
	item := page.Items[0]
	assert.Equal(t, 1, item.ID)
	assert.Equal(t, "Lorem ipsum", item.Name)
	assert.Equal(t, "Lorem ipsum dolar sit amet", item.Description)
	assert.True(t, item.IsActive)
	assert.False(t, item.IsStandard)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, item.CreatedBy)
	assert.Nil(t, item.UpdatedBy)
	assert.Equal(t, "2026-10-03 09:40:57", item.Created.String())
	assert.Nil(t, item.Updated)
}

func TestTimeCodesListAllUsesListPath(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "timecodes_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	var ids []int
	for item, err := range client.TimeCodes.ListAll(context.Background(), TimeCodeListParams{}) {
		require.NoError(t, err)
		ids = append(ids, item.ID)
	}

	assert.Equal(t, []int{1}, ids)
	assert.Equal(t, "/v1/Admin/TimeCodes", handler.last().Path)
}
