package blikk

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSupplierInvoicesListEncodesEveryFilter(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "supplierinvoices_list.json"))
	moment := DateTimeOf(time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC))

	_, err := client.SupplierInvoices.List(context.Background(), SupplierInvoiceListParams{
		ListOptions:           ListOptions{SortBy: SupplierInvoiceSortByInvoiceCostAddition, SortOrder: SortDescending},
		Query:                 "solar",
		InvoiceDateFrom:       DateTimeOf(time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)),
		InvoiceDateTo:         DateTimeOf(time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)),
		CreatedDateFrom:       moment,
		CreatedDateTo:         moment,
		UpdatedDateFrom:       moment,
		UpdatedDateTo:         moment,
		CreatedOrUpdatedSince: moment,
		ProjectID:             444,
		SupplierID:            44,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/core/supplierinvoices", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                         {"1"},
		"sortBy":                       {"invoiceCostAddition"},
		"sortOrder":                    {"descending"},
		"filter.query":                 {"solar"},
		"filter.invoiceDateFrom":       {"2026-06-01 00:00:00"},
		"filter.invoiceDateTo":         {"2026-06-30 00:00:00"},
		"filter.createdDateFrom":       {"2026-09-01 12:00:00"},
		"filter.createdDateTo":         {"2026-09-01 12:00:00"},
		"filter.updatedDateFrom":       {"2026-09-01 12:00:00"},
		"filter.updatedDateTo":         {"2026-09-01 12:00:00"},
		"filter.createdOrUpdatedSince": {"2026-09-01 12:00:00"},
		"filter.projectId":             {"444"},
		"filter.supplierId":            {"44"},
	}, request.Query)
}

func TestSupplierInvoicesListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, http.MethodGet, "/v1/core/supplierinvoices", "supplierinvoices_list.json"))

	page, err := client.SupplierInvoices.List(context.Background(), SupplierInvoiceListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	invoice := page.Items[0]
	assert.Equal(t, 1, invoice.ID)
	assert.Equal(t, "11232", invoice.InvoiceNumber)
	assert.Equal(t, NewDate(2026, time.September, 30), *invoice.InvoiceDate)
	assert.Equal(t, "Construction time spent", invoice.Title)
	assert.Equal(t, "A brief description.", invoice.Description)
	assert.Equal(t, 2.0, invoice.ExternalHours)
	assert.Equal(t, 1100.0, invoice.TotalAmount)
	assert.Equal(t, 100.0, invoice.InvoiceCostAddition)
	assert.Equal(t, 1000.0, invoice.Cost)
	assert.Equal(t, "1234567890", invoice.OCRNumber)
	assert.True(t, invoice.ToBeInvoiced)
	assert.Zero(t, invoice.InvoiceID)
	assert.Nil(t, invoice.InvoicedDate)
	assert.Zero(t, invoice.InvoiceDraftID)
	assert.Equal(t, &NameID{ID: 44, Name: "Builders Inc."}, invoice.Supplier)
	assert.Equal(t, ProjectRef{ID: 123, Name: "My project", Number: "213"}, invoice.Project)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, invoice.CreatedBy)
	assert.Equal(t, &NameID{ID: 2, Name: "Jane Doe"}, invoice.UpdatedBy)
	assert.Equal(t, "2026-10-01 07:16:57", invoice.CreatedDate.String())
	assert.Equal(t, "2026-10-12 04:09:45", invoice.UpdatedDate.String())
}

func TestSupplierInvoicesGetDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, http.MethodGet, "/v1/core/supplierinvoices/42", "supplierinvoice.json"))

	invoice, err := client.SupplierInvoices.Get(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, 42, invoice.ID)
	assert.Equal(t, "Supplier invoice 1", invoice.Title)
	assert.Equal(t, "Description of supplier invoice 1", invoice.Description)
	assert.Equal(t, &NameID{ID: 1, Name: "Supplier 1"}, invoice.Supplier)
	assert.Equal(t, "123", invoice.InvoiceNumber)
	assert.Equal(t, NewDate(2026, time.August, 13), *invoice.InvoiceDate)
	assert.Equal(t, 200.0, invoice.ProjectCost)
	assert.Equal(t, 100.0, invoice.Addition)
	assert.Equal(t, 300.0, invoice.TotalAmount)
	assert.True(t, invoice.ToBeInvoiced)
	assert.Equal(t, 1, invoice.InvoiceID)
	assert.Equal(t, ProjectRef{ID: 1, Name: "Project 1", Number: "123"}, invoice.Project)
	assert.Equal(t, &NameID{ID: 1, Name: "User 1"}, invoice.CreatedBy)
	assert.Equal(t, &NameID{ID: 2, Name: "User 2"}, invoice.UpdatedBy)
	assert.Equal(t, "2026-07-13 09:40:57", invoice.CreatedDate.String())
	assert.Equal(t, "2026-08-13 09:40:57", invoice.UpdatedDate.String())
	assert.Equal(t, "2026-09-13 09:40:57", invoice.InvoicedDate.String())
	assert.Equal(t, &CostCenterRef{ID: 1, Name: "Cost center 1"}, invoice.CostCenter)
	assert.Equal(t, &NameID{ID: 1, Name: "Project 1"}, invoice.BookingProject)
}

