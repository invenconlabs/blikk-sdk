//go:build integration

package blikk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationOffers(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.Offers.List(ctx, OfferListParams{ListOptions: ListOptions{PageSize: 10, SortBy: OfferSortByCreatedDate, SortOrder: SortDescending}})

	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("tenant has no offers")
	}
	offer := page.Items[0]
	assert.NotZero(t, offer.ID)
	assert.NotZero(t, offer.User.ID)
	assert.False(t, offer.CreatedDate.IsZero())
}

func TestIntegrationOpportunities(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.Opportunities.List(ctx, OpportunityListParams{ListOptions: ListOptions{PageSize: 10, SortBy: OpportunitySortByUpdatedDate, SortOrder: SortDescending}})

	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("tenant has no opportunities")
	}
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)
	assert.NotZero(t, listed.User.ID)
	assert.False(t, listed.CreatedDate.IsZero())

	opportunity, err := client.Opportunities.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, opportunity.ID)
	assert.Equal(t, listed.Title, opportunity.Title)
	assert.Equal(t, listed.User, opportunity.User)
	assert.Equal(t, listed.OpportunityState, opportunity.OpportunityState)
}

func TestIntegrationOpportunityStatuses(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.OpportunityStatuses.List(ctx, OpportunityStatusListParams{ListOptions: ListOptions{PageSize: 10, SortBy: OpportunityStatusSortBySortOrder}})

	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("tenant has no opportunity statuses")
	}
	status := page.Items[0]
	assert.NotZero(t, status.ID)
	assert.NotEmpty(t, status.Name)
	assert.False(t, status.CreatedDate.IsZero())
}

func TestIntegrationOpportunityTags(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.OpportunityTags.List(ctx, OpportunityTagListParams{ListOptions: ListOptions{PageSize: 10, SortBy: OpportunityTagSortByName}})

	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("tenant has no opportunity tags")
	}
	tag := page.Items[0]
	assert.NotZero(t, tag.ID)
	assert.NotEmpty(t, tag.Title)
	assert.False(t, tag.CreatedDate.IsZero())
}
