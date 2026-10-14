package blikk

import (
	"context"
	"iter"
	"strconv"
)

const (
	timeReportsPath       = "v1/Core/TimeReports"
	userDayStatisticsPath = "v1/Core/TimeReports/UserDayStatistics"
)

type TimeReportsService struct {
	client *Client
}

type TimeReportProjectType string

const (
	TimeReportProjectTypeProject         TimeReportProjectType = "Project"
	TimeReportProjectTypeInternalProject TimeReportProjectType = "InternalProject"
	TimeReportProjectTypeAbsenceProject  TimeReportProjectType = "AbsenceProject"
)

func (t *TimeReportProjectType) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, t, TimeReportProjectTypeProject, TimeReportProjectTypeInternalProject, TimeReportProjectTypeAbsenceProject)
}

type InvoiceStatus string

const (
	InvoiceStatusPending  InvoiceStatus = "Pending"
	InvoiceStatusDraft    InvoiceStatus = "Draft"
	InvoiceStatusInvoiced InvoiceStatus = "Invoiced"
)

func (s *InvoiceStatus) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, s, InvoiceStatusPending, InvoiceStatusDraft, InvoiceStatusInvoiced)
}

const (
	TimeReportSortByDate             = "date"
	TimeReportSortByUserName         = "userName"
	TimeReportSortByHours            = "hours"
	TimeReportSortByInvoiceableHours = "invoiceableHours"
	TimeReportSortBySentToAttestDate = "sentToAttestDate"
	TimeReportSortByAttestedDate     = "attestedDate"
	TimeReportSortByCreatedDate      = "createdDate"
	TimeReportSortByUpdatedDate      = "updatedDate"
)

// TimeReportListParams uses *bool for its boolean filters because false is a meaningful filter; leave them nil to not filter.
type TimeReportListParams struct {
	ListOptions
	CreatedDateFrom       DateTime              `query:"filter.createdDateFrom"`
	CreatedDateTo         DateTime              `query:"filter.createdDateTo"`
	UpdatedDateFrom       DateTime              `query:"filter.updatedDateFrom"`
	UpdatedDateTo         DateTime              `query:"filter.updatedDateTo"`
	CreatedOrUpdatedSince DateTime              `query:"filter.createdOrUpdatedSince"`
	From                  Date                  `query:"filter.from"`
	To                    Date                  `query:"filter.to"`
	UserIDs               []int                 `query:"filter.userIds"`
	ProjectID             int                   `query:"filter.projectId"`
	InternalProjectID     int                   `query:"filter.internalProjectId"`
	AbsenceProjectID      int                   `query:"filter.absenceProjectId"`
	ContactID             int                   `query:"filter.contactId"`
	ActivityID            int                   `query:"filter.activityId"`
	IsSentToAttest        *bool                 `query:"filter.isSentToAttest"`
	IsAttested            *bool                 `query:"filter.isAttested"`
	ProjectType           TimeReportProjectType `query:"filter.projectType"`
	HasAdditions          *bool                 `query:"filter.hasAdditions"`
	HasEquipment          *bool                 `query:"filter.hasEquipment"`
	HasAllowanceReports   *bool                 `query:"filter.hasAllowanceReports"`
	HasTravelReports      *bool                 `query:"filter.hasTravelReports"`
	TaskID                int                   `query:"filter.taskId"`
	InvoiceStatuses       []InvoiceStatus       `query:"filter.invoiceStatus"`
}

