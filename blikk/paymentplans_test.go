package blikk

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentPlansListEncodesEveryFilter(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "paymentplans_list.json"))
	moment := DateTimeOf(time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC))

	_, err := client.PaymentPlans.List(context.Background(), PaymentPlanListParams{
		ListOptions:           ListOptions{SortBy: PaymentPlanSortByDescription, SortOrder: SortAscending},
		BatchID:               "c5816c0a-6145-4de2-9e7d-ea8187eb3a43",
		CreatedDateFrom:       moment,
		CreatedDateTo:         moment,
		UpdatedDateFrom:       moment,
		UpdatedDateTo:         moment,
		CreatedOrUpdatedSince: moment,
		ProjectID:             1,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Core/Paymentplans", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                         {"1"},
		"sortBy":                       {"description"},
		"sortOrder":                    {"ascending"},
		"filter.batchId":               {"c5816c0a-6145-4de2-9e7d-ea8187eb3a43"},
		"filter.createdDateFrom":       {"2026-09-01 12:00:00"},
		"filter.createdDateTo":         {"2026-09-01 12:00:00"},
		"filter.updatedDateFrom":       {"2026-09-01 12:00:00"},
		"filter.updatedDateTo":         {"2026-09-01 12:00:00"},
		"filter.createdOrUpdatedSince": {"2026-09-01 12:00:00"},
		"filter.projectId":             {"1"},
	}, request.Query)
}

func TestPaymentPlansListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Core/Paymentplans", "paymentplans_list.json"))

	page, err := client.PaymentPlans.List(context.Background(), PaymentPlanListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 2)
	plan := page.Items[0]
	assert.Equal(t, 1, plan.ID)
	assert.Equal(t, "PP01", plan.ArticleNumber)
	assert.Equal(t, "Installation + material", plan.Description)
	assert.Equal(t, 12000.0, plan.Price)
	assert.Equal(t, 1.0, plan.Units)
	assert.Equal(t, "c5816c0a-6145-4de2-9e7d-ea8187eb3a43", plan.BatchID)
	assert.Equal(t, "2026-10-11 00:00:00", plan.Date.String())
	assert.Equal(t, 560, plan.InvoiceID)
	assert.Nil(t, plan.InvoicedDate)
	assert.Zero(t, plan.InvoiceDraftID)
	assert.Equal(t, &NameID{ID: 12, Name: "My article"}, plan.Article)
	assert.Equal(t, ProjectRef{ID: 1, Name: "My project", Number: "123"}, plan.Project)
	assert.Equal(t, "2026-09-28 07:16:57", plan.CreatedDate.String())
	assert.Zero(t, plan.CreatedBy)
	assert.Equal(t, "2026-10-10 04:09:45", plan.UpdatedDate.String())
	assert.Nil(t, plan.UpdatedBy)
	assert.Nil(t, page.Items[1].Article)
	assert.Nil(t, page.Items[1].UpdatedDate)
}

func TestPaymentPlansCreateSendsDocumentedBodyAndDecodesResult(t *testing.T) {
	client, handler := newRecordingClient(t, routes(t, map[string]http.HandlerFunc{
		"POST /v1/Core/Paymentplans": serveFixture(t, "paymentplans_create_result.json"),
	}))

	result, err := client.PaymentPlans.Create(context.Background(), PaymentPlanInput{
		ProjectID:   1,
		Date:        DateTimeOf(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)),
		Amount:      12000,
		Description: "Installation",
		ArticleRows: []PaymentPlanArticleRowInput{{ArticleNumber: "PP01", Description: "Installation + material", Quantity: 1, UnitPrice: 12000}},
	})

	require.NoError(t, err)
	assert.Equal(t, &PaymentPlanCreateResult{BatchID: "c5816c0a-6145-4de2-9e7d-ea8187eb3a43", PaymentPlanIDs: []int{101}}, result)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, map[string]any{
		"projectId":   1.0,
		"date":        "2026-08-01 00:00:00",
		"amount":      12000.0,
		"description": "Installation",
		"articleRows": []any{map[string]any{"articleNumber": "PP01", "description": "Installation + material", "quantity": 1.0, "unitPrice": 12000.0}},
	}, request.jsonBody(t))
}

func TestPaymentPlansCreateMinimalBodyOmitsOptionalFields(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "paymentplans_create_result.json"))

	_, err := client.PaymentPlans.Create(context.Background(), PaymentPlanInput{
		ProjectID:   1,
		ArticleRows: []PaymentPlanArticleRowInput{{}},
	})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"projectId":   1.0,
		"articleRows": []any{map[string]any{}},
	}, handler.last().jsonBody(t))
}

func TestPaymentPlanInputReadsDocumentedRequest(t *testing.T) {
	var input PaymentPlanInput
	require.NoError(t, json.Unmarshal(loadFixture(t, "paymentplans_create_input.json"), &input))

	assert.Equal(t, PaymentPlanInput{
		ProjectID:   1,
		Date:        DateTimeOf(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)),
		ArticleRows: []PaymentPlanArticleRowInput{{ArticleNumber: "PP01", Description: "Installation + material", Quantity: 1, UnitPrice: 12000}},
	}, input)
}
