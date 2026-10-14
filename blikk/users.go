package blikk

import (
	"context"
	"encoding/json"
	"iter"
	"strconv"
)

const usersPath = "v1/Admin/Users"

type UsersService struct {
	client *Client
}

const UserSortByName = "name"

type UserListParams struct {
	ListOptions
	Query             string `query:"filter.query"`
	IncludeDeleted    bool   `query:"filter.includeDeleted"`
	IncludeRestricted bool   `query:"filter.includeRestricted"`
}

type License string

const (
	LicenseTotal    License = "Total"
	LicenseStandard License = "Standard"
	LicenseBasic    License = "Basic"
	LicenseLight    License = "Light"
)

func (l *License) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, l, LicenseTotal, LicenseStandard, LicenseBasic, LicenseLight)
}

type SalaryType string

const (
	SalaryTypeNone    SalaryType = "None"
	SalaryTypeMonthly SalaryType = "Monthly"
	SalaryTypeHourly  SalaryType = "Hourly"
)

// The user list sends salaryType as a number while the user detail sends its name.
var salaryTypesByNumber = map[int]SalaryType{
	1: SalaryTypeNone,
	2: SalaryTypeMonthly,
	3: SalaryTypeHourly,
}

func (s *SalaryType) UnmarshalJSON(data []byte) error {
	var number int
	if err := json.Unmarshal(data, &number); err == nil {
		*s = salaryTypeFromNumber(number)
		return nil
	}
	return unmarshalEnum(data, s, SalaryTypeNone, SalaryTypeMonthly, SalaryTypeHourly)
}

func salaryTypeFromNumber(number int) SalaryType {
	if number == 0 {
		return ""
	}
	if salaryType, known := salaryTypesByNumber[number]; known {
		return salaryType
	}
	return SalaryType(strconv.Itoa(number))
}

type UserListItem struct {
	ID                   int            `json:"id"`
	FirstName            string         `json:"firstName"`
	LastName             string         `json:"lastName"`
	Email                string         `json:"email"`
	StartDate            *Date          `json:"startDate"`
	EndDate              *Date          `json:"endDate"`
	EmployeeNumber       string         `json:"employeeNumber"`
	TimeReportingProfile *NameID        `json:"timeReportingProfile"`
	Department           *NameID        `json:"department"`
	License              License        `json:"license"`
	Permissions          []string       `json:"permissions"`
	CreatedBy            *NameID        `json:"createdBy"`
	UpdatedBy            *NameID        `json:"updatedBy"`
	CreatedDate          *DateTime      `json:"createdDate"`
	UpdatedDate          *DateTime      `json:"updatedDate"`
	SalaryType           SalaryType     `json:"salaryType"`
	CostCenter           *CostCenterRef `json:"costCenter"`
	IsRestricted         bool           `json:"isRestricted"`
}

type User struct {
	ID                        int            `json:"id"`
	FirstName                 string         `json:"firstName"`
	LastName                  string         `json:"lastName"`
	License                   License        `json:"license"`
	Email                     string         `json:"email"`
	SocialSecurityNumber      string         `json:"socialSecurityNumber"`
	MobilePhoneNumber         string         `json:"mobilePhoneNumber"`
	WorkPhoneNumber           string         `json:"workPhoneNumber"`
	Note                      string         `json:"note"`
	StartDate                 *Date          `json:"startDate"`
	EndDate                   *Date          `json:"endDate"`
	Department                *NameID        `json:"department"`
	CostCenter                *CostCenterRef `json:"costCenter"`
	SalaryType                SalaryType     `json:"salaryType"`
	CostPerHour               float64        `json:"costPerHour"`
	EmployeeNumber            string         `json:"employeeNumber"`
	SettlementAccount         string         `json:"settlementAccount"`
	Address                   *Address       `json:"address"`
	NextOfKin                 string         `json:"nextOfKin"`
	NextOfKinRelation         string         `json:"nextOfKinRelation"`
	NextOfKinPhoneNumber      string         `json:"nextOfKinPhoneNumber"`
	Manager                   *NameID        `json:"manager"`
	Schedule                  *NameID        `json:"schedule"`
	StandardTimeArticle       *NameID        `json:"standardTimeArticle"`
	StandardActivity          *NameID        `json:"standardActivity"`
	TimeReportingProfile      *NameID        `json:"timeReportingProfile"`
	TimeBankEnabled           bool           `json:"timeBankEnabled"`
	CurrentTimeBank           float64        `json:"currentTimeBank"`
	PlanningCapacityInPercent float64        `json:"planningCapacityInPercent"`
	Tags                      []TagRef       `json:"tags"`
	Permissions               []string       `json:"permissions"`
	CreatedBy                 *NameID        `json:"createdBy"`
	UpdatedBy                 *NameID        `json:"updatedBy"`
	CreatedDate               *DateTime      `json:"createdDate"`
	UpdatedDate               *DateTime      `json:"updatedDate"`
	IsRestricted              bool           `json:"isRestricted"`
}

func (s *UsersService) List(ctx context.Context, params UserListParams) (*Page[UserListItem], error) {
	return listPage[UserListItem](ctx, s.client, usersPath, params)
}

func (s *UsersService) ListAll(ctx context.Context, params UserListParams) iter.Seq2[UserListItem, error] {
	return listAll[UserListItem](ctx, s.client, usersPath, params)
}

func (s *UsersService) Get(ctx context.Context, id int) (*User, error) {
	return getResource[User](ctx, s.client, usersPath, strconv.Itoa(id))
}