type TimeReportListItem struct {
	ID                int            `json:"id"`
	Date              Date           `json:"date"`
	ClockStart        string         `json:"clockStart"`
	ClockEnd          string         `json:"clockEnd"`
	Hours             float64        `json:"hours"`
	InvoiceableHours  float64        `json:"invoiceableHours"`
	BreakMinutes      int            `json:"breakMinutes"`
	Cost              float64        `json:"cost"`
	Rate              float64        `json:"rate"`
	Discount          float64        `json:"discount"`
	Comment           string         `json:"comment"`
	InternalComment   string         `json:"internalComment"`
	SentToAttestDate  *DateTime      `json:"sentToAttestDate"`
	AttestedDate      *DateTime      `json:"attestedDate"`
	User              NameID         `json:"user"`
	Project           *ProjectRef    `json:"project"`
	InternalProject   *NameID        `json:"internalProject"`
	AbsenceProject    *NameID        `json:"absenceProject"`
	Contact           NameID         `json:"contact"`
	Activity          NameID         `json:"activity"`
	TimeCode          NameID         `json:"timeCode"`
	TimeArticle       NameID         `json:"timeArticle"`
	CostCenter        *CostCenterRef `json:"costCenter"`
	InvoiceID         int            `json:"invoiceId"`
	InvoicedDate      *DateTime      `json:"invoicedDate"`
	InvoiceDraftID    int            `json:"invoiceDraftId"`
	TravelReportID    int            `json:"travelReportId"`
	AllowanceReportID int            `json:"allowanceReportId"`
	CreatedDate       DateTime       `json:"createdDate"`
	UpdatedDate       *DateTime      `json:"updatedDate"`
	HasAdditions      bool           `json:"hasAdditions"`
	HasEquipment      bool           `json:"hasEquipment"`
	CreatedBy         *NameID        `json:"createdBy"`
	UpdatedBy         *NameID        `json:"updatedBy"`
	TaskID            int            `json:"taskId"`
}

type TimeReport struct {
	ID               int                  `json:"id"`
	Date             Date                 `json:"date"`
	ClockStart       string               `json:"clockStart"`
	ClockEnd         string               `json:"clockEnd"`
	Hours            float64              `json:"hours"`
	InvoiceableHours float64              `json:"invoiceableHours"`
	BreakMinutes     int                  `json:"breakMinutes"`
	Cost             float64              `json:"cost"`
	Rate             float64              `json:"rate"`
	Discount         float64              `json:"discount"`
	Comment          string               `json:"comment"`
	InternalComment  string               `json:"internalComment"`
	SentToAttestDate *DateTime            `json:"sentToAttestDate"`
	AttestedDate     *DateTime            `json:"attestedDate"`
	AttestedByUser   *NameID              `json:"attestedByUser"`
	User             NameID               `json:"user"`
	Contact          NameID               `json:"contact"`
	Project          *ProjectRef          `json:"project"`
	InternalProject  *NameID              `json:"internalProject"`
	AbsenceProject   *NameID              `json:"absenceProject"`
	Activity         NameID               `json:"activity"`
	TimeCode         NameID               `json:"timeCode"`
	TimeArticle      NameID               `json:"timeArticle"`
	CostCenter       *CostCenterRef       `json:"costCenter"`
	BookingProject   *NameID              `json:"bookingProject"`
	InvoiceID        int                  `json:"invoiceId"`
	InvoicedDate     *DateTime            `json:"invoicedDate"`
	InvoiceDraftID   int                  `json:"invoiceDraftId"`
	Task             *NameID              `json:"task"`
	TravelReport     *TravelReport        `json:"travelReport"`
	AllowanceReport  *AllowanceReport     `json:"allowanceReport"`
	Materials        []TimeReportMaterial `json:"materials"`
	Equipments       []EquipmentReport    `json:"equipments"`
	Additions        []AdditionReport     `json:"additions"`
	CreatedDate      DateTime             `json:"createdDate"`
	UpdatedDate      *DateTime            `json:"updatedDate"`
}

