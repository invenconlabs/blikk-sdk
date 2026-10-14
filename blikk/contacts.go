package blikk

import (
	"context"
	"iter"
	"strconv"
)

const contactsPath = "v1/Core/Contacts"

type ContactsService struct {
	client *Client
}

type ContactType string

const (
	ContactTypePerson  ContactType = "Person"
	ContactTypeCompany ContactType = "Company"
)

func (t *ContactType) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, t, ContactTypePerson, ContactTypeCompany)
}

type ContactRelation string

const (
	ContactRelationCustomer      ContactRelation = "customer"
	ContactRelationSupplier      ContactRelation = "supplier"
	ContactRelationWholesale     ContactRelation = "wholesale"
	ContactRelationRetail        ContactRelation = "retail"
	ContactRelationSubcontractor ContactRelation = "subcontractor"
	ContactRelationPartner       ContactRelation = "partner"
)

type ContactVatType string

const (
	ContactVatTypeSweden              ContactVatType = "Sweden"
	ContactVatTypeSwedenReverseCharge ContactVatType = "SwedenReverseCharge"
	ContactVatTypeEu                  ContactVatType = "Eu"
	ContactVatTypeEuReverseCharge     ContactVatType = "EuReverseCharge"
	ContactVatTypeExport              ContactVatType = "Export"
)

func (t *ContactVatType) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, t, ContactVatTypeSweden, ContactVatTypeSwedenReverseCharge, ContactVatTypeEu, ContactVatTypeEuReverseCharge, ContactVatTypeExport)
}

const (
	ContactSortByName           = "name"
	ContactSortByCustomerNumber = "customerNumber"
	ContactSortByCreated        = "created"
	ContactSortByUpdated        = "updated"
)

type ContactListParams struct {
	ListOptions
	Query                 string            `query:"filter.query"`
	PhoneNumber           string            `query:"filter.phoneNumber"`
	TagIDs                []int             `query:"filter.tagIds"`
	CreatedOrUpdatedSince DateTime          `query:"filter.createdOrUpdatedSince"`
	ContactType           ContactType       `query:"filter.contactType"`
	Relations             []ContactRelation `query:"filter.relations"`
}

type ContactListItem struct {
	ID                                   int         `json:"id"`
	GUID                                 string      `json:"guid"`
	Name                                 string      `json:"name"`
	PhoneNumber                          string      `json:"phoneNumber"`
	CellPhoneNumber                      string      `json:"cellPhoneNumber"`
	Email                                string      `json:"email"`
	CustomerNumber                       string      `json:"customerNumber"`
	ContactType                          ContactType `json:"contactType"`
	OrganizationalOrSocialSecurityNumber string      `json:"organizationalOrSocialSecurityNumber"`
	VisitingAddress                      *Address    `json:"visitingAddress"`
	Company                              *NameID     `json:"company"`
	ParentCompany                        *NameID     `json:"parentCompany"`
	Created                              *DateTime   `json:"created"`
	LastUpdated                          *DateTime   `json:"lastUpdated"`
	ExternalID                           string      `json:"externalId"`
	IsKycCustomer                        bool        `json:"isKycCustomer"`
}

type Contact struct {
	ID                                   int             `json:"id"`
	GUID                                 string          `json:"guid"`
	Name                                 string          `json:"name"`
	CustomerNumber                       string          `json:"customerNumber"`
	ContactType                          ContactType     `json:"contactType"`
	OrganizationalOrSocialSecurityNumber string          `json:"organizationalOrSocialSecurityNumber"`
	Company                              *NameID         `json:"company"`
	ExternalID                           string          `json:"externalId"`
	ParentCompany                        *NameID         `json:"parentCompany"`
	IsActive                             bool            `json:"isActive"`
	Email                                string          `json:"email"`
	InvoiceEmail                         string          `json:"invoiceEmail"`
	PhoneNumber                          string          `json:"phoneNumber"`
	CellPhoneNumber                      string          `json:"cellPhoneNumber"`
	ImportantInformation                 string          `json:"importantInformation"`
	IsSupplier                           bool            `json:"isSupplier"`
	IsWholeSale                          bool            `json:"isWholeSale"`
	IsRetailer                           bool            `json:"isRetailer"`
	IsSubContractor                      bool            `json:"isSubContractor"`
	IsPartner                            bool            `json:"isPartner"`
	IsCustomer                           bool            `json:"isCustomer"`
	UseReverseBuildTax                   bool            `json:"useReverseBuildTax"`
	IsKycCustomer                        bool            `json:"isKycCustomer"`
	Responsible                          *NameID         `json:"responsible"`
	Tags                                 []NameID        `json:"tags"`
	VisitingAddress                      *Address        `json:"visitingAddress"`
	DeliveryAddress                      *Address        `json:"deliveryAddress"`
	InvoiceAddress                       *Address        `json:"invoiceAddress"`
	Created                              *DateTime       `json:"created"`
	CreatedBy                            *NameID         `json:"createdBy"`
	LastUpdated                          *DateTime       `json:"lastUpdated"`
	LastUpdatedBy                        *NameID         `json:"lastUpdatedBy"`
	ContactPersons                       []ContactPerson `json:"contactPersons"`
	VatNumber                            string          `json:"vatNumber"`
	InvoicingNotes                       string          `json:"invoicingNotes"`
	VatType                              ContactVatType  `json:"vatType"`
	ElectronicInvoiceReferenceCode       string          `json:"electronicInvoiceReferenceCode"`
}

