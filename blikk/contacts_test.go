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

func TestContactsListEncodesEveryFilter(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Core/Contacts", "contacts_list.json"))

	_, err := client.Contacts.List(context.Background(), ContactListParams{
		ListOptions:           ListOptions{PageSize: 10, SortBy: ContactSortByCustomerNumber, SortOrder: SortAscending},
		Query:                 "anders",
		PhoneNumber:           "070",
		TagIDs:                []int{10, 16},
		CreatedOrUpdatedSince: DateTimeOf(time.Date(2021, 1, 1, 14, 0, 0, 0, time.UTC)),
		ContactType:           ContactTypeCompany,
		Relations:             []ContactRelation{ContactRelationSupplier, ContactRelationWholesale},
	})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"page":                         {"1"},
		"pageSize":                     {"10"},
		"sortBy":                       {"customerNumber"},
		"sortOrder":                    {"ascending"},
		"filter.query":                 {"anders"},
		"filter.phoneNumber":           {"070"},
		"filter.tagIds":                {"10,16"},
		"filter.createdOrUpdatedSince": {"2021-01-01 14:00:00"},
		"filter.contactType":           {"Company"},
		"filter.relations":             {"supplier,wholesale"},
	}, handler.last().Query)
}

func TestContactsListDecodesFixture(t *testing.T) {
	client, _ := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Core/Contacts", "contacts_list.json"))

	page, err := client.Contacts.List(context.Background(), ContactListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	contact := page.Items[0]
	assert.Equal(t, 1, contact.ID)
	assert.Empty(t, contact.GUID)
	assert.Equal(t, "John Doe", contact.Name)
	assert.Empty(t, contact.PhoneNumber)
	assert.Equal(t, "John.doe@gmail.com", contact.Email)
	assert.Equal(t, "1", contact.CustomerNumber)
	assert.Equal(t, ContactTypePerson, contact.ContactType)
	assert.Equal(t, &Address{StreetAddress: "Trädvägen 123", PostalCode: "94100", City: "Piteå", State: "Norrbotten", CountryID: 1, CountryName: "Sverige"}, contact.VisitingAddress)
	assert.Equal(t, &NameID{ID: 2, Name: "My company"}, contact.Company)
	assert.Nil(t, contact.ParentCompany)
	assert.Equal(t, "2026-10-13 09:40:57", contact.Created.String())
	assert.Nil(t, contact.LastUpdated)
	assert.Empty(t, contact.ExternalID)
	assert.False(t, contact.IsKycCustomer)
}

func TestContactsGetDecodesFixture(t *testing.T) {
	client, _ := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Core/Contacts/32", "contact.json"))

	contact, err := client.Contacts.Get(context.Background(), 32)

	require.NoError(t, err)
	assert.Equal(t, 32, contact.ID)
	assert.Equal(t, "Lotta Larsson", contact.Name)
	assert.Equal(t, "111", contact.CustomerNumber)
	assert.Equal(t, ContactTypePerson, contact.ContactType)
	assert.Equal(t, "860311-7063", contact.OrganizationalOrSocialSecurityNumber)
	assert.Nil(t, contact.Company)
	assert.Equal(t, "12345", contact.ExternalID)
	assert.True(t, contact.IsActive)
	assert.Equal(t, "lotta@familjenlarsson.se", contact.Email)
	assert.Equal(t, "075-9902281", contact.PhoneNumber)
	assert.Equal(t, "0707788994", contact.CellPhoneNumber)
	assert.False(t, contact.IsCustomer)
	assert.Equal(t, "SE556677889901", contact.VatNumber)
	assert.Equal(t, "Mark the invoice with the project number", contact.InvoicingNotes)
	assert.Equal(t, ContactVatTypeSweden, contact.VatType)
	assert.Equal(t, "REF-4711", contact.ElectronicInvoiceReferenceCode)
	assert.Equal(t, &NameID{ID: 5, Name: "Jonas lindberg"}, contact.Responsible)
	assert.Empty(t, contact.ContactPersons)
	assert.Empty(t, contact.Tags)
	assert.Equal(t, &Address{StreetAddress: "Solberga", PostalCode: "92372", CountryID: 1, CountryName: "Sverige"}, contact.VisitingAddress)
	assert.Equal(t, &Address{CountryID: 1, CountryName: "Sverige"}, contact.DeliveryAddress)
	assert.Equal(t, &Address{CountryID: 1, CountryName: "Sverige"}, contact.InvoiceAddress)
	assert.Equal(t, "2021-06-17 05:43:54", contact.Created.String())
	assert.Equal(t, &NameID{ID: 5, Name: "Jonas lindberg"}, contact.CreatedBy)
	assert.Nil(t, contact.LastUpdated)
	assert.Nil(t, contact.LastUpdatedBy)
}

func TestContactDecodesFieldsMissingFromDocsExample(t *testing.T) {
	var contact Contact
	err := json.Unmarshal([]byte(`{"guid":"80e191f0-fd53-463c-9856-78f6bc088616","contactType":"Company","isKycCustomer":true,"vatType":null,"tags":[{"objectName":"generic.name_id","id":3,"name":"VIP"}],"contactPersons":[{"objectName":"contact.person","id":7,"name":"Anna","isPrimary":true}]}`), &contact)

	require.NoError(t, err)
	assert.Equal(t, "80e191f0-fd53-463c-9856-78f6bc088616", contact.GUID)
	assert.Equal(t, ContactTypeCompany, contact.ContactType)
	assert.True(t, contact.IsKycCustomer)
	assert.Equal(t, ContactVatType(""), contact.VatType)
	assert.Equal(t, []NameID{{ID: 3, Name: "VIP"}}, contact.Tags)
	assert.Equal(t, []ContactPerson{{ID: 7, Name: "Anna", IsPrimary: true}}, contact.ContactPersons)
}

func TestContactsGetByExternalIDAndCustomerNumberUseLookupPaths(t *testing.T) {
	client, handler := newRecordingClient(t, routes(t, map[string]http.HandlerFunc{
		"GET /v1/Core/Contacts/byExternalId/ext%2F1": serveFixture(t, "contact.json"),
		"GET /v1/Core/Contacts/byCustomerNumber/111": serveFixture(t, "contact.json"),
	}))

	byExternalID, err := client.Contacts.GetByExternalID(context.Background(), "ext/1")
	require.NoError(t, err)
	byCustomerNumber, err := client.Contacts.GetByCustomerNumber(context.Background(), "111")
	require.NoError(t, err)

	assert.Equal(t, 32, byExternalID.ID)
	assert.Equal(t, 32, byCustomerNumber.ID)
	assert.Len(t, handler.all(), 2)
}

func TestContactsCreateSendsFullBodyAndReturnsContact(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusCreated, string(loadFixture(t, "contact.json"))))

	contact, err := client.Contacts.Create(context.Background(), fullContactInput())

	require.NoError(t, err)
	assert.Equal(t, 32, contact.ID)
	assert.Equal(t, "Lotta Larsson", contact.Name)
	assert.Equal(t, ContactTypePerson, contact.ContactType)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/Core/Contacts", request.Path)
	assert.JSONEq(t, `{
		"name": "Lotta Larsson",
		"customerNumber": "111",
		"contactType": "Person",
		"organizationalOrSocialSecurityNumber": "860311-7063",
		"companyId": 2,
		"parentCompanyId": 3,
		"externalId": "12345",
		"isActive": true,
		"email": "lotta@familjenlarsson.se",
		"invoiceEmail": "invoice@familjenlarsson.se",
		"phoneNumber": "075-9902281",
		"cellPhoneNumber": "0707788994",
		"importantInformation": "Call first",
		"isSupplier": true,
		"isWholeSale": true,
		"isRetailer": true,
		"isSubContractor": true,
		"isPartner": true,
		"isCustomer": true,
		"useReverseBuildTax": true,
		"responsibleId": 15,
		"tagIds": [10, 16],
		"visitingAddress": {"streetAddress": "Solberga", "additionalAddress": "c/o Larsson", "postalCode": "92372", "city": "Lycksele", "state": "Västerbotten", "countryId": 1},
		"deliveryAddress": {"countryId": 1},
		"invoiceAddress": {"city": "Umeå"}
	}`, string(request.Body))
}

