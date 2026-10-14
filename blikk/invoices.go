package blikk

import (
	"context"
	"iter"
	"net/url"
	"strconv"
)

const invoicesPath = "v1/core/invoices"

type InvoicesService struct {
	client *Client
}

const (
	InvoiceSortByCreatedDate         = "createdDate"
	InvoiceSortByInvoiceDate         = "invoiceDate"
	InvoiceSortByCustomer            = "customer"
	InvoiceSortBySum                 = "sum"
	InvoiceSortBySentToEconomySystem = "sentToEconomySystem"
)

type InvoiceListParams struct {
	ListOptions
	CustomerID   int      `query:"filter.customerId"`
	CreatedSince DateTime `query:"filter.createdSince"`
	SinceID      int      `query:"filter.sinceId"`
}

type InvoiceListItem struct {
	ID                  int       `json:"id"`
	CreatedDate         DateTime  `json:"createdDate"`
	InvoiceDate         DateTime  `json:"invoiceDate"`
	Customer            *NameID   `json:"customer"`
	FromDate            *DateTime `json:"fromDate"`
	ToDate              *DateTime `json:"toDate"`
	Sum                 float64   `json:"sum"`
	SentToEconomySystem *DateTime `json:"sentToEconomySystem"`
}

type Invoice struct {
	ID                             int                   `json:"id"`
	Archived                       bool                  `json:"archived"`
	Cadastral                      string                `json:"cadastral"`
	CreatedDate                    DateTime              `json:"createdDate"`
	Customer                       *NameID               `json:"customer"`
	CustomerReferenceMarking       string                `json:"customerReferenceMarking"`
	EconomySystemInvoiceNumber     string                `json:"economySystemInvoiceNumber"`
	ElectronicInvoiceReferenceCode string                `json:"electronicInvoiceReferenceCode"`
	FromDate                       DateTime              `json:"fromDate"`
	InvoiceDate                    DateTime              `json:"invoiceDate"`
	InvoiceNumber                  int                   `json:"invoiceNumber"`
	InvoiceText                    string                `json:"invoiceText"`
	InvoiceReference               string                `json:"invoiceReference"`
	MaxTaxReduction                float64               `json:"maxTaxReduction"`
	OurReference                   string                `json:"ourReference"`
	YourReference                  string                `json:"yourReference"`
	RecipientName                  string                `json:"recipientName"`
	RecipientStreetAddress         string                `json:"recipientStreetAddress"`
	RecipientPostalCode            string                `json:"recipientPostalCode"`
	RecipientCity                  string                `json:"recipientCity"`
	RecipientCountryID             int                   `json:"recipientCountryId"`
	RotRutText                     string                `json:"rotRutText"`
	SentToEconomySystem            *DateTime             `json:"sentToEconomySystem"`
	Sum                            float64               `json:"sum"`
	TaxReductionAmountInPercent    float64               `json:"taxReductionAmountInPercent"`
	TaxReductionMethod             string                `json:"taxReductionMethod"`
	TenantOrganisationalNumber     string                `json:"tenantOrganisationalNumber"`
	ToDate                         DateTime              `json:"toDate"`
	Rows                           []InvoiceRow          `json:"rows"`
	TaxReductions                  []InvoiceTaxReduction `json:"taxReductions"`
}

type InvoiceRowType string

const (
	InvoiceRowTypeNone      InvoiceRowType = "None"
	InvoiceRowTypeTime      InvoiceRowType = "Time"
	InvoiceRowTypeAllowance InvoiceRowType = "Allowance"
	InvoiceRowTypePurchase  InvoiceRowType = "Purchase"
)

func (t *InvoiceRowType) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, t, InvoiceRowTypeNone, InvoiceRowTypeTime, InvoiceRowTypeAllowance, InvoiceRowTypePurchase)
}

type InvoiceRow struct {
	ID                   int            `json:"id"`
	Article              string         `json:"article"`
	ArticleNumber        string         `json:"articleNumber"`
	BookingProjectNumber string         `json:"bookingProjectNumber"`
	CostCenter           *CostCenterRef `json:"costCenter"`
	Comment              string         `json:"comment"`
	DateStr              string         `json:"dateStr"`
	Discount             float64        `json:"discount"`
	Header1              string         `json:"header1"`
	Header2              string         `json:"header2"`
	Header3              string         `json:"header3"`
	Header4              string         `json:"header4"`
	Header5              string         `json:"header5"`
	Header6              string         `json:"header6"`
	Header7              string         `json:"header7"`
	Header8              string         `json:"header8"`
	Header9              string         `json:"header9"`
	HouseWorkType        *NameID        `json:"houseWorkType"`
	IndentLevel          int            `json:"indentLevel"`
	IsHeader             bool           `json:"isHeader"`
	IsSumRow             bool           `json:"isSumRow"`
	SortOrder            int            `json:"sortOrder"`
	Rate                 float64        `json:"rate"`
	Sum                  float64        `json:"sum"`
	Type                 InvoiceRowType `json:"type"`
	VatCode              *NameID        `json:"vatCode"`
	Value                float64        `json:"value"`
	Unit                 string         `json:"unit"`
}