type ContactPerson struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IsPrimary bool   `json:"isPrimary"`
}

// ContactInput keeps the relation flags and isActive without omitempty because Update replaces the whole contact, so false must be sent.
type ContactInput struct {
	Name                                 string               `json:"name"`
	CustomerNumber                       string               `json:"customerNumber,omitempty"`
	ContactType                          ContactType          `json:"contactType"`
	OrganizationalOrSocialSecurityNumber string               `json:"organizationalOrSocialSecurityNumber,omitempty"`
	CompanyID                            int                  `json:"companyId,omitempty"`
	ParentCompanyID                      int                  `json:"parentCompanyId,omitempty"`
	ExternalID                           string               `json:"externalId,omitempty"`
	IsActive                             bool                 `json:"isActive"`
	Email                                string               `json:"email,omitempty"`
	InvoiceEmail                         string               `json:"invoiceEmail,omitempty"`
	PhoneNumber                          string               `json:"phoneNumber,omitempty"`
	CellPhoneNumber                      string               `json:"cellPhoneNumber,omitempty"`
	ImportantInformation                 string               `json:"importantInformation,omitempty"`
	IsSupplier                           bool                 `json:"isSupplier"`
	IsWholeSale                          bool                 `json:"isWholeSale"`
	IsRetailer                           bool                 `json:"isRetailer"`
	IsSubContractor                      bool                 `json:"isSubContractor"`
	IsPartner                            bool                 `json:"isPartner"`
	IsCustomer                           bool                 `json:"isCustomer"`
	UseReverseBuildTax                   bool                 `json:"useReverseBuildTax"`
	ResponsibleID                        int                  `json:"responsibleId,omitempty"`
	TagIDs                               []int                `json:"tagIds,omitempty"`
	VisitingAddress                      *ContactAddressInput `json:"visitingAddress,omitempty"`
	DeliveryAddress                      *ContactAddressInput `json:"deliveryAddress,omitempty"`
	InvoiceAddress                       *ContactAddressInput `json:"invoiceAddress,omitempty"`
}

type ContactAddressInput struct {
	StreetAddress     string `json:"streetAddress,omitempty"`
	AdditionalAddress string `json:"additionalAddress,omitempty"`
	PostalCode        string `json:"postalCode,omitempty"`
	City              string `json:"city,omitempty"`
	State             string `json:"state,omitempty"`
	CountryID         int    `json:"countryId,omitempty"`
}

func (s *ContactsService) List(ctx context.Context, params ContactListParams) (*Page[ContactListItem], error) {
	return listPage[ContactListItem](ctx, s.client, contactsPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *ContactsService) ListAll(ctx context.Context, params ContactListParams) iter.Seq2[ContactListItem, error] {
	return listAll[ContactListItem](ctx, s.client, contactsPath, params)
}

func (s *ContactsService) Get(ctx context.Context, id int) (*Contact, error) {
	return getResource[Contact](ctx, s.client, contactsPath, strconv.Itoa(id))
}

// GetByExternalID fails with a server error when more than one contact shares the external id.
func (s *ContactsService) GetByExternalID(ctx context.Context, externalID string) (*Contact, error) {
	return getResource[Contact](ctx, s.client, contactsPath, "byExternalId", externalID)
}

// GetByCustomerNumber fails with a server error when more than one contact shares the customer number.
func (s *ContactsService) GetByCustomerNumber(ctx context.Context, customerNumber string) (*Contact, error) {
	return getResource[Contact](ctx, s.client, contactsPath, "byCustomerNumber", customerNumber)
}

func (s *ContactsService) Create(ctx context.Context, input ContactInput) (*Contact, error) {
	return createResource[Contact](ctx, s.client, input, contactsPath)
}

// Update replaces the whole contact; the API requires ResponsibleID here and clears ExternalID when it is left empty.
func (s *ContactsService) Update(ctx context.Context, id int, input ContactInput) error {
	return updateResource(ctx, s.client, input, contactsPath, strconv.Itoa(id))
}

func (s *ContactsService) Delete(ctx context.Context, id int) error {
	return deleteResource(ctx, s.client, contactsPath, strconv.Itoa(id))
}