func TestContactsCreateSendsMinimalBody(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusCreated, `{"id":33}`))

	contact, err := client.Contacts.Create(context.Background(), ContactInput{Name: "Lotta Larsson", ContactType: ContactTypePerson})

	require.NoError(t, err)
	assert.Equal(t, 33, contact.ID)
	assert.JSONEq(t, `{
		"name": "Lotta Larsson",
		"contactType": "Person",
		"isActive": false,
		"isSupplier": false,
		"isWholeSale": false,
		"isRetailer": false,
		"isSubContractor": false,
		"isPartner": false,
		"isCustomer": false,
		"useReverseBuildTax": false
	}`, string(handler.last().Body))
}

func TestContactInputMatchesDocumentedRequest(t *testing.T) {
	var documented map[string]any
	require.NoError(t, json.Unmarshal(loadFixture(t, "contact_input.json"), &documented))
	var input ContactInput
	require.NoError(t, json.Unmarshal(loadFixture(t, "contact_input.json"), &input))

	encoded, err := json.Marshal(input)
	require.NoError(t, err)

	var reencoded map[string]any
	require.NoError(t, json.Unmarshal(encoded, &reencoded))
	assert.Equal(t, withoutEmptyValues(documented), withoutEmptyValues(reencoded))
}

func TestContactsUpdateSendsPut(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	err := client.Contacts.Update(context.Background(), 32, fullContactInput())

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodPut, request.Method)
	assert.Equal(t, "/v1/Core/Contacts/32", request.Path)
	body := request.jsonBody(t)
	assert.Equal(t, "Lotta Larsson", body["name"])
	assert.Equal(t, 15.0, body["responsibleId"])
	assert.Equal(t, []any{10.0, 16.0}, body["tagIds"])
	assert.NotContains(t, body, "tags")
}

