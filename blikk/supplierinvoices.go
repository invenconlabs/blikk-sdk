package blikk

import (
	"context"
	"iter"
	"strconv"
	"strings"
)

const supplierInvoicesPath = "v1/core/supplierinvoices"

type SupplierInvoicesService struct {
	client *Client
}

const (
	SupplierInvoiceSortByInvoiceNumber       = "invoiceNumber"
	SupplierInvoiceSortByInvoiceDate         = "invoiceDate"
	SupplierInvoiceSortByTotalAmount         = "totalAmount"
	SupplierInvoiceSortByInvoiceCostAddition = "invoiceCostAddition"
	SupplierInvoiceSortByCost                = "cost"
	SupplierInvoiceSortByCreatedDate         = "createdDate"
	SupplierInvoiceSortByUpdatedDate         = "updatedDate"
)

type SupplierInvoiceListParams struct {
	ListOptions
	Query                 string   `query:"filter.query"`
	InvoiceDateFrom       DateTime `query:"filter.invoiceDateFrom"`
	InvoiceDateTo         DateTime `query:"filter.invoiceDateTo"`
	CreatedDateFrom       DateTime `query:"filter.createdDateFrom"`
	CreatedDateTo         DateTime `query:"filter.createdDateTo"`
	UpdatedDateFrom       DateTime `query:"filter.updatedDateFrom"`
	UpdatedDateTo         DateTime `query:"filter.updatedDateTo"`
	CreatedOrUpdatedSince DateTime `query:"filter.createdOrUpdatedSince"`
	ProjectID             int      `query:"filter.projectId"`
	SupplierID            int      `query:"filter.supplierId"`
}

type SupplierInvoiceListItem struct {
	ID                  int        `json:"id"`
	InvoiceNumber       string     `json:"invoiceNumber"`
	InvoiceDate         *Date      `json:"invoiceDate"`
	Title               string     `json:"title"`
	Description         string     `json:"description"`
	ExternalHours       float64    `json:"externalHours"`
	TotalAmount         float64    `json:"totalAmount"`
	InvoiceCostAddition float64    `json:"invoiceCostAddition"`
	Cost                float64    `json:"cost"`
	OCRNumber           string     `json:"ocrNumber"`
	ToBeInvoiced        bool       `json:"toBeInvoiced"`
	InvoiceID           int        `json:"invoiceId"`
	InvoicedDate        *DateTime  `json:"invoicedDate"`
	InvoiceDraftID      int        `json:"invoiceDraftId"`
	Supplier            *NameID    `json:"supplier"`
	Project             ProjectRef `json:"project"`
	CreatedBy           *NameID    `json:"createdBy"`
	UpdatedBy           *NameID    `json:"updatedBy"`
	CreatedDate         DateTime   `json:"createdDate"`
	UpdatedDate         *DateTime  `json:"updatedDate"`
}

