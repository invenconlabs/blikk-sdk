package blikk

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaterialReportsListEncodesEveryFilter(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "materialreports_list.json"))
	moment := DateTimeOf(time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC))

	_, err := client.MaterialReports.List(context.Background(), MaterialReportListParams{
		ListOptions:           ListOptions{SortBy: MaterialReportSortBySortOrder, SortOrder: SortAscending},
		CreatedDateFrom:       moment,
		CreatedDateTo:         moment,
		UpdatedDateFrom:       moment,
		UpdatedDateTo:         moment,
		CreatedOrUpdatedSince: moment,
		Query:                 "cable",
		OnlyText:              true,
		OnlyInvoiced:          true,
		ProjectID:             407,
		TimeReportID:          270,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Core/MaterialReports", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                         {"1"},
		"sortBy":                       {"sortOrder"},
		"sortOrder":                    {"ascending"},
		"filter.createdDateFrom":       {"2026-09-01 12:00:00"},
		"filter.createdDateTo":         {"2026-09-01 12:00:00"},
		"filter.updatedDateFrom":       {"2026-09-01 12:00:00"},
		"filter.updatedDateTo":         {"2026-09-01 12:00:00"},
		"filter.createdOrUpdatedSince": {"2026-09-01 12:00:00"},
		"filter.query":                 {"cable"},
		"filter.onlyText":              {"true"},
		"filter.onlyInvoiced":          {"true"},
		"filter.projectId":             {"407"},
		"filter.timeReportId":          {"270"},
	}, request.Query)
}

func TestMaterialReportsListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Core/MaterialReports", "materialreports_list.json"))

	page, err := client.MaterialReports.List(context.Background(), MaterialReportListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	assertDocumentedMaterialReport(t, page.Items[0])
}

func TestMaterialReportsGetDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Core/MaterialReports/1", "materialreport.json"))

	report, err := client.MaterialReports.Get(context.Background(), 1)

	require.NoError(t, err)
	assertDocumentedMaterialReport(t, *report)
}

func assertDocumentedMaterialReport(t *testing.T, report MaterialReport) {
	t.Helper()
	assert.Equal(t, 1, report.ID)
	assert.Zero(t, report.SortOrder)
	assert.Equal(t, "PP01", report.ArticleNumber)
	assert.Equal(t, "My article", report.Title)
	assert.Equal(t, 1.0, report.ConsumedUnits)
	assert.Equal(t, 2.0, report.InvoiceableUnits)
	assert.Equal(t, "pcs", report.UnitText)
	assert.Equal(t, 200.0, report.CostPerUnit)
	assert.Equal(t, 250.0, report.PricePerUnit)
	assert.Equal(t, 10.0, report.DiscountPercentage)
	assert.Equal(t, 25.0, report.VatPercentage)
	assert.Equal(t, 450.0, report.TotalAmount)
	assert.True(t, report.ToBeInvoiced)
	assert.Equal(t, 560, report.InvoiceID)
	assert.Nil(t, report.InvoicedDate)
	assert.Zero(t, report.InvoiceDraftID)
	assert.Zero(t, report.TimeReportID)
	assert.False(t, report.IsTextOnly)
	assert.Equal(t, &NameID{ID: 12, Name: "My article"}, report.Article)
	assert.Equal(t, &CostCenterRef{ID: 1, Name: "IT Department"}, report.CostCenter)
	assert.Equal(t, &NameID{ID: 1, Name: "BP 1"}, report.BookingProject)
	assert.Equal(t, ProjectRef{ID: 1, Name: "My project", Number: "123"}, report.Project)
	assert.Equal(t, "2026-09-28 07:16:57", report.CreatedDate.String())
	assert.Equal(t, NameID{ID: 1, Name: "John Doe"}, report.CreatedBy)
	assert.Equal(t, "2026-10-10 04:09:45", report.UpdatedDate.String())
	assert.Nil(t, report.UpdatedBy)
}

func TestMaterialReportsCreateSendsFullBody(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"id":8}`))

	id, err := client.MaterialReports.Create(context.Background(), MaterialReportCreateInput{
		ArticleID:        12,
		ProjectID:        1,
		Units:            1,
		InvoiceableUnits: 2,
		Title:            "My article",
		UnitText:         "pcs",
		Cost:             200,
		UnitPrice:        250,
		Discount:         10,
		Created:          DateTimeOf(time.Date(2026, 9, 28, 7, 16, 57, 0, time.UTC)),
	})

	require.NoError(t, err)
	assert.Equal(t, 8, id)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/Core/MaterialReports", request.Path)
	assert.Equal(t, map[string]any{
		"articleId":       12.0,
		"workOrderId":     1.0,
		"units":           1.0,
		"invoicableUnits": 2.0,
		"title":           "My article",
		"unitText":        "pcs",
		"cost":            200.0,
		"unitPrice":       250.0,
		"discount":        10.0,
		"created":         "2026-09-28 07:16:57",
	}, request.jsonBody(t))
}

func TestMaterialReportsCreateMinimalBodyKeepsRequiredFields(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"id":8}`))

	_, err := client.MaterialReports.Create(context.Background(), MaterialReportCreateInput{ArticleID: 12, ProjectID: 1})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"articleId":       12.0,
		"workOrderId":     1.0,
		"units":           0.0,
		"invoicableUnits": 0.0,
	}, handler.last().jsonBody(t))
}

func TestMaterialReportsUpdateSendsFullBody(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	err := client.MaterialReports.Update(context.Background(), 1, MaterialReportUpdateInput{
		ArticleNumber:    "PP02",
		Units:            2,
		InvoiceableUnits: 3,
		Title:            "Updated material report title",
		UnitText:         "Updated unit text",
		Cost:             100,
		UnitPrice:        300,
		Discount:         20,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodPut, request.Method)
	assert.Equal(t, "/v1/Core/MaterialReports/1", request.Path)
	assert.Equal(t, map[string]any{
		"articleNumber":   "PP02",
		"units":           2.0,
		"invoicableUnits": 3.0,
		"title":           "Updated material report title",
		"unitText":        "Updated unit text",
		"cost":            100.0,
		"unitPrice":       300.0,
		"discount":        20.0,
	}, request.jsonBody(t))
}

func TestMaterialReportsUpdateMinimalBodyIsEmpty(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	err := client.MaterialReports.Update(context.Background(), 1, MaterialReportUpdateInput{Title: "Renamed"})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"title": "Renamed"}, handler.last().jsonBody(t))
}

func TestMaterialReportsDelete(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	err := client.MaterialReports.Delete(context.Background(), 1)

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodDelete, request.Method)
	assert.Equal(t, "/v1/Core/MaterialReports/1", request.Path)
}
