package blikk

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpportunitiesListEncodesEveryFilter(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "opportunities_list.json")}
	client := newTestClient(t, handler.ServeHTTP)
	moment := DateTimeOf(time.Date(2021, 1, 1, 14, 0, 0, 0, time.UTC))

	_, err := client.Opportunities.List(context.Background(), OpportunityListParams{
		ListOptions:              ListOptions{SortBy: OpportunitySortByTitle, SortOrder: SortAscending},
		CreatedDateFrom:          moment,
		CreatedDateTo:            moment,
		UpdatedDateFrom:          moment,
		UpdatedDateTo:            moment,
		ClosedDateFrom:           moment,
		ClosedDateTo:             moment,
		EstimatedClosingDateFrom: moment,
		EstimatedClosingDateTo:   moment,
		HasOffers:                Ptr(false),
		OpportunityStatusIDs:     []int{1, 2},
		OpportunityTagIDs:        []int{3, 4},
		OpportunityState:         OpportunityStateWon,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Core/Opportunities", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                            {"1"},
		"sortBy":                          {"title"},
		"sortOrder":                       {"ascending"},
		"filter.createdDateFrom":          {"2021-01-01 14:00:00"},
		"filter.createdDateTo":            {"2021-01-01 14:00:00"},
		"filter.updatedDateFrom":          {"2021-01-01 14:00:00"},
		"filter.updatedDateTo":            {"2021-01-01 14:00:00"},
		"filter.closedDateFrom":           {"2021-01-01 14:00:00"},
		"filter.closedDateTo":             {"2021-01-01 14:00:00"},
		"filter.estimatedClosingDateFrom": {"2021-01-01 14:00:00"},
		"filter.estimatedClosingDateTo":   {"2021-01-01 14:00:00"},
		"filter.hasOffers":                {"false"},
		"filter.opportunityStatusIds":     {"1,2"},
		"filter.opportunityTagIds":        {"3,4"},
		"filter.opportunityState":         {"Won"},
	}, request.Query)
}

func TestOpportunitiesListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixture(t, "opportunities_list.json"))

	page, err := client.Opportunities.List(context.Background(), OpportunityListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	item := page.Items[0]
	assert.Equal(t, 1, item.ID)
	assert.Equal(t, "Lorem Ipsum", item.Title)
	assert.InDelta(t, 100000.0, item.Turnover, 0)
	assert.True(t, item.HasOffers)
	assert.Equal(t, NameID{ID: 1, Name: "John Doe"}, item.User)
	assert.Equal(t, &NameID{ID: 1, Name: "Ipsum"}, item.Customer)
	assert.Nil(t, item.Contact)
	assert.Equal(t, OpportunityStateWon, item.OpportunityState)
	assert.Equal(t, "2026-10-13 00:00:00", item.ClosedDate.String())
	assert.Nil(t, item.EstimatedClosingDate)
	assert.Equal(t, &NameID{ID: 1, Name: "Lorem"}, item.Status)
	assert.Equal(t, "2026-10-13 09:40:57", item.CreatedDate.String())
	assert.Equal(t, []Tag{{ID: 1, Title: "Ipsum", Color: "#41b18b"}}, item.Tags)
}

func TestOpportunitiesGetDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "opportunity.json")}
	client := newTestClient(t, handler.ServeHTTP)

	opportunity, err := client.Opportunities.Get(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, "/v1/Core/Opportunities/1", handler.last().Path)
	assert.Equal(t, 1, opportunity.ID)
	assert.Equal(t, "Lorem Ipsum", opportunity.Title)
	assert.Empty(t, opportunity.Description)
	assert.InDelta(t, 100000.0, opportunity.Turnover, 0)
	assert.Zero(t, opportunity.Probability)
	assert.True(t, opportunity.HasOffers)
	assert.Equal(t, NameID{ID: 1, Name: "John Doe"}, opportunity.User)
	assert.Nil(t, opportunity.Contact)
	assert.Nil(t, opportunity.EstimatedClosingDate)
	assert.Empty(t, opportunity.OpportunityState)
	assert.Equal(t, "2026-10-13 00:00:00", opportunity.ClosedDate.String())
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, opportunity.ClosedByUser)
	assert.Nil(t, opportunity.Project)
	assert.Equal(t, &NameID{ID: 1, Name: "Lorem"}, opportunity.Status)
	assert.Equal(t, []Tag{{ID: 1, Title: "Ipsum", Color: "#41b18b"}}, opportunity.Tags)
}