type TravelReport struct {
	ID                  int                     `json:"id"`
	Place               string                  `json:"place"`
	InvoiceablePlace    string                  `json:"invoicablePlace"`
	Distance            float64                 `json:"distance"`
	InvoiceableDistance float64                 `json:"invoicableDistance"`
	ToSalary            bool                    `json:"toSalary"`
	ToBeInvoiced        bool                    `json:"toBeInvoiced"`
	Cost                float64                 `json:"cost"`
	Rate                float64                 `json:"rate"`
	Discount            float64                 `json:"discount"`
	SentToAttestDate    *DateTime               `json:"sentToAttestDate"`
	AttestedDate        *DateTime               `json:"attestedDate"`
	CostCenter          *CostCenterRef          `json:"costCenter"`
	BookingProject      *NameID                 `json:"bookingProject"`
	InvoiceID           int                     `json:"invoiceId"`
	InvoicedDate        *DateTime               `json:"invoicedDate"`
	InvoiceDraftID      int                     `json:"invoiceDraftId"`
	CompanyCar          *TravelReportCompanyCar `json:"companyCar"`
	Article             *NameID                 `json:"article"`
}

type TravelReportCompanyCar struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	TripStart    float64 `json:"tripStart"`
	TripEnd      float64 `json:"tripEnd"`
	AddressStart string  `json:"addressStart"`
	AddressGoal  string  `json:"addressGoal"`
	AddressEnd   string  `json:"addressEnd"`
}

type AllowanceReport struct {
	ID                       int            `json:"id"`
	Description              string         `json:"description"`
	Cost                     float64        `json:"cost"`
	Amount                   float64        `json:"amount"`
	InvoiceableAmount        float64        `json:"invoicableAmount"`
	Discount                 float64        `json:"discount"`
	InvoiceableSpecification string         `json:"invoicableSpecification"`
	ToBeInvoiced             bool           `json:"toBeInvoiced"`
	SentToAttestDate         *DateTime      `json:"sentToAttestDate"`
	AttestedDate             *DateTime      `json:"attestedDate"`
	CostCenter               *CostCenterRef `json:"costCenter"`
	BookingProject           *NameID        `json:"bookingProject"`
	InvoiceID                int            `json:"invoiceId"`
	InvoicedDate             *DateTime      `json:"invoicedDate"`
	InvoiceDraftID           int            `json:"invoiceDraftId"`
	Article                  *NameID        `json:"article"`
}

type TimeReportMaterial struct {
	ID                   int            `json:"id"`
	Title                string         `json:"title"`
	Units                float64        `json:"units"`
	CostPerUnit          float64        `json:"costPerUnit"`
	DiscountInPercent    float64        `json:"discountInPercent"`
	InvoiceableUnits     float64        `json:"invoicableUnits"`
	UnitPrice            float64        `json:"unitPrice"`
	InvoiceableUnitPrice float64        `json:"invoicableUnitPrice"`
	TaxPercent           float64        `json:"taxPercent"`
	Amount               float64        `json:"amount"`
	InvoiceableAmount    float64        `json:"invoicableAmount"`
	Cost                 float64        `json:"cost"`
	CalculatedUnits      float64        `json:"calculatedUnits"`
	Tax                  float64        `json:"tax"`
	ToBeInvoiced         bool           `json:"toBeInvoiced"`
	Article              *NameID        `json:"article"`
	UnitText             string         `json:"unitText"`
	InvoiceableUnitText  string         `json:"invoicableUnitText"`
	IsRotRut             bool           `json:"isRotRut"`
	InvoiceID            int            `json:"invoiceId"`
	InvoicedDate         *DateTime      `json:"invoicedDate"`
	InvoiceDraftID       int            `json:"invoiceDraftId"`
	ArticleNumber        string         `json:"articleNumber"`
	CostCenter           *CostCenterRef `json:"costCenter"`
	BookingProject       *NameID        `json:"bookingProject"`
}

type EquipmentReport struct {
	ID               int            `json:"id"`
	Hours            float64        `json:"hours"`
	Cost             float64        `json:"cost"`
	Rate             float64        `json:"rate"`
	InvoiceableRate  float64        `json:"invoicableRate"`
	InvoiceableHours float64        `json:"invoicableHours"`
	Article          NameID         `json:"article"`
	InvoiceID        int            `json:"invoiceId"`
	InvoicedDate     *DateTime      `json:"invoicedDate"`
	InvoiceDraftID   int            `json:"invoiceDraftId"`
	Discount         float64        `json:"discount"`
	Comment          string         `json:"comment"`
	CostCenter       *CostCenterRef `json:"costCenter"`
	BookingProject   *NameID        `json:"bookingProject"`
}

