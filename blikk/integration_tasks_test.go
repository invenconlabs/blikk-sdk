//go:build integration

package blikk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationTaskBoards(t *testing.T) {
	client := integrationClient(t)

	page, err := client.TaskBoards.List(integrationContext(t), TaskBoardListParams{ListOptions: ListOptions{PageSize: 10, SortBy: TaskBoardSortByName}})

	require.NoError(t, err)
	assert.Equal(t, 1, page.Page)
	if len(page.Items) == 0 {
		t.Skip("no task boards in the tenant")
	}
	board := page.Items[0]
	assert.NotZero(t, board.ID)
	assert.NotNil(t, board.CreatedDate)
	for _, column := range board.Columns {
		assert.NotZero(t, column.ID)
	}
}

func TestIntegrationTasks(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.Tasks.List(ctx, TaskListParams{ListOptions: ListOptions{PageSize: 10, SortBy: TaskSortByCreatedDate, SortOrder: SortDescending}})

	require.NoError(t, err)
	assert.Equal(t, 1, page.Page)
	if len(page.Items) == 0 {
		t.Skip("no tasks in the tenant")
	}
	task := page.Items[0]
	assert.NotZero(t, task.ID)
	assert.NotNil(t, task.CreatedDate)

	_, err = client.Tasks.ListComments(ctx, task.ID, ListOptions{PageSize: 5})
	require.NoError(t, err)

	_, err = client.Tasks.ListChecklists(ctx, task.ID, ListOptions{PageSize: 5})
	require.NoError(t, err)

	_, err = client.Tasks.ListAttachments(ctx, task.ID, ListOptions{PageSize: 5})
	require.NoError(t, err)
}

func TestIntegrationTaskSubresourceOfMissingTask(t *testing.T) {
	client := integrationClient(t)

	_, err := client.Tasks.ListComments(integrationContext(t), 999999999, ListOptions{})

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 400, apiErr.StatusCode, "the API answers a missing task with 400, not 404")
	assert.True(t, IsNotFound(err), "a 400 saying the task could not be found counts as not found: %v", err)
}
