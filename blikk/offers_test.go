package blikk

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOffersListEncodesEveryFilter(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "offers_list.json")}
	client := newTestClient(t, handler.ServeHTTP)
	moment := DateTimeOf(time.Date(2021, 1, 1, 14, 0, 0, 0, time.UTC))

	_, err := client.Offers.List(context.Background(), OfferListParams{
		ListOptions:     ListOptions{SortBy: OfferSortByOfferNumber, SortOrder: SortDescending},
		OfferNumber:     12,
		CreatedDateFrom: moment,
		CreatedDateTo:   moment,
		UpdatedDateFrom: moment,
		UpdatedDateTo:   moment,
		OfferState:      OfferStateSentToCustomer,
		OpportunityID:   3,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Core/Offers", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                   {"1"},
		"sortBy":                 {"offerNumber"},
		"sortOrder":              {"descending"},
		"filter.offerNumber":     {"12"},
		"filter.createdDateFrom": {"2021-01-01 14:00:00"},
		"filter.createdDateTo":   {"2021-01-01 14:00:00"},
		"filter.updatedDateFrom": {"2021-01-01 14:00:00"},
		"filter.updatedDateTo":   {"2021-01-01 14:00:00"},
		"filter.offerState":      {"SentToCustomer"},
		"filter.opportunityId":   {"3"},
	}, request.Query)
}

func TestOffersListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixture(t, "offers_list.json"))

	page, err := client.Offers.List(context.Background(), OfferListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	offer := page.Items[0]
	assert.Equal(t, 1, offer.ID)
	assert.Equal(t, 1, offer.OfferNumber)
	assert.Equal(t, 1, offer.Version)
	assert.Equal(t, "Lorem Ipsum", offer.Title)
	assert.InDelta(t, 10000.0, offer.OfferValue, 0)
	assert.True(t, offer.HasPdf)
	assert.Equal(t, NameID{ID: 1, Name: "John Doe"}, offer.User)
	assert.Equal(t, &NameID{ID: 1, Name: "Ipsum"}, offer.Customer)
	assert.Nil(t, offer.Opportunity)
	assert.Equal(t, OfferStateSentToCustomer, offer.OfferState)
	assert.Equal(t, "2026-10-13 09:40:57", offer.UpdatedDate.String())
	assert.Equal(t, "2026-10-13 09:40:57", offer.CreatedDate.String())
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, offer.CreatedByUser)
}

func TestOfferStateAcceptsLowercaseDocsSpelling(t *testing.T) {
	var state OfferState

	require.NoError(t, state.UnmarshalJSON([]byte(`"acceptedSigned"`)))

	assert.Equal(t, OfferStateAcceptedSigned, state)
}

func TestOffersListAllWalksPages(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "offers_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	var ids []int
	for offer, err := range client.Offers.ListAll(context.Background(), OfferListParams{}) {
		require.NoError(t, err)
		ids = append(ids, offer.ID)
	}

	assert.Equal(t, []int{1}, ids)
	assert.Equal(t, []string{"100"}, handler.last().Query["pageSize"])
}

func TestOffersDownloadStreamsPdf(t *testing.T) {
	disposition := `attachment; filename="Offer 1.pdf"`
	client, handler := newRecordingClient(t, routes(t, map[string]http.HandlerFunc{
		"GET /v1/Core/Offers/7/Download": serveFile("application/pdf", disposition, []byte("%PDF-1.7")),
	}))

	download, err := client.Offers.Download(context.Background(), 7)

	require.NoError(t, err)
	defer download.Content.Close()
	content, err := io.ReadAll(download.Content)
	require.NoError(t, err)
	assert.Equal(t, "%PDF-1.7", string(content))
	assert.Equal(t, "application/pdf", download.ContentType)
	assert.Equal(t, "Offer 1.pdf", download.Filename)
	assert.Equal(t, "/v1/Core/Offers/7/Download", handler.last().Path)
}
