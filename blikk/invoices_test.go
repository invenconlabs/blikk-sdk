package blikk

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoicesListEncodesEveryFilter(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "invoices_list.json"))

	_, err := client.Invoices.List(context.Background(), InvoiceListParams{
		ListOptions:  ListOptions{PageSize: 10, SortBy: InvoiceSortBySentToEconomySystem, SortOrder: SortAscending},
		CustomerID:   92,
		CreatedSince: DateTimeOf(time.Date(2026, 10, 1, 8, 30, 0, 0, time.UTC)),
		SinceID:      41,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodGet, request.Method)
	assert.Equal(t, "/v1/core/invoices", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                {"1"},
		"pageSize":            {"10"},
		"sortBy":              {"sentToEconomySystem"},
		"sortOrder":           {"ascending"},
		"filter.customerId":   {"92"},
		"filter.createdSince": {"2026-10-01 08:30:00"},
		"filter.sinceId":      {"41"},
	}, request.Query)
}

func TestInvoicesListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, http.MethodGet, "/v1/core/invoices", "invoices_list.json"))

	page, err := client.Invoices.List(context.Background(), InvoiceListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	invoice := page.Items[0]
	assert.Equal(t, 42, invoice.ID)
	assert.Equal(t, "2026-10-13 09:40:57", invoice.CreatedDate.String())
	assert.Equal(t, "2026-10-13 09:40:57", invoice.InvoiceDate.String())
	assert.Equal(t, &NameID{ID: 123, Name: "Toyota"}, invoice.Customer)
	assert.Equal(t, "2026-10-13", invoice.FromDate.String())
	assert.Equal(t, "2026-10-13", invoice.ToDate.String())
	assert.Equal(t, 1200.50, invoice.Sum)
	assert.Equal(t, "2026-10-13 09:40:57", invoice.SentToEconomySystem.String())
}

func TestInvoicesGetDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, http.MethodGet, "/v1/core/invoices/42", "invoice.json"))

	invoice, err := client.Invoices.Get(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, 42, invoice.ID)
	assert.False(t, invoice.Archived)
	assert.Equal(t, "BRF-123", invoice.Cadastral)
	assert.Equal(t, "2026-10-13 09:40:57", invoice.CreatedDate.String())
	assert.Equal(t, &NameID{ID: 92, Name: "Söderåsens BRF"}, invoice.Customer)
	assert.Equal(t, "ref-123", invoice.CustomerReferenceMarking)
	assert.Equal(t, "88-88", invoice.EconomySystemInvoiceNumber)
	assert.Equal(t, "EFH12345", invoice.ElectronicInvoiceReferenceCode)
	assert.Equal(t, "2026-10-13 09:40:57", invoice.FromDate.String())
	assert.Equal(t, "2026-10-13 09:40:57", invoice.InvoiceDate.String())
	assert.Equal(t, "2026-10-13 09:40:57", invoice.ToDate.String())
	assert.Equal(t, 1001, invoice.InvoiceNumber)
	assert.Equal(t, "Rent", invoice.InvoiceText)
	assert.Equal(t, "Mr. Bean", invoice.InvoiceReference)
	assert.Equal(t, 10.0, invoice.MaxTaxReduction)
	assert.Equal(t, "Sr. Rodriguez", invoice.OurReference)
	assert.Equal(t, "Mr. Bean", invoice.YourReference)
	assert.Equal(t, "Söderåsens BRF", invoice.RecipientName)
	assert.Equal(t, "Storgatan 1", invoice.RecipientStreetAddress)
	assert.Equal(t, "12345", invoice.RecipientPostalCode)
	assert.Equal(t, "Stockholm", invoice.RecipientCity)
	assert.Equal(t, 209, invoice.RecipientCountryID)
	assert.Equal(t, "Rot", invoice.RotRutText)
	assert.Equal(t, "2026-10-13 09:40:57", invoice.SentToEconomySystem.String())
	assert.Equal(t, 7800.0, invoice.Sum)
	assert.Equal(t, 10.0, invoice.TaxReductionAmountInPercent)
	assert.Equal(t, "Rut", invoice.TaxReductionMethod)
	assert.Equal(t, "123456-7890", invoice.TenantOrganisationalNumber)

	require.Len(t, invoice.Rows, 1)
	assert.Equal(t, InvoiceRow{
		ID:                   1,
		Article:              "Article #4",
		ArticleNumber:        "4",
		BookingProjectNumber: "1234",
		CostCenter:           &CostCenterRef{ID: 1, Name: "Cost center 1", Code: "CC1"},
		Comment:              "It's a comment",
		DateStr:              "2026-10-13",
		Discount:             10,
		Header1:              "Header 1",
		Header2:              "Header 2",
		Header3:              "Header 3",
		Header4:              "Header 4",
		Header5:              "Header 5",
		Header6:              "Header 6",
		Header7:              "Header 7",
		Header8:              "Header 8",
		Header9:              "Header 9",
		HouseWorkType:        &NameID{ID: 1, Name: "Garden work"},
		IndentLevel:          0,
		IsHeader:             false,
		IsSumRow:             false,
		SortOrder:            1,
		Rate:                 100,
		Sum:                  7500,
		Type:                 InvoiceRowTypeNone,
		VatCode:              &NameID{ID: 1, Name: "VAT code 1"},
		Value:                75,
		Unit:                 "m2",
	}, invoice.Rows[0])

	assert.Equal(t, []InvoiceTaxReduction{{ID: 1, Name: "Rut", SocialSecurityNumber: "123456-7890", TaxReduction: 10}}, invoice.TaxReductions)
}

