package blikk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskBoardsListEncodesEveryFilter(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, "GET", "/v1/Core/TaskBoards", "taskboards_list.json"))

	_, err := client.TaskBoards.List(context.Background(), TaskBoardListParams{
		ListOptions: ListOptions{SortBy: TaskBoardSortByUpdatedDate, SortOrder: SortDescending},
		Query:       "sprint",
		ProjectIDs:  []int{15, 16},
	})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"page":              {"1"},
		"sortBy":            {"updatedDate"},
		"sortOrder":         {"descending"},
		"filter.query":      {"sprint"},
		"filter.projectIds": {"15,16"},
	}, handler.last().Query)
}

func TestTaskBoardsListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, "GET", "/v1/Core/TaskBoards", "taskboards_list.json"))

	page, err := client.TaskBoards.List(context.Background(), TaskBoardListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	board := page.Items[0]
	assert.Equal(t, 1, board.ID)
	assert.Equal(t, "Sprint 2.8", board.Name)
	assert.Equal(t, &ProjectRef{ID: 15, Name: "Utvecklingsprojektet", Number: "26"}, board.Project)
	assert.Equal(t, []TagRef{{ID: 1, Color: "#41b18b"}, {ID: 2, Title: "Bug", Color: "#f2d600"}}, board.Tags)
	assert.Equal(t, []NameID{{ID: 1, Name: "John Doe"}, {ID: 2, Name: "Jane Doe"}}, board.Members)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, board.CreatedByUser)
	assert.Nil(t, board.UpdatedByUser)
	assert.Equal(t, "2026-10-03 09:40:57", board.CreatedDate.String())
	assert.Nil(t, board.UpdatedDate)
	require.Len(t, board.Columns, 3)
	column := board.Columns[1]
	assert.Equal(t, 2, column.ID)
	assert.Equal(t, "Pågår", column.Name)
	assert.Equal(t, 7, column.TaskCount)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, column.CreatedByUser)
	assert.Nil(t, column.UpdatedByUser)
	assert.Equal(t, "2026-10-06 09:40:57", column.CreatedDate.String())
	assert.Nil(t, column.UpdatedDate)
}
