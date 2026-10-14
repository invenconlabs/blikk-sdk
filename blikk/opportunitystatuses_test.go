package blikk

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpportunityStatusesListEncodesEveryFilter(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "opportunitystatuses_list.json")}
	client := newTestClient(t, handler.ServeHTTP)
	moment := DateTimeOf(time.Date(2021, 1, 1, 14, 0, 0, 0, time.UTC))

	_, err := client.OpportunityStatuses.List(context.Background(), OpportunityStatusListParams{
		ListOptions:     ListOptions{SortBy: OpportunityStatusSortBySortOrder, SortOrder: SortAscending},
		CreatedDateFrom: moment,
		CreatedDateTo:   moment,
		UpdatedDateFrom: moment,
		UpdatedDateTo:   moment,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/OpportunityStatuses", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                   {"1"},
		"sortBy":                 {"sortOrder"},
		"sortOrder":              {"ascending"},
		"filter.createdDateFrom": {"2021-01-01 14:00:00"},
		"filter.createdDateTo":   {"2021-01-01 14:00:00"},
		"filter.updatedDateFrom": {"2021-01-01 14:00:00"},
		"filter.updatedDateTo":   {"2021-01-01 14:00:00"},
	}, request.Query)
}

func TestOpportunityStatusesListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixture(t, "opportunitystatuses_list.json"))

	page, err := client.OpportunityStatuses.List(context.Background(), OpportunityStatusListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	status := page.Items[0]
	assert.Equal(t, 1, status.ID)
	assert.Equal(t, "Started", status.Name)
	assert.Equal(t, 1, status.SortOrder)
	assert.Nil(t, status.UpdatedDate)
	assert.Nil(t, status.UpdatedBy)
	assert.Equal(t, "2026-09-28 07:16:57", status.CreatedDate.String())
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, status.CreatedBy)
}
