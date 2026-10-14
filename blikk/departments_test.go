package blikk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDepartmentsListDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "departments_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	page, err := client.Departments.List(context.Background(), DepartmentListParams{
		ListOptions: ListOptions{SortBy: DepartmentSortByName},
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/Departments", request.Path)
	assert.Equal(t, map[string][]string{"page": {"1"}, "sortBy": {"name"}}, request.Query)

	require.Len(t, page.Items, 1)
	item := page.Items[0]
	assert.Equal(t, 1, item.ID)
	assert.Equal(t, "Lorem ipsum", item.Name)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, item.CreatedBy)
	assert.Equal(t, "2026-10-03 09:40:57", item.Created.String())
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, item.UpdatedBy)
	assert.Equal(t, "2026-10-08 09:40:57", item.Updated.String())
}

func TestDepartmentsListAllUsesListPath(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "departments_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	var ids []int
	for item, err := range client.Departments.ListAll(context.Background(), DepartmentListParams{}) {
		require.NoError(t, err)
		ids = append(ids, item.ID)
	}

	assert.Equal(t, []int{1}, ids)
	assert.Equal(t, "/v1/Admin/Departments", handler.last().Path)
}
