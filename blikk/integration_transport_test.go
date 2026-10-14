//go:build integration

package blikk

import (
	"io"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type deletedEntityProbe struct {
	EntityType string `json:"entityType"`
	EntityID   int    `json:"entityId"`
}

func TestIntegrationListPageWithBody(t *testing.T) {
	client := integrationClient(t)
	now := time.Now().UTC()
	filter := map[string]string{
		"deletedFromDate": now.AddDate(0, 0, -30).Format(time.RFC3339),
		"deletedToDate":   now.Format(time.RFC3339),
	}

	page, err := listPageWithBody[deletedEntityProbe](integrationContext(t), client, "v1/Core/DeletedEntities", pagedBody{Filter: filter, PageSize: 5})

	require.NoError(t, err)
	assert.Equal(t, 1, page.Page)
	assert.Equal(t, 5, page.PageSize)
	for _, entity := range page.Items {
		assert.NotEmpty(t, entity.EntityType)
		assert.NotZero(t, entity.EntityID)
	}
}

func TestIntegrationGetResourceWithQuery(t *testing.T) {
	client := integrationClient(t)

	country, err := getResourceWithQuery[NameID](integrationContext(t), client, "v1/admin/countries", url.Values{"code": {"SWE"}})

	require.NoError(t, err)
	assert.NotZero(t, country.ID)
	assert.NotEmpty(t, country.Name)
}

func TestIntegrationDownloadOfferPDF(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)
	offers, err := listPage[NameID](ctx, client, "v1/Core/Offers", ListOptions{PageSize: 1})
	require.NoError(t, err)
	if len(offers.Items) == 0 {
		t.Skip("no offers to download")
	}

	download, err := client.download(ctx, "v1/Core/Offers/"+strconv.Itoa(offers.Items[0].ID)+"/Download", nil)

	require.NoError(t, err)
	defer download.Content.Close()
	header := make([]byte, 4)
	_, err = io.ReadFull(download.Content, header)
	require.NoError(t, err)
	assert.Equal(t, "%PDF", string(header))
	assert.Contains(t, download.ContentType, "pdf")
	t.Logf("offer download filename %q content type %q", download.Filename, download.ContentType)
}

func TestIntegrationDownloadMissingInvoiceAttachment(t *testing.T) {
	client := integrationClient(t)

	download, err := client.download(integrationContext(t), "v1/core/invoices/999999999/attachments/999999999/download", nil)

	assert.Nil(t, download)
	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	t.Logf("missing attachment: %v", err)
	assert.True(t, IsNotFound(err), "expected not found, got %v", err)
}