type SupplierInvoice struct {
	ID             int            `json:"id"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Supplier       *NameID        `json:"supplier"`
	InvoiceNumber  string         `json:"invoiceNumber"`
	InvoiceDate    *Date          `json:"invoiceDate"`
	ProjectCost    float64        `json:"projectCost"`
	Addition       float64        `json:"addition"`
	TotalAmount    float64        `json:"totalAmount"`
	ToBeInvoiced   bool           `json:"toBeInvoiced"`
	InvoiceID      int            `json:"invoiceId"`
	Project        ProjectRef     `json:"project"`
	CreatedBy      *NameID        `json:"createdBy"`
	UpdatedBy      *NameID        `json:"updatedBy"`
	CreatedDate    DateTime       `json:"createdDate"`
	UpdatedDate    *DateTime      `json:"updatedDate"`
	InvoicedDate   *DateTime      `json:"invoicedDate"`
	CostCenter     *CostCenterRef `json:"costCenter"`
	BookingProject *NameID        `json:"bookingProject"`
}

type SupplierInvoiceUnhandledInput struct {
	Title                string
	InvoiceDate          Date
	Amount               float64
	SourceSystemID       string
	Description          string
	SupplierName         string
	BookingProjectNumber string
	InvoiceNumber        string
	Files                []UploadFile
}

const supplierInvoiceFilesField = "files"

func (input SupplierInvoiceUnhandledInput) formFields() map[string]string {
	fields := map[string]string{
		"title":          input.Title,
		"invoiceDate":    input.InvoiceDate.String(),
		"amount":         formatDecimal(input.Amount),
		"sourceSystemId": input.SourceSystemID,
	}
	setFieldIfNotEmpty(fields, "description", input.Description)
	setFieldIfNotEmpty(fields, "supplierName", input.SupplierName)
	setFieldIfNotEmpty(fields, "bookingProjectNumber", input.BookingProjectNumber)
	setFieldIfNotEmpty(fields, "invoiceNumber", input.InvoiceNumber)
	return fields
}

func (input SupplierInvoiceUnhandledInput) uploadFiles() []UploadFile {
	files := make([]UploadFile, len(input.Files))
	for i, file := range input.Files {
		files[i] = withDefaultFieldName(file, supplierInvoiceFilesField)
	}
	return files
}

// The documented example sends "500.00"; more decimals are kept rather than rounded away.
func formatDecimal(value float64) string {
	formatted := strconv.FormatFloat(value, 'f', -1, 64)
	if decimalPlaces(formatted) < 2 {
		return strconv.FormatFloat(value, 'f', 2, 64)
	}
	return formatted
}

func decimalPlaces(formatted string) int {
	point := strings.IndexByte(formatted, '.')
	if point < 0 {
		return 0
	}
	return len(formatted) - point - 1
}

type SupplierInvoiceUpdateInput struct {
	Title            string   `json:"title"`
	Description      string   `json:"description,omitempty"`
	SupplierID       int      `json:"supplierId,omitempty"`
	InvoiceNumber    string   `json:"invoiceNumber,omitempty"`
	InvoiceDate      DateTime `json:"invoiceDate,omitzero"`
	ProjectCost      float64  `json:"projectCost,omitempty"`
	Addition         float64  `json:"addition,omitempty"`
	ToBeInvoiced     bool     `json:"toBeInvoiced"`
	ProjectID        int      `json:"projectId"`
	CostCenterID     int      `json:"costCenterId,omitempty"`
	BookingProjectID int      `json:"bookingProjectId,omitempty"`
}

func (s *SupplierInvoicesService) List(ctx context.Context, params SupplierInvoiceListParams) (*Page[SupplierInvoiceListItem], error) {
	return listPage[SupplierInvoiceListItem](ctx, s.client, supplierInvoicesPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *SupplierInvoicesService) ListAll(ctx context.Context, params SupplierInvoiceListParams) iter.Seq2[SupplierInvoiceListItem, error] {
	return listAll[SupplierInvoiceListItem](ctx, s.client, supplierInvoicesPath, params)
}

func (s *SupplierInvoicesService) Get(ctx context.Context, id int) (*SupplierInvoice, error) {
	return getResource[SupplierInvoice](ctx, s.client, supplierInvoicesPath, strconv.Itoa(id))
}

// AddUnhandled returns the ID of the new supplier invoice.
func (s *SupplierInvoicesService) AddUnhandled(ctx context.Context, input SupplierInvoiceUnhandledInput) (int, error) {
	return uploadResourceID(ctx, s.client, input.formFields(), input.uploadFiles(), supplierInvoicesPath, "addunhandled")
}

func (s *SupplierInvoicesService) Update(ctx context.Context, id int, input SupplierInvoiceUpdateInput) error {
	return patchSupplierInvoice(ctx, s.client, id, input)
}

func (s *SupplierInvoicesService) Delete(ctx context.Context, id int) error {
	return deleteResource(ctx, s.client, supplierInvoicesPath, strconv.Itoa(id))
}

func patchSupplierInvoice(ctx context.Context, c *Client, id int, input SupplierInvoiceUpdateInput) error {
	path, err := resourcePath(supplierInvoicesPath, strconv.Itoa(id))
	if err != nil {
		return err
	}
	return patchResource(ctx, c, path, nil, input)
}