func TestOpportunityStateAcceptsLowercaseDocsSpelling(t *testing.T) {
	var state OpportunityState

	require.NoError(t, state.UnmarshalJSON([]byte(`"lost"`)))

	assert.Equal(t, OpportunityStateLost, state)
}

func TestOpportunitiesCreateSendsFullInput(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusCreated, `{
		"objectName": "opportunity",
		"id": 42,
		"customerId": 1,
		"title": "Lorem",
		"statusId": 2,
		"turnover": 1500.5,
		"estimatedClosingDate": "2026-12-31 00:00:00",
		"description": "Ipsum",
		"userId": 3,
		"createdDate": "2026-10-14 08:30:00",
		"createdById": 7
	}`))
	closing := DateTimeOf(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC))

	created, err := client.Opportunities.Create(context.Background(), OpportunityCreateInput{
		CustomerID:           1,
		Title:                "Lorem",
		StatusID:             2,
		UserID:               3,
		Turnover:             1500.5,
		EstimatedClosingDate: closing,
		Description:          "Ipsum",
	})

	require.NoError(t, err)
	assert.Equal(t, 42, created.ID)
	assert.Equal(t, "Lorem", created.Title)
	assert.Equal(t, &closing, created.EstimatedClosingDate)
	assert.Equal(t, 7, created.CreatedByID)
	assert.Equal(t, DateTimeOf(time.Date(2026, 10, 14, 8, 30, 0, 0, time.UTC)), created.CreatedDate)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/Core/Opportunities", request.Path)
	assert.Equal(t, map[string]any{
		"customerId":           float64(1),
		"title":                "Lorem",
		"statusId":             float64(2),
		"userId":               float64(3),
		"turnover":             1500.5,
		"estimatedClosingDate": "2026-12-31 00:00:00",
		"description":          "Ipsum",
	}, request.jsonBody(t))
}

func TestOpportunitiesCreateOmitsUnsetOptionalFields(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusCreated, `{"id": 42}`))

	_, err := client.Opportunities.Create(context.Background(), OpportunityCreateInput{CustomerID: 1, Title: "Lorem", StatusID: 2, UserID: 3})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"customerId": float64(1),
		"title":      "Lorem",
		"statusId":   float64(2),
		"userId":     float64(3),
	}, handler.last().jsonBody(t))
}

func TestOpportunitiesUpdateSendsFullInput(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ``))
	closing := DateTimeOf(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC))

	err := client.Opportunities.Update(context.Background(), 9, OpportunityUpdateInput{
		Title:                "Lorem",
		UserID:               3,
		Probability:          75,
		Description:          "Ipsum",
		CustomerID:           1,
		ContactID:            4,
		Turnover:             1500.5,
		EstimatedClosingDate: closing,
		TagIDs:               []int{5, 6},
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodPut, request.Method)
	assert.Equal(t, "/v1/Core/Opportunities/9", request.Path)
	assert.Equal(t, map[string]any{
		"title":                "Lorem",
		"userId":               float64(3),
		"probability":          float64(75),
		"description":          "Ipsum",
		"customerId":           float64(1),
		"contactId":            float64(4),
		"turnover":             1500.5,
		"estimatedClosingDate": "2026-12-31 00:00:00",
		"tagIds":               []any{float64(5), float64(6)},
	}, request.jsonBody(t))
}

func TestOpportunitiesUpdateSendsOnlyRequiredFieldsForMinimalInput(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ``))

	err := client.Opportunities.Update(context.Background(), 9, OpportunityUpdateInput{Title: "Lorem", UserID: 3, Probability: 0})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"title":       "Lorem",
		"userId":      float64(3),
		"probability": float64(0),
	}, handler.last().jsonBody(t))
}
