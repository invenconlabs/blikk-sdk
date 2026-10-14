package blikk

import (
	"context"
	"iter"
	"net/http"
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

type ProjectWorkSiteAddressInput struct {
	StreetAddress string  `json:"streetAddress,omitempty"`
	PostalCode    string  `json:"postalCode,omitempty"`
	City          string  `json:"city,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`
	Latitude      float64 `json:"latitude,omitempty"`
}

type ProjectCreateInput struct {
	ProjectNumber            string                       `json:"projectNumber,omitempty"`
	Title                    string                       `json:"title"`
	Description              string                       `json:"description,omitempty"`
	ContactID                int                          `json:"contactId"`
	ResponsibleID            int                          `json:"responsibleId"`
	StatusID                 int                          `json:"statusId,omitempty"`
	ProjectCollectionID      int                          `json:"projectCollectionId,omitempty"`
	StartDate                DateTime                     `json:"startDate,omitzero"`
	EndDate                  DateTime                     `json:"endDate,omitzero"`
	WorkSiteAddress          *ProjectWorkSiteAddressInput `json:"workSiteAddress,omitempty"`
	CategoryID               int                          `json:"categoryId,omitempty"`
	CostCenterID             int                          `json:"costCenterId,omitempty"`
	InvoiceText              string                       `json:"invoiceText,omitempty"`
	YourReference            string                       `json:"yourReference,omitempty"`
	OurReference             string                       `json:"ourReference,omitempty"`
	CustomerReferenceMarking string                       `json:"customerReferenceMarking,omitempty"`
	AllIsMember              bool                         `json:"allIsMember,omitempty"`
}

type ProjectCreateResult struct {
	ID                int    `json:"id"`
	NextProjectNumber string `json:"nextProjectNumber"`
}

// ProjectUpdateInput.ID is filled in by Update from its id argument.
type ProjectUpdateInput struct {
	ID                       int                          `json:"id"`
	Title                    string                       `json:"title,omitempty"`
	ProjectCollectionID      int                          `json:"projectCollectionId,omitempty"`
	StartDate                DateTime                     `json:"startDate,omitzero"`
	EndDate                  DateTime                     `json:"endDate,omitzero"`
	TagIDs                   []int                        `json:"tagIds,omitempty"`
	ContactID                int                          `json:"contactId,omitempty"`
	WorkSiteAddress          *ProjectWorkSiteAddressInput `json:"workSiteAddress,omitempty"`
	ResponsibleID            int                          `json:"responsibleId,omitempty"`
	SalesResponsibleID       int                          `json:"salesResponsibleId,omitempty"`
	Description              string                       `json:"description,omitempty"`
	StatusID                 int                          `json:"statusId,omitempty"`
	ReadyForInvoice          bool                         `json:"readyForInvoice,omitempty"`
	CostCenterID             int                          `json:"costCenterId,omitempty"`
	BookingProjectID         int                          `json:"bookingProjectId,omitempty"`
	CategoryID               int                          `json:"categoryId,omitempty"`
	InvoiceText              string                       `json:"invoiceText,omitempty"`
	YourReference            string                       `json:"yourReference,omitempty"`
	OurReference             string                       `json:"ourReference,omitempty"`
	CustomerReferenceMarking string                       `json:"customerReferenceMarking,omitempty"`
	AllIsMember              bool                         `json:"allIsMember,omitempty"`
}

type projectNextNumber struct {
	NextProjectNumber string `json:"nextProjectNumber"`
}

const (
	ProjectCommentSortByCreatedDate = "createdDate"
	ProjectCommentSortByUpdatedDate = "updatedDate"
)

type ProjectCommentListParams struct {
	ListOptions
	CreatedFrom      DateTime `query:"filter.createdFrom"`
	CreatedTo        DateTime `query:"filter.createdTo"`
	CreatedByUserIDs []int    `query:"filter.createdByUserIds"`
}

type projectCommentInput struct {
	Text string `json:"text"`
}

type ProjectTimeCalculation struct {
	Total                    float64                          `json:"total"`
	IsTotal                  bool                             `json:"isTotal"`
	IsByActivity             bool                             `json:"isByActivity"`
	IsByActivityAndUser      bool                             `json:"isByActivityAndUser"`
	ActivityCalculations     []ProjectActivityCalculation     `json:"activityCalculations"`
	UserActivityCalculations []ProjectUserActivityCalculation `json:"userActivityCalculations"`
}

type ProjectActivityCalculation struct {
	Activity    *NameID `json:"activity"`
	Calculation float64 `json:"calculation"`
}

// The docs give no example of this shape and the live account has none; the fields are a guess modelled on ProjectActivityCalculation.
type ProjectUserActivityCalculation struct {
	User        *NameID `json:"user"`
	Activity    *NameID `json:"activity"`
	Calculation float64 `json:"calculation"`
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

func (s *ProjectsService) Create(ctx context.Context, input ProjectCreateInput) (*ProjectCreateResult, error) {
	return createResource[ProjectCreateResult](ctx, s.client, input, projectsPath)
}

func (s *ProjectsService) Update(ctx context.Context, id int, input ProjectUpdateInput) error {
	input.ID = id
	return updateResource(ctx, s.client, input, projectsPath, strconv.Itoa(id))
}

// NextProjectNumber fails with a 400 when suggesting project numbers is turned off in the account settings.
func (s *ProjectsService) NextProjectNumber(ctx context.Context) (string, error) {
	next, err := getResource[projectNextNumber](ctx, s.client, projectsPath, "nextProjectNumber")
	if err != nil {
		return "", err
	}
	return next.NextProjectNumber, nil
}

func (s *ProjectsService) ListComments(ctx context.Context, projectID int, params ProjectCommentListParams) (*Page[Comment], error) {
	return listPage[Comment](ctx, s.client, projectCommentsPath(projectID), params)
}

// ListAllComments walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *ProjectsService) ListAllComments(ctx context.Context, projectID int, params ProjectCommentListParams) iter.Seq2[Comment, error] {
	return listAll[Comment](ctx, s.client, projectCommentsPath(projectID), params)
}

func (s *ProjectsService) CreateComment(ctx context.Context, projectID int, text string) (int, error) {
	return createResourceID(ctx, s.client, projectCommentInput{Text: text}, projectCommentsPath(projectID))
}

func (s *ProjectsService) UpdateComment(ctx context.Context, projectID, commentID int, text string) error {
	return updateResource(ctx, s.client, projectCommentInput{Text: text}, projectCommentsPath(projectID), strconv.Itoa(commentID))
}

func (s *ProjectsService) DeleteComment(ctx context.Context, projectID, commentID int) error {
	return deleteResource(ctx, s.client, projectCommentsPath(projectID), strconv.Itoa(commentID))
}

func (s *ProjectsService) TimeCalculation(ctx context.Context, projectID int) (*ProjectTimeCalculation, error) {
	return getResource[ProjectTimeCalculation](ctx, s.client, projectsPath, strconv.Itoa(projectID), "TimeCalculation")
}

// AddMember fails with a 409 when the user already is a member.
func (s *ProjectsService) AddMember(ctx context.Context, projectID, userID int) error {
	return s.client.do(ctx, http.MethodPost, projectMemberPath(projectID, userID), nil, nil, nil)
}

func (s *ProjectsService) RemoveMember(ctx context.Context, projectID, userID int) error {
	return deleteResource(ctx, s.client, projectMemberPath(projectID, userID))
}

func projectCommentsPath(projectID int) string {
	return projectsPath + "/" + strconv.Itoa(projectID) + "/Comments"
}

func projectMemberPath(projectID, userID int) string {
	return projectsPath + "/" + strconv.Itoa(projectID) + "/members/" + strconv.Itoa(userID)
}