func TestSupplierInvoicesAddUnhandledSendsFullBody(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusCreated, `{"Id":1451}`))

	id, err := client.SupplierInvoices.AddUnhandled(context.Background(), SupplierInvoiceUnhandledInput{
		Title:                "Solar Sverige",
		InvoiceDate:          DateTimeOf(time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)),
		Amount:               805.42,
		SourceSystemID:       "erp-113",
		Description:          "Cables",
		SupplierName:         "Solar Sverige Aktiebolag",
		BookingProjectNumber: "1269",
		InvoiceNumber:        "2511980411",
	})

	require.NoError(t, err)
	assert.Equal(t, 1451, id)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/core/supplierinvoices/addunhandled", request.Path)
	assert.Equal(t, map[string]any{
		"title":                "Solar Sverige",
		"invoiceDate":          "2026-06-30 00:00:00",
		"amount":               805.42,
		"sourceSystemId":       "erp-113",
		"description":          "Cables",
		"supplierName":         "Solar Sverige Aktiebolag",
		"bookingProjectNumber": "1269",
		"invoiceNumber":        "2511980411",
	}, request.jsonBody(t))
}

func TestSupplierInvoicesAddUnhandledMinimalBodyOmitsOptionalFields(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusCreated, `{"id":1451}`))

	id, err := client.SupplierInvoices.AddUnhandled(context.Background(), SupplierInvoiceUnhandledInput{
		Title:          "Solar Sverige",
		InvoiceDate:    DateTimeOf(time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)),
		Amount:         805.42,
		SourceSystemID: "erp-113",
	})

	require.NoError(t, err)
	assert.Equal(t, 1451, id)
	assert.Equal(t, map[string]any{
		"title":          "Solar Sverige",
		"invoiceDate":    "2026-06-30 00:00:00",
		"amount":         805.42,
		"sourceSystemId": "erp-113",
	}, handler.last().jsonBody(t))
}

func TestSupplierInvoicesUpdatePatchesFullBody(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusNoContent, ""))

	err := client.SupplierInvoices.Update(context.Background(), 42, SupplierInvoiceUpdateInput{
		Title:            "Supplier Invoice 2",
		Description:      "Description of supplier invoice 2",
		SupplierID:       2,
		InvoiceNumber:    "321",
		InvoiceDate:      DateTimeOf(time.Date(2026, 8, 13, 0, 0, 0, 0, time.UTC)),
		ProjectCost:      100,
		Addition:         100,
		ToBeInvoiced:     true,
		ProjectID:        6,
		CostCenterID:     7,
		BookingProjectID: 8,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodPatch, request.Method)
	assert.Equal(t, "/v1/core/supplierinvoices/42", request.Path)
	assert.Empty(t, request.Query)
	assert.Equal(t, map[string]any{
		"title":            "Supplier Invoice 2",
		"description":      "Description of supplier invoice 2",
		"supplierId":       2.0,
		"invoiceNumber":    "321",
		"invoiceDate":      "2026-08-13 00:00:00",
		"projectCost":      100.0,
		"addition":         100.0,
		"toBeInvoiced":     true,
		"projectId":        6.0,
		"costCenterId":     7.0,
		"bookingProjectId": 8.0,
	}, request.jsonBody(t))
}

func TestSupplierInvoicesUpdateMinimalBodyKeepsRequiredFields(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusNoContent, ""))

	err := client.SupplierInvoices.Update(context.Background(), 42, SupplierInvoiceUpdateInput{Title: "Supplier Invoice 2", ProjectID: 6})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"title":        "Supplier Invoice 2",
		"toBeInvoiced": false,
		"projectId":    6.0,
	}, handler.last().jsonBody(t))
}

func TestSupplierInvoicesDelete(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusNoContent, ""))

	err := client.SupplierInvoices.Delete(context.Background(), 42)

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodDelete, request.Method)
	assert.Equal(t, "/v1/core/supplierinvoices/42", request.Path)
}
