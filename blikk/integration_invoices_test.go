//go:build integration

package blikk

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationInvoices(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.Invoices.List(ctx, InvoiceListParams{ListOptions: ListOptions{PageSize: 10, SortBy: InvoiceSortByCreatedDate, SortOrder: SortDescending}})
	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("no invoices in the tenant")
	}
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)
	assert.False(t, listed.CreatedDate.IsZero())
	assert.False(t, listed.InvoiceDate.IsZero())

	invoice, err := client.Invoices.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, invoice.ID)
	assert.Equal(t, listed.Sum, invoice.Sum)
	assert.Equal(t, listed.InvoiceDate, invoice.InvoiceDate)
	assert.NotZero(t, invoice.InvoiceNumber)
	for _, row := range invoice.Rows {
		assert.NotZero(t, row.ID)
		assert.NotEmpty(t, row.Type)
	}
}

func TestIntegrationInvoiceAttachments(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	invoiceID, attachment := firstInvoiceAttachment(t, ctx, client)
	if attachment == nil {
		t.Skip("none of the 100 newest invoices has an attachment")
	}
	assert.NotZero(t, attachment.ID)
	assert.NotEmpty(t, attachment.Name)

	download, err := client.Invoices.DownloadAttachment(ctx, invoiceID, attachment.ID)
	require.NoError(t, err)
	defer download.Content.Close()
	content, err := io.ReadAll(download.Content)
	require.NoError(t, err)
	assert.Len(t, content, attachment.SizeInBytes)
	t.Logf("attachment %q content type %q, download filename %q content type %q", attachment.Name, attachment.ContentType, download.Filename, download.ContentType)
}

func firstInvoiceAttachment(t *testing.T, ctx context.Context, client *Client) (int, *InvoiceAttachment) {
	t.Helper()
	invoices, err := client.Invoices.List(ctx, InvoiceListParams{ListOptions: ListOptions{PageSize: MaxPageSize, SortBy: InvoiceSortByCreatedDate, SortOrder: SortDescending}})
	require.NoError(t, err)
	for _, invoice := range invoices.Items {
		attachments, err := client.Invoices.ListAttachments(ctx, invoice.ID, ListOptions{PageSize: 1})
		require.NoError(t, err)
		if len(attachments.Items) > 0 {
			return invoice.ID, &attachments.Items[0]
		}
	}
	return 0, nil
}

func TestIntegrationSupplierInvoices(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.SupplierInvoices.List(ctx, SupplierInvoiceListParams{ListOptions: ListOptions{PageSize: 10, SortBy: SupplierInvoiceSortByCreatedDate, SortOrder: SortDescending}})
	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("no supplier invoices in the tenant")
	}
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)
	assert.NotEmpty(t, listed.Title)
	assert.NotZero(t, listed.Project.ID)
	assert.False(t, listed.CreatedDate.IsZero())

	invoice, err := client.SupplierInvoices.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, invoice.ID)
	assert.Equal(t, listed.Title, invoice.Title)
	assert.Equal(t, listed.InvoiceDate, invoice.InvoiceDate)
	assert.Equal(t, listed.TotalAmount, invoice.TotalAmount)
	assert.Equal(t, listed.Project, invoice.Project)
}

func TestIntegrationPaymentPlans(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.PaymentPlans.List(ctx, PaymentPlanListParams{ListOptions: ListOptions{PageSize: 10, SortBy: PaymentPlanSortByCreatedDate, SortOrder: SortDescending}})
	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("no payment plans in the tenant")
	}
	plan := page.Items[0]
	assert.NotZero(t, plan.ID)
	assert.NotEmpty(t, plan.BatchID)
	assert.NotZero(t, plan.Project.ID)
	assert.False(t, plan.CreatedDate.IsZero())

	byBatch, err := client.PaymentPlans.List(ctx, PaymentPlanListParams{BatchID: plan.BatchID})
	require.NoError(t, err)
	require.NotEmpty(t, byBatch.Items)
	for _, item := range byBatch.Items {
		assert.Equal(t, plan.BatchID, item.BatchID)
	}
}

func TestIntegrationMaterialReports(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.MaterialReports.List(ctx, MaterialReportListParams{ListOptions: ListOptions{PageSize: 10, SortBy: MaterialReportSortByCreatedDate, SortOrder: SortDescending}})
	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("no material reports in the tenant")
	}
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)
	assert.NotZero(t, listed.Project.ID)
	assert.False(t, listed.CreatedDate.IsZero())

	report, err := client.MaterialReports.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed, *report)
}
