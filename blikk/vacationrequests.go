package blikk

import (
	"context"
	"iter"
)

const vacationRequestsPath = "v1/Core/VacationRequests"

type VacationRequestsService struct {
	client *Client
}

type VacationStatus string

const (
	VacationStatusAttested       VacationStatus = "Attested"
	VacationStatusAwaitingAttest VacationStatus = "AwaitingAttest"
	VacationStatusDenied         VacationStatus = "Denied"
)

func (s *VacationStatus) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, s, VacationStatusAttested, VacationStatusAwaitingAttest, VacationStatusDenied)
}

const (
	VacationRequestSortByFromDate     = "fromDate"
	VacationRequestSortByToDate       = "toDate"
	VacationRequestSortByCreatedDate  = "createdDate"
	VacationRequestSortByUpdatedDate  = "updatedDate"
	VacationRequestSortByAttestedDate = "attestedDate"
	VacationRequestSortByDeniedDate   = "deniedDate"
)

type VacationRequestListParams struct {
	ListOptions
	CreatedDateFrom DateTime       `query:"filter.createdDateFrom"`
	CreatedDateTo   DateTime       `query:"filter.createdDateTo"`
	UpdatedDateFrom DateTime       `query:"filter.updatedDateFrom"`
	UpdatedDateTo   DateTime       `query:"filter.updatedDateTo"`
	FromDate        DateTime       `query:"filter.fromDate"`
	ToDate          DateTime       `query:"filter.toDate"`
	UserIDs         []int          `query:"filter.userIds"`
	DepartmentIDs   []int          `query:"filter.departmentIds"`
	ManagerIDs      []int          `query:"filter.managerIds"`
	AbsenceIDs      []int          `query:"filter.absenceIds"`
	VacationStatus  VacationStatus `query:"filter.vacationStatus"`
}

type VacationRequest struct {
	ID             int       `json:"id"`
	User           NameID    `json:"user"`
	FromDate       Date      `json:"fromDate"`
	ToDate         *Date     `json:"toDate"`
	Hours          float64   `json:"hours"`
	Description    string    `json:"description"`
	AbsenceProject NameID    `json:"absenceProject"`
	AttestedDate   *DateTime `json:"attestedDate"`
	AttestedByUser *NameID   `json:"attestedByUser"`
	DeniedDate     *DateTime `json:"deniedDate"`
	DeniedByUser   *NameID   `json:"deniedByUser"`
	UpdatedDate    *DateTime `json:"updatedDate"`
	UpdatedByUser  *NameID   `json:"updatedByUser"`
	CreatedDate    DateTime  `json:"createdDate"`
}

func (s *VacationRequestsService) List(ctx context.Context, params VacationRequestListParams) (*Page[VacationRequest], error) {
	return listPage[VacationRequest](ctx, s.client, vacationRequestsPath, params)
}

func (s *VacationRequestsService) ListAll(ctx context.Context, params VacationRequestListParams) iter.Seq2[VacationRequest, error] {
	return listAll[VacationRequest](ctx, s.client, vacationRequestsPath, params)
}