func TestContactsUpdateSendsMinimalBody(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	err := client.Contacts.Update(context.Background(), 32, ContactInput{Name: "Lotta Larsson", ContactType: ContactTypeCompany, ResponsibleID: 15})

	require.NoError(t, err)
	assert.JSONEq(t, `{
		"name": "Lotta Larsson",
		"contactType": "Company",
		"responsibleId": 15,
		"isActive": false,
		"isSupplier": false,
		"isWholeSale": false,
		"isRetailer": false,
		"isSubContractor": false,
		"isPartner": false,
		"isCustomer": false,
		"useReverseBuildTax": false
	}`, string(handler.last().Body))
}

func TestContactsDeleteSendsDelete(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	err := client.Contacts.Delete(context.Background(), 32)

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodDelete, request.Method)
	assert.Equal(t, "/v1/Core/Contacts/32", request.Path)
}

func fullContactInput() ContactInput {
	return ContactInput{
		Name:                                 "Lotta Larsson",
		CustomerNumber:                       "111",
		ContactType:                          ContactTypePerson,
		OrganizationalOrSocialSecurityNumber: "860311-7063",
		CompanyID:                            2,
		ParentCompanyID:                      3,
		ExternalID:                           "12345",
		IsActive:                             true,
		Email:                                "lotta@familjenlarsson.se",
		InvoiceEmail:                         "invoice@familjenlarsson.se",
		PhoneNumber:                          "075-9902281",
		CellPhoneNumber:                      "0707788994",
		ImportantInformation:                 "Call first",
		IsSupplier:                           true,
		IsWholeSale:                          true,
		IsRetailer:                           true,
		IsSubContractor:                      true,
		IsPartner:                            true,
		IsCustomer:                           true,
		UseReverseBuildTax:                   true,
		ResponsibleID:                        15,
		TagIDs:                               []int{10, 16},
		VisitingAddress:                      &ContactAddressInput{StreetAddress: "Solberga", AdditionalAddress: "c/o Larsson", PostalCode: "92372", City: "Lycksele", State: "Västerbotten", CountryID: 1},
		DeliveryAddress:                      &ContactAddressInput{CountryID: 1},
		InvoiceAddress:                       &ContactAddressInput{City: "Umeå"},
	}
}
