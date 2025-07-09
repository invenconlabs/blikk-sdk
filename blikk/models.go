package blikk

import (
	"time"

	"github.com/invenconlabs/blikk-sdk/dateutils"
)

type accessTokenResponse struct {
	ObjectName  string `json:"objectName"`
	AccessToken string `json:"accessToken"`
	Expires     string `json:"expires"`
}

type ListOptions struct {
	// Pagination
	Page     int `paramName:"page"`
	PageSize int `paramName:"pageSize"`

	// Filtering
	UserIDs  []uint16            `paramName:"filter.userIds"`
	FromDate *dateutils.DateOnly `paramName:"filter.from"`
	ToDate   *dateutils.DateOnly `paramName:"filter.to"`
}

func NewListOptions() ListOptions {
	return ListOptions{
		Page:     1,
		PageSize: 100,
	}
}

type blikkObject struct {
	ObjectName string `json:"objectName"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
}

type ListItem interface {
	path() string
	validFilter(options *ListOptions) bool
}

type ListResponse[T ListItem] struct {
	ObjectName     string `json:"objectName"`
	Page           int    `json:"page"`
	PageSize       int    `json:"pageSize"`
	ItemCount      int    `json:"itemCount"`
	TotalItemCount int    `json:"totalItemCount"`
	TotalPages     int    `json:"totalPages"`
	Items          []T    `json:"items"`
}

type Users struct {
	ObjectName           string             `json:"objectName"`
	ID                   int                `json:"id"`
	FirstName            string             `json:"firstName"`
	LastName             string             `json:"lastName"`
	Email                string             `json:"email"`
	StartDate            dateutils.DateOnly `json:"startDate"`
	EndDate              dateutils.DateOnly `json:"endDate,omitempty"`
	EmployeeNumber       string             `json:"employeeNumber"`
	TimeReportingProfile blikkObject        `json:"timeReportingProfile"`
	Department           blikkObject        `json:"department"`
	License              string             `json:"license"`
	Permissions          []string           `json:"permissions"`
	CreatedBy            blikkObject        `json:"createdBy"`
	UpdatedBy            blikkObject        `json:"updatedBy"`
	CreatedDate          string             `json:"createdDate"`
	UpdatedDate          string             `json:"updatedDate"`
	SalaryType           int                `json:"salaryType"`
	CostCenter           blikkObject        `json:"costCenter"`
}

func (Users) path() string {
	return "v1/Admin/Users"
}

func (Users) validFilter(options *ListOptions) bool {
	if len(options.UserIDs) > 0 ||
		options.FromDate != nil ||
		options.ToDate != nil {
		return false
	}
	return true
}

type UserDayStatistics struct {
	ObjectName string      `json:"objectName"`
	UserID     uint16      `json:"userId"`
	Name       string      `json:"name"`
	Department blikkObject `json:"department"`
	Dates      []struct {
		ObjectName     string             `json:"objectName"`
		Date           dateutils.DateOnly `json:"date"`
		ReportedHours  float64            `json:"reportedHours"`
		ScheduledHours float64            `json:"scheduledHours"`
		LockedDate     string             `json:"lockedDate"`
		AttestedDate   string             `json:"attestedDate"`
	} `json:"dates"`
}

func (UserDayStatistics) path() string {
	return "v1/Core/TimeReports/UserDayStatistics"
}

func (UserDayStatistics) validFilter(options *ListOptions) bool {
	if options.FromDate != nil && options.ToDate != nil {
		if options.FromDate.After(options.ToDate.Time) {
			return false
		}

		if options.ToDate.Sub(options.FromDate.Time) > 31*24*time.Hour {
			return false
		}
	}

	return true
}

type TimeReports struct {
	ObjectName       string             `json:"objectName"`
	ID               int                `json:"id"`
	Date             dateutils.DateOnly `json:"date"`
	ClockStart       string             `json:"clockStart"`
	ClockEnd         string             `json:"clockEnd"`
	Hours            float64            `json:"hours"`
	InvoiceableHours float64            `json:"invoiceableHours"`
	BreakMinutes     int                `json:"breakMinutes"`
	Cost             float64            `json:"cost"`
	Rate             float64            `json:"rate"`
	Discount         float64            `json:"discount"`
	Comment          string             `json:"comment"`
	InternalComment  string             `json:"internalComment"`
	SentToAttestDate string             `json:"sentToAttestDate"`
	AttestedDate     string             `json:"attestedDate"`
	User             blikkObject        `json:"user"`
	Project          struct {
		blikkObject
		Number string `json:"number"`
	} `json:"project"`
	InternalProject   blikkObject `json:"internalProject"`
	AbsenceProject    blikkObject `json:"absenceProject"`
	Contact           blikkObject `json:"contact"`
	Activity          blikkObject `json:"activity"`
	TimeCode          blikkObject `json:"timeCode"`
	TimeArticle       blikkObject `json:"timeArticle"`
	CostCenter        blikkObject `json:"costCenter"`
	InvoiceID         int         `json:"invoiceId"`
	InvoicedDate      string      `json:"invoicedDate"`
	InvoiceDraftID    int         `json:"invoiceDraftId"`
	TravelReportID    int         `json:"travelReportId"`
	AllowanceReportID int         `json:"allowanceReportId"`
	CreatedDate       string      `json:"createdDate"`
	UpdatedDate       string      `json:"updatedDate"`
	HasAdditions      bool        `json:"hasAdditions"`
	HasEquipment      bool        `json:"hasEquipment"`
	CreatedBy         blikkObject `json:"createdBy"`
	UpdatedBy         blikkObject `json:"updatedBy"`
	TaskID            int         `json:"taskId"`
}

func (TimeReports) path() string {
	return "v1/Core/TimeReports"
}

func (TimeReports) validFilter(options *ListOptions) bool {
	if options.FromDate != nil && options.ToDate != nil {
		if options.FromDate.After(options.ToDate.Time) {
			return false
		}
	}
	return true
}

type Projects struct {
	ObjectName  string `json:"objectName"`
	ID          int    `json:"id"`
	OrderNumber string `json:"orderNumber"`
	Title       string `json:"title"`
	Status      struct {
		blikkObject
		IsCompletedStatus bool `json:"isCompletedStatus"`
	} `json:"status"`
	Category struct {
		blikkObject
		Color string `json:"color"`
	} `json:"category"`
	SalesResponsible blikkObject        `json:"salesResponsible"`
	StartDate        dateutils.DateOnly `json:"startDate"`
	EndDate          dateutils.DateOnly `json:"endDate"`
	InvoiceType      string             `json:"invoiceType"`
	Location         struct {
		ObjectName    string  `json:"objectName"`
		Longitude     float64 `json:"longitude"`
		Latitude      float64 `json:"latitude"`
		StreetAddress string  `json:"streetAddress"`
		PostalCode    string  `json:"postalCode"`
		City          string  `json:"city"`
		CountryName   string  `json:"countryName"`
	} `json:"location"`
	ProjectManager    blikkObject `json:"projectManager"`
	Customer          blikkObject `json:"customer"`
	ProjectCollection struct {
		blikkObject
		Number string `json:"number"`
	} `json:"projectCollection"`
	Tags []struct {
		ObjectName string `json:"objectName"`
		ID         int    `json:"id"`
		Title      string `json:"title"`
		Color      string `json:"color"`
	} `json:"tags"`
	CostCenter struct {
		blikkObject
		Code string `json:"code"`
	} `json:"costCenter"`
	CreatedBy blikkObject `json:"createdBy"`
	UpdatedBy blikkObject `json:"updatedBy"`
	Created   string      `json:"created"`
	Updated   string      `json:"updated"`
}

func (Projects) path() string {
	return "v1/Core/Projects"
}

func (Projects) validFilter(options *ListOptions) bool {
	if len(options.UserIDs) > 0 ||
		options.FromDate != nil ||
		options.ToDate != nil {
		return false
	}
	return true
}

type GetItem interface {
	path(query string) string
}

type User struct {
	ObjectName           string             `json:"objectName"`
	ID                   int                `json:"id"`
	FirstName            string             `json:"firstName"`
	LastName             string             `json:"lastName"`
	License              string             `json:"license"`
	Email                string             `json:"email"`
	SocialSecurityNumber string             `json:"socialSecurityNumber"`
	MobilePhoneNumber    string             `json:"mobilePhoneNumber"`
	WorkPhoneNumber      string             `json:"workPhoneNumber"`
	Note                 string             `json:"note"`
	StartDate            dateutils.DateOnly `json:"startDate"`
	EndDate              dateutils.DateOnly `json:"endDate"`
	Department           blikkObject        `json:"department"`
	CostCenter           blikkObject        `json:"costCenter"`
	SalaryType           string             `json:"salaryType"`
	CostPerHour          float64            `json:"costPerHour"`
	EmployeeNumber       string             `json:"employeeNumber"`
	SettlementAccount    string             `json:"settlementAccount"`
	Address              struct {
		ObjectName        string `json:"objectName"`
		StreetAddress     string `json:"streetAddress"`
		AdditionalAddress string `json:"additionalAddress"`
		PostalCode        string `json:"postalCode"`
		City              string `json:"city"`
		State             string `json:"state"`
		CountryID         int    `json:"countryId"`
		CountryName       string `json:"countryName"`
	} `json:"address"`
	NextOfKin                 string      `json:"nextOfKin"`
	NextOfKinRelation         string      `json:"nextOfKinRelation"`
	NextOfKinPhoneNumber      string      `json:"nextOfKinPhoneNumber"`
	Manager                   blikkObject `json:"manager"`
	Schedule                  blikkObject `json:"schedule"`
	StandardTimeArticle       blikkObject `json:"standardTimeArticle"`
	StandardActivity          blikkObject `json:"standardActivity"`
	TimeReportingProfile      blikkObject `json:"timeReportingProfile"`
	TimeBankEnabled           bool        `json:"timeBankEnabled"`
	CurrentTimeBank           float64     `json:"currentTimeBank"`
	PlanningCapacityInPercent float64     `json:"planningCapacityInPercent"`
	Tags                      []string    `json:"tags"`
	Permissions               []string    `json:"permissions"`
	CreatedBy                 blikkObject `json:"createdBy"`
	UpdatedBy                 blikkObject `json:"updatedBy"`
	CreatedDate               string      `json:"createdDate"`
	UpdatedDate               string      `json:"updatedDate"`
}

func (User) path(query string) string {
	return "v1/Admin/Users/" + query
}
