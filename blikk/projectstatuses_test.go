package blikk

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectStatusesListEncodesFiltersAndDecodesFixture(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Admin/ProjectStatuses", "projectstatuses_list.json"))
	since := DateTimeOf(time.Date(2021, 1, 1, 14, 0, 0, 0, time.UTC))

	page, err := client.ProjectStatuses.List(context.Background(), ProjectStatusListParams{
		ListOptions:     ListOptions{SortBy: ProjectStatusSortBySortOrder},
		CreatedDateFrom: since,
		CreatedDateTo:   since,
		UpdatedDateFrom: since,
		UpdatedDateTo:   since,
	})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"page":                   {"1"},
		"sortBy":                 {"sortOrder"},
		"filter.createdDateFrom": {"2021-01-01 14:00:00"},
		"filter.createdDateTo":   {"2021-01-01 14:00:00"},
		"filter.updatedDateFrom": {"2021-01-01 14:00:00"},
		"filter.updatedDateTo":   {"2021-01-01 14:00:00"},
	}, handler.last().Query)
	require.Len(t, page.Items, 1)
	status := page.Items[0]
	assert.Equal(t, 1, status.ID)
	assert.Equal(t, "Ongoing", status.Name)
	assert.Equal(t, "#e3168a", status.Color)
	assert.Equal(t, 1, status.SortOrder)
	assert.False(t, status.IsCompleted)
	assert.True(t, status.IsStandard)
	assert.Nil(t, status.UpdatedDate)
	assert.Nil(t, status.UpdatedBy)
	assert.Equal(t, "2026-09-28 07:16:57", status.CreatedDate.String())
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, status.CreatedBy)
}