type AdditionReport struct {
	ID               int            `json:"id"`
	SalaryCount      float64        `json:"salaryCount"`
	InvoiceableCount float64        `json:"invoicableCount"`
	Cost             float64        `json:"cost"`
	InvoiceablePrice float64        `json:"invoicablePrice"`
	Discount         float64        `json:"discount"`
	InvoiceID        int            `json:"invoiceId"`
	InvoicedDate     *DateTime      `json:"invoicedDate"`
	InvoiceDraftID   int            `json:"invoiceDraftId"`
	Article          NameID         `json:"article"`
	CostCenter       *CostCenterRef `json:"costCenter"`
	BookingProject   *NameID        `json:"bookingProject"`
}

type TimeReportInput struct {
	Date              Date                  `json:"date"`
	ClockStart        string                `json:"clockStart,omitempty"`
	ClockEnd          string                `json:"clockEnd,omitempty"`
	Hours             float64               `json:"hours"`
	InvoiceableHours  float64               `json:"invoiceableHours"`
	BreakMinutes      int                   `json:"breakMinutes,omitempty"`
	Comment           string                `json:"comment,omitempty"`
	InternalComment   string                `json:"internalComment,omitempty"`
	UserID            int                   `json:"userId"`
	ProjectID         int                   `json:"projectId,omitempty"`
	InternalProjectID int                   `json:"internalProjectId,omitempty"`
	AbsenceProjectID  int                   `json:"absenceProjectId,omitempty"`
	ActivityID        int                   `json:"activityId"`
	TimeCodeID        int                   `json:"timeCodeId"`
	TimeArticleID     int                   `json:"timeArticleId"`
	TaskID            int                   `json:"taskId,omitempty"`
	TravelReport      *TravelReportInput    `json:"travelReport,omitempty"`
	AllowanceReport   *AllowanceReportInput `json:"allowanceReport,omitempty"`
	Materials         []MaterialInput       `json:"materials,omitempty"`
	Equipments        []EquipmentInput      `json:"equipments,omitempty"`
	Additions         []AdditionInput       `json:"additions,omitempty"`
}

type TravelReportInput struct {
	Place               string           `json:"place"`
	Distance            float64          `json:"distance"`
	InvoiceableDistance float64          `json:"invoiceableDistance,omitempty"`
	ToSalary            bool             `json:"toSalary"`
	CompanyCar          *CompanyCarInput `json:"companyCar,omitempty"`
}

type CompanyCarInput struct {
	ID           int     `json:"id"`
	TripStart    float64 `json:"tripStart,omitempty"`
	TripEnd      float64 `json:"tripEnd,omitempty"`
	AddressStart string  `json:"addressStart,omitempty"`
	AddressGoal  string  `json:"addressGoal,omitempty"`
	AddressEnd   string  `json:"addressEnd,omitempty"`
}

type AllowanceReportInput struct {
	AllowanceTypeID int    `json:"allowanceTypeId"`
	Description     string `json:"description,omitempty"`
	ToBeInvoiced    bool   `json:"toBeInvoiced"`
}

type MaterialInput struct {
	ID                   int     `json:"id,omitempty"`
	ArticleID            int     `json:"articleId,omitempty"`
	ArticleNumber        string  `json:"articleNumber,omitempty"`
	Title                string  `json:"title"`
	UnitPrice            float64 `json:"unitPrice"`
	InvoiceableUnitPrice float64 `json:"invoiceableUnitPrice,omitempty"`
	Units                float64 `json:"units"`
	InvoiceableUnits     float64 `json:"invoiceableUnits,omitempty"`
	CalculatedUnits      float64 `json:"calculatedUnits,omitempty"`
	UnitText             string  `json:"unitText,omitempty"`
	Cost                 float64 `json:"cost,omitempty"`
	DiscountInPercent    float64 `json:"discountInPercent,omitempty"`
	IsTextOnly           bool    `json:"isTextOnly"`
}