type InvoiceTaxReduction struct {
	ID                   int     `json:"id"`
	Name                 string  `json:"name"`
	SocialSecurityNumber string  `json:"socialSecurityNumber"`
	TaxReduction         float64 `json:"taxReducation"`
}

type InvoiceAttachment struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	SizeInBytes int    `json:"sizeInBytes"`
	ContactID   int    `json:"contactId"`
	// The API calls projects work orders.
	ProjectID     int       `json:"workOrderId"`
	IsTimeJournal bool      `json:"isTimeJournal"`
	SortIndex     int       `json:"sortIndex"`
	CreatedByUser *NameID   `json:"createdByUser"`
	CreatedDate   *DateTime `json:"createdDate"`
}

func (s *InvoicesService) List(ctx context.Context, params InvoiceListParams) (*Page[InvoiceListItem], error) {
	return listPage[InvoiceListItem](ctx, s.client, invoicesPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *InvoicesService) ListAll(ctx context.Context, params InvoiceListParams) iter.Seq2[InvoiceListItem, error] {
	return listAll[InvoiceListItem](ctx, s.client, invoicesPath, params)
}

func (s *InvoicesService) Get(ctx context.Context, id int) (*Invoice, error) {
	return getResource[Invoice](ctx, s.client, invoicesPath, strconv.Itoa(id))
}

// SetSentToEconomySystem leaves economySystemInvoiceNumber out of the request when it is empty.
func (s *InvoicesService) SetSentToEconomySystem(ctx context.Context, id int, date Date, economySystemInvoiceNumber string) error {
	return setInvoiceSentToEconomySystem(ctx, s.client, id, sentToEconomySystemQuery(date, economySystemInvoiceNumber))
}

func (s *InvoicesService) ListAttachments(ctx context.Context, invoiceID int, options ListOptions) (*Page[InvoiceAttachment], error) {
	return listInvoiceAttachments(ctx, s.client, invoiceID, options)
}

// ListAllAttachments walks every page from options.Page onwards; on failure it yields a single error and stops.
func (s *InvoicesService) ListAllAttachments(ctx context.Context, invoiceID int, options ListOptions) iter.Seq2[InvoiceAttachment, error] {
	return listAllInvoiceAttachments(ctx, s.client, invoiceID, options)
}

// DownloadAttachment streams the file; the caller must close Content.
func (s *InvoicesService) DownloadAttachment(ctx context.Context, invoiceID, attachmentID int) (*Download, error) {
	return downloadInvoiceAttachment(ctx, s.client, invoiceID, attachmentID)
}

func sentToEconomySystemQuery(date Date, economySystemInvoiceNumber string) url.Values {
	query := url.Values{"date": {date.String()}}
	if economySystemInvoiceNumber != "" {
		query.Set("economySystemInvoiceNumber", economySystemInvoiceNumber)
	}
	return query
}

func setInvoiceSentToEconomySystem(ctx context.Context, c *Client, id int, query url.Values) error {
	path, err := resourcePath(invoicesPath, strconv.Itoa(id), "setsenttoeconomysystem")
	if err != nil {
		return err
	}
	return patchResource(ctx, c, path, query, nil)
}

func invoiceAttachmentsPath(invoiceID int) (string, error) {
	return resourcePath(invoicesPath, strconv.Itoa(invoiceID), "attachments")
}

func listInvoiceAttachments(ctx context.Context, c *Client, invoiceID int, options ListOptions) (*Page[InvoiceAttachment], error) {
	path, err := invoiceAttachmentsPath(invoiceID)
	if err != nil {
		return nil, err
	}
	return listPage[InvoiceAttachment](ctx, c, path, options)
}

func listAllInvoiceAttachments(ctx context.Context, c *Client, invoiceID int, options ListOptions) iter.Seq2[InvoiceAttachment, error] {
	path, err := invoiceAttachmentsPath(invoiceID)
	if err != nil {
		return failedSeq[InvoiceAttachment](err)
	}
	return listAll[InvoiceAttachment](ctx, c, path, options)
}

func downloadInvoiceAttachment(ctx context.Context, c *Client, invoiceID, attachmentID int) (*Download, error) {
	path, err := resourcePath(invoicesPath, strconv.Itoa(invoiceID), "attachments", strconv.Itoa(attachmentID), "download")
	if err != nil {
		return nil, err
	}
	return c.download(ctx, path, nil)
}
