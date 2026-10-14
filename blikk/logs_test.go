package blikk

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogsListEncodesEveryFilter(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Core/Logs", "logs_list.json"))

	_, err := client.Logs.List(context.Background(), LogListParams{
		ListOptions:     ListOptions{SortBy: LogSortByCreatedDate, SortOrder: SortDescending},
		LogType:         LogTypeProjectChangedStatus,
		CreatedDateFrom: DateTimeOf(time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)),
		CreatedDateTo:   DateTimeOf(time.Date(2026, 9, 30, 23, 59, 59, 0, time.UTC)),
		ProjectID:       24,
	})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"page":                   {"1"},
		"sortBy":                 {"createdDate"},
		"sortOrder":              {"descending"},
		"filter.logType":         {"Project_Changed_Status"},
		"filter.createdDateFrom": {"2026-09-01 00:00:00"},
		"filter.createdDateTo":   {"2026-09-30 23:59:59"},
		"filter.projectId":       {"24"},
	}, handler.last().Query)
}

func TestLogsListDecodesFixture(t *testing.T) {
	client, _ := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Core/Logs", "logs_list.json"))

	page, err := client.Logs.List(context.Background(), LogListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 2)
	statusChange := page.Items[0]
	assert.Equal(t, 1, statusChange.ID)
	assert.Equal(t, "Changed status to Finished", statusChange.Text)
	assert.Empty(t, statusChange.Note)
	assert.Equal(t, LogTypeProjectChangedStatus, statusChange.LogType)
	assert.Equal(t, LogEntity{EntityType: EntityTypeProject, ID: 24, Name: "My project"}, statusChange.Entity)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, statusChange.CreatedBy)
	assert.Equal(t, "2026-10-01 07:16:57", statusChange.CreatedDate.String())
	mention := page.Items[1]
	assert.Empty(t, mention.Text)
	assert.Equal(t, "@johndoe you're up!", mention.Note)
}