type EquipmentInput struct {
	ID        int     `json:"id,omitempty"`
	Hours     float64 `json:"hours"`
	ArticleID int     `json:"articleId"`
	Comment   string  `json:"comment,omitempty"`
}

type AdditionInput struct {
	ID          int     `json:"id,omitempty"`
	ArticleID   int     `json:"articleId"`
	YesNo       bool    `json:"yesNo"`
	NumberCount float64 `json:"numberCount,omitempty"`
}

const UserDayStatisticsSortByName = "name"

type UserDayStatisticsListParams struct {
	ListOptions
	From          Date  `query:"filter.from"`
	To            Date  `query:"filter.to"`
	UserIDs       []int `query:"filter.userIds"`
	DepartmentIDs []int `query:"filter.departmentIds"`
}

type UserDayStatistics struct {
	UserID     int                      `json:"userId"`
	Name       string                   `json:"name"`
	Department *NameID                  `json:"department"`
	Dates      []UserDayStatisticsEntry `json:"dates"`
}

type UserDayStatisticsEntry struct {
	Date           Date      `json:"date"`
	ReportedHours  float64   `json:"reportedHours"`
	ScheduledHours float64   `json:"scheduledHours"`
	LockedDate     *DateTime `json:"lockedDate"`
	AttestedDate   *DateTime `json:"attestedDate"`
}

func (r TimeReportListItem) ProjectName() string {
	return firstProjectName(r.Project, r.InternalProject, r.AbsenceProject)
}

func (r TimeReport) ProjectName() string {
	return firstProjectName(r.Project, r.InternalProject, r.AbsenceProject)
}

func firstProjectName(project *ProjectRef, internalProject, absenceProject *NameID) string {
	switch {
	case project != nil:
		return project.Name
	case internalProject != nil:
		return internalProject.Name
	case absenceProject != nil:
		return absenceProject.Name
	default:
		return ""
	}
}

func (s *TimeReportsService) List(ctx context.Context, params TimeReportListParams) (*Page[TimeReportListItem], error) {
	return listPage[TimeReportListItem](ctx, s.client, timeReportsPath, params)
}

func (s *TimeReportsService) ListAll(ctx context.Context, params TimeReportListParams) iter.Seq2[TimeReportListItem, error] {
	return listAll[TimeReportListItem](ctx, s.client, timeReportsPath, params)
}

func (s *TimeReportsService) Get(ctx context.Context, id int) (*TimeReport, error) {
	return getResource[TimeReport](ctx, s.client, timeReportsPath, strconv.Itoa(id))
}

// Create returns the ID of the new time report.
func (s *TimeReportsService) Create(ctx context.Context, input TimeReportInput) (int, error) {
	return createResourceID(ctx, s.client, input, timeReportsPath)
}

// Update replaces the whole time report; any field left unset in input is cleared.
func (s *TimeReportsService) Update(ctx context.Context, id int, input TimeReportInput) error {
	return updateResource(ctx, s.client, input, timeReportsPath, strconv.Itoa(id))
}

func (s *TimeReportsService) Delete(ctx context.Context, id int) error {
	return deleteResource(ctx, s.client, timeReportsPath, strconv.Itoa(id))
}

// ListUserDayStatistics returns one page; the API rejects ranges longer than 31 days.
func (s *TimeReportsService) ListUserDayStatistics(ctx context.Context, params UserDayStatisticsListParams) (*Page[UserDayStatistics], error) {
	return listPage[UserDayStatistics](ctx, s.client, userDayStatisticsPath, params)
}

func (s *TimeReportsService) ListAllUserDayStatistics(ctx context.Context, params UserDayStatisticsListParams) iter.Seq2[UserDayStatistics, error] {
	return listAll[UserDayStatistics](ctx, s.client, userDayStatisticsPath, params)
}