func TestInvoiceRowTypeMatchesLowercaseAndKeepsUnknown(t *testing.T) {
	var rows []InvoiceRow
	require.NoError(t, json.Unmarshal([]byte(`[{"type":"time"},{"type":"Purchase"},{"type":"travel"},{"type":"PeriodRow"},{"type":"Fee"}]`), &rows))

	assert.Equal(t, InvoiceRowTypeTime, rows[0].Type)
	assert.Equal(t, InvoiceRowTypePurchase, rows[1].Type)
	assert.Equal(t, InvoiceRowTypeTravel, rows[2].Type)
	assert.Equal(t, InvoiceRowTypePeriodRow, rows[3].Type)
	assert.Equal(t, InvoiceRowType("Fee"), rows[4].Type)
}

func TestInvoicesSetSentToEconomySystemSendsQueryWithoutBody(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusNoContent, ""))

	err := client.Invoices.SetSentToEconomySystem(context.Background(), 42, NewDate(2025, time.January, 1), "invoice2")

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodPatch, request.Method)
	assert.Equal(t, "/v1/core/invoices/42/setsenttoeconomysystem", request.Path)
	assert.Equal(t, map[string][]string{"date": {"2025-01-01"}, "economySystemInvoiceNumber": {"invoice2"}}, request.Query)
	assert.Empty(t, request.Body)
}

func TestInvoicesSetSentToEconomySystemOmitsEmptyInvoiceNumber(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusNoContent, ""))

	err := client.Invoices.SetSentToEconomySystem(context.Background(), 42, NewDate(2025, time.January, 1), "")

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{"date": {"2025-01-01"}}, handler.last().Query)
}

func TestInvoicesListAttachmentsDecodesFixture(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/core/invoices/42/attachments", "invoices_attachments_list.json"))

	page, err := client.Invoices.ListAttachments(context.Background(), 42, ListOptions{PageSize: 5})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{"page": {"1"}, "pageSize": {"5"}}, handler.last().Query)
	require.Len(t, page.Items, 2)
	drawing := page.Items[0]
	assert.Equal(t, 1, drawing.ID)
	assert.Equal(t, "Drawing.pdf", drawing.Name)
	assert.Equal(t, "application/pdf", drawing.ContentType)
	assert.Equal(t, 245760, drawing.SizeInBytes)
	assert.Equal(t, 92, drawing.ContactID)
	assert.Zero(t, drawing.ProjectID)
	assert.False(t, drawing.IsTimeJournal)
	assert.Equal(t, 1, drawing.SortIndex)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, drawing.CreatedByUser)
	assert.Equal(t, "2026-10-11 09:40:57", drawing.CreatedDate.String())
	journal := page.Items[1]
	assert.Zero(t, journal.ContactID)
	assert.Equal(t, 17, journal.ProjectID)
	assert.True(t, journal.IsTimeJournal)
}

func TestInvoicesListAllAttachmentsWalksPages(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "invoices_attachments_list.json"))

	var names []string
	for attachment, err := range client.Invoices.ListAllAttachments(context.Background(), 42, ListOptions{}) {
		require.NoError(t, err)
		names = append(names, attachment.Name)
	}

	assert.Equal(t, []string{"Drawing.pdf", "TimeJournal.pdf"}, names)
	request := handler.last()
	assert.Equal(t, "/v1/core/invoices/42/attachments", request.Path)
	assert.Equal(t, map[string][]string{"page": {"1"}, "pageSize": {"100"}}, request.Query)
}

func TestInvoicesDownloadAttachment(t *testing.T) {
	client := newTestClient(t, routes(t, map[string]http.HandlerFunc{
		"GET /v1/core/invoices/42/attachments/1/download": serveFile("application/pdf", `attachment; filename="Drawing.pdf"`, []byte("%PDF-1.7 drawing")),
	}))

	download, err := client.Invoices.DownloadAttachment(context.Background(), 42, 1)

	require.NoError(t, err)
	defer download.Content.Close()
	content, err := io.ReadAll(download.Content)
	require.NoError(t, err)
	assert.Equal(t, "%PDF-1.7 drawing", string(content))
	assert.Equal(t, "application/pdf", download.ContentType)
	assert.Equal(t, "Drawing.pdf", download.Filename)
}
