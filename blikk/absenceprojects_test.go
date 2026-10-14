package blikk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAbsenceProjectsListDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "absenceprojects_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	page, err := client.AbsenceProjects.List(context.Background(), AbsenceProjectListParams{
		ListOptions: ListOptions{SortBy: AbsenceProjectSortByName, SortOrder: SortDescending},
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/AbsenceProjects", request.Path)
	assert.Equal(t, map[string][]string{"page": {"1"}, "sortBy": {"name"}, "sortOrder": {"descending"}}, request.Query)

	require.Len(t, page.Items, 1)
	item := page.Items[0]
	assert.Equal(t, 1, item.ID)
	assert.Equal(t, "Lorem ipsum", item.Name)
	assert.False(t, item.CommentRequiredWhenTimeReporting)
	assert.Equal(t, "040", item.SalaryCodeForMonthlySalary)
	assert.Equal(t, "041", item.SalaryCodeForHourlySalary)
	assert.True(t, item.RefersToSickness)
	assert.True(t, item.IsActive)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, item.CreatedBy)
	assert.Nil(t, item.UpdatedBy)
	assert.Equal(t, "2026-10-03 09:40:57", item.Created.String())
	assert.Nil(t, item.Updated)
}

func TestAbsenceProjectsListAllUsesListPath(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "absenceprojects_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	var ids []int
	for item, err := range client.AbsenceProjects.ListAll(context.Background(), AbsenceProjectListParams{}) {
		require.NoError(t, err)
		ids = append(ids, item.ID)
	}

	assert.Equal(t, []int{1}, ids)
	assert.Equal(t, "/v1/Admin/AbsenceProjects", handler.last().Path)
	assert.Equal(t, []string{"100"}, handler.last().Query["pageSize"])
}
