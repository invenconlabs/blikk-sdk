package blikk

import (
	"context"
	"iter"
	"strconv"
)

const projectsPath = "v1/Core/Projects"

type ProjectsService struct {
	client *Client
}

type InvoiceType string

const (
	InvoiceTypeOngoing    InvoiceType = "Ongoing"
	InvoiceTypeFixedPrice InvoiceType = "FixedPrice"
	InvoiceTypeNone       InvoiceType = "None"
)

func (t *InvoiceType) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, t, InvoiceTypeOngoing, InvoiceTypeFixedPrice, InvoiceTypeNone)
}

const (
	ProjectSortByOrderNumber = "orderNumber"
	ProjectSortByTitle       = "title"
	ProjectSortByCreatedDate = "createdDate"
	ProjectSortByUpdatedDate = "updatedDate"
)

type ProjectListParams struct {
	ListOptions
	CreatedDateFrom       DateTime    `query:"filter.createdDateFrom"`
	CreatedDateTo         DateTime    `query:"filter.createdDateTo"`
	UpdatedDateFrom       DateTime    `query:"filter.updatedDateFrom"`
	UpdatedDateTo         DateTime    `query:"filter.updatedDateTo"`
	CreatedOrUpdatedSince DateTime    `query:"filter.createdOrUpdatedSince"`
	Query                 string      `query:"filter.query"`
	ProjectTagIDs         []int       `query:"filter.projectTagIds"`
	ProjectCategoryIDs    []int       `query:"filter.projectCategoryIds"`
	ProjectStatusIDs      []int       `query:"filter.projectStatusIds"`
	ProjectType           InvoiceType `query:"filter.projectType"`
}

type ProjectListItem struct {
	ID                int              `json:"id"`
	OrderNumber       string           `json:"orderNumber"`
	Title             string           `json:"title"`
	Status            *ProjectStatus   `json:"status"`
	Category          *ProjectCategory `json:"category"`
	SalesResponsible  *NameID          `json:"salesResponsible"`
	StartDate         *Date            `json:"startDate"`
	EndDate           *Date            `json:"endDate"`
	InvoiceType       InvoiceType      `json:"invoiceType"`
	Location          *Location        `json:"location"`
	ProjectManager    *NameID          `json:"projectManager"`
	Customer          *NameID          `json:"customer"`
	ProjectCollection *ProjectRef      `json:"projectCollection"`
	Tags              []Tag            `json:"tags"`
	CostCenter        *CostCenterRef   `json:"costCenter"`
	CreatedBy         *NameID          `json:"createdBy"`
	UpdatedBy         *NameID          `json:"updatedBy"`
	Created           *DateTime        `json:"created"`
	Updated           *DateTime        `json:"updated"`
}

type Project struct {
	ID                             int              `json:"id"`
	OrderNumber                    string           `json:"orderNumber"`
	Title                          string           `json:"title"`
	Status                         *ProjectStatus   `json:"status"`
	ReadyForInvoice                bool             `json:"readyForInvoice"`
	Category                       *ProjectCategory `json:"category"`
	InvoiceText                    string           `json:"invoiceText"`
	YourReference                  string           `json:"yourReference"`
	OurReference                   string           `json:"ourReference"`
	CustomerReferenceMarking       string           `json:"customerReferenceMarking"`
	SalesResponsible               *NameID          `json:"salesResponsible"`
	Description                    string           `json:"description"`
	StartDate                      *Date            `json:"startDate"`
	EndDate                        *Date            `json:"endDate"`
	InvoiceType                    InvoiceType      `json:"invoiceType"`
	Location                       *Location        `json:"location"`
	ProjectManager                 *NameID          `json:"projectManager"`
	Customer                       *NameID          `json:"customer"`
	ProjectCollection              *ProjectRef      `json:"projectCollection"`
	Tags                           []Tag            `json:"tags"`
	CostCenter                     *CostCenterRef   `json:"costCenter"`
	BookingProject                 *NameID          `json:"bookingProject"`
	CreatedBy                      *NameID          `json:"createdBy"`
	UpdatedBy                      *NameID          `json:"updatedBy"`
	Created                        *DateTime        `json:"created"`
	Updated                        *DateTime        `json:"updated"`
	CompletedDate                  *DateTime        `json:"completedDate"`
	ElectronicInvoiceReferenceCode string           `json:"electronicInvoiceReferenceCode"`
	ImportantInvoiceInformation    string           `json:"importantInvoiceInformation"`
	AllIsMember                    bool             `json:"allIsMember"`
	Members                        []NameID         `json:"members"`
}

func (s *ProjectsService) List(ctx context.Context, params ProjectListParams) (*Page[ProjectListItem], error) {
	return listPage[ProjectListItem](ctx, s.client, projectsPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *ProjectsService) ListAll(ctx context.Context, params ProjectListParams) iter.Seq2[ProjectListItem, error] {
	return listAll[ProjectListItem](ctx, s.client, projectsPath, params)
}

func (s *ProjectsService) Get(ctx context.Context, id int) (*Project, error) {
	return getResource[Project](ctx, s.client, projectsPath, strconv.Itoa(id))
}

func (s *ProjectsService) GetByOrderNumber(ctx context.Context, orderNumber string) (*Project, error) {
	return getResource[Project](ctx, s.client, projectsPath, "order", orderNumber)
}
