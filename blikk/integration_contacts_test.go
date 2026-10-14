//go:build integration

package blikk

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationContacts(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.Contacts.List(ctx, ContactListParams{ListOptions: ListOptions{PageSize: 10, SortBy: ContactSortByName}})
	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("no contacts in the tenant")
	}
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)
	assert.NotEmpty(t, listed.Name)
	assert.Contains(t, []ContactType{ContactTypePerson, ContactTypeCompany}, listed.ContactType)

	contact, err := client.Contacts.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, contact.ID)
	assert.Equal(t, listed.Name, contact.Name)
	assert.Equal(t, listed.ContactType, contact.ContactType)
	assert.Equal(t, listed.GUID, contact.GUID)

	if listed.CustomerNumber != "" {
		byCustomerNumber, err := client.Contacts.GetByCustomerNumber(ctx, listed.CustomerNumber)
		require.NoError(t, err)
		assert.Equal(t, listed.ID, byCustomerNumber.ID)
	}
}

func TestIntegrationContactsFilterByType(t *testing.T) {
	client := integrationClient(t)

	page, err := client.Contacts.List(integrationContext(t), ContactListParams{ListOptions: ListOptions{PageSize: 10}, ContactType: ContactTypeCompany})

	require.NoError(t, err)
	for _, contact := range page.Items {
		assert.Equal(t, ContactTypeCompany, contact.ContactType)
	}
}

func TestIntegrationContactByExternalID(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	var externalID string
	for contact, err := range client.Contacts.ListAll(ctx, ContactListParams{}) {
		require.NoError(t, err)
		if contact.ExternalID != "" {
			externalID = contact.ExternalID
			break
		}
	}
	if externalID == "" {
		t.Skip("no contact with an external id in the tenant")
	}

	contact, err := client.Contacts.GetByExternalID(ctx, externalID)
	require.NoError(t, err)
	assert.Equal(t, externalID, contact.ExternalID)
}

func TestIntegrationContactTags(t *testing.T) {
	client := integrationClient(t)

	page, err := client.ContactTags.List(integrationContext(t), ContactTagListParams{ListOptions: ListOptions{PageSize: 10, SortBy: ContactTagSortByName}})

	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("no contact tags in the tenant")
	}
	assert.NotZero(t, page.Items[0].ID)
	assert.False(t, page.Items[0].CreatedDate.IsZero())
}

func TestIntegrationDeletedEntities(t *testing.T) {
	client := integrationClient(t)
	now := time.Now()

	page, err := client.DeletedEntities.List(integrationContext(t), DeletedEntityListParams{
		PageSize: 10,
		Filter: DeletedEntityFilter{
			DeletedFromDate: DateTimeOf(now.AddDate(0, 0, -90)),
			DeletedToDate:   DateTimeOf(now),
		},
	})

	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("no entities deleted in the last 90 days")
	}
	deleted := page.Items[0]
	assert.NotZero(t, deleted.EntityID)
	assert.NotEmpty(t, deleted.EntityType)
	require.NotNil(t, deleted.DeletedDate)
	assert.False(t, deleted.DeletedDate.Time().Before(DateOf(now.AddDate(0, 0, -91)).Time()))
}

func TestIntegrationDeletedEntitiesFilterByType(t *testing.T) {
	client := integrationClient(t)

	page, err := client.DeletedEntities.List(integrationContext(t), DeletedEntityListParams{
		PageSize: 10,
		Filter:   DeletedEntityFilter{EntityTypes: []EntityType{EntityTypeTimeReport}},
	})

	require.NoError(t, err)
	for _, deleted := range page.Items {
		assert.Equal(t, EntityTypeTimeReport, deleted.EntityType)
	}
}

func TestIntegrationLogs(t *testing.T) {
	client := integrationClient(t)

	page, err := client.Logs.List(integrationContext(t), LogListParams{
		ListOptions: ListOptions{PageSize: 10, SortBy: LogSortByCreatedDate, SortOrder: SortDescending},
		LogType:     LogTypeProjectAll,
	})

	require.NoError(t, err)
	if len(page.Items) == 0 {
		t.Skip("no logs in the tenant")
	}
	log := page.Items[0]
	assert.NotZero(t, log.ID)
	assert.NotEmpty(t, log.LogType)
	assert.NotZero(t, log.Entity.ID)
	assert.Equal(t, EntityTypeProject, log.Entity.EntityType)
	assert.False(t, log.CreatedDate.IsZero())
}

func TestIntegrationDeletedEntitiesWithoutFilter(t *testing.T) {
	client := integrationClient(t)

	page, err := client.DeletedEntities.List(integrationContext(t), DeletedEntityListParams{PageSize: 1})

	require.NoError(t, err)
	assert.Equal(t, 1, page.Page)
}
