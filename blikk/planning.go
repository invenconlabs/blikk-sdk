package blikk

import (
	"context"
	"iter"
)

const (
	projectsWithPlanningPath     = "v1/Core/Planning/HasResourcePlanning/Projects"
	usersWithPlanningPath        = "v1/Core/Planning/HasResourcePlanning/Users"
	planningProjectSummariesPath = "v1/Core/Planning/GetPlanningSummaries/Projects"
	planningUserSummariesPath    = "v1/Core/Planning/GetPlanningSummaries/Users"
)

type PlanningService struct {
	client *Client
}

const (
	PlanningGroupByDay   = "day"
	PlanningGroupByWeek  = "week"
	PlanningGroupByMonth = "month"
)

type PlanningProjectsParams struct {
	ListOptions
	FromDate        Date `query:"fromDate"`
	ToDate          Date `query:"toDate"`
	ExcludeProjects bool `query:"excludeProjects"`
	ExcludeAbsence  bool `query:"excludeAbsence"`
	ExcludeInternal bool `query:"excludeInternal"`
}

// PlanningUsersParams uses *bool because the API excludes deleted and restricted users unless told false.
type PlanningUsersParams struct {
	ListOptions
	FromDate          Date  `query:"fromDate"`
	ToDate            Date  `query:"toDate"`
	ExcludeDeleted    *bool `query:"excludeDeleted"`
	ExcludeRestricted *bool `query:"excludeRestricted"`
}

type PlanningProjectSummaryParams struct {
	ListOptions
	UserID          int    `query:"userId"`
	FromDate        Date   `query:"fromDate"`
	ToDate          Date   `query:"toDate"`
	ExcludeProjects bool   `query:"excludeProjects"`
	ExcludeAbsence  bool   `query:"excludeAbsence"`
	ExcludeInternal bool   `query:"excludeInternal"`
	GroupBy         string `query:"groupBy"`
}

type PlanningUserSummaryParams struct {
	ListOptions
	ProjectID int    `query:"projectId"`
	FromDate  Date   `query:"fromDate"`
	ToDate    Date   `query:"toDate"`
	GroupBy   string `query:"groupBy"`
}

type PlanningProject struct {
	ID          int    `json:"id"`
	OrderNumber string `json:"orderNumber"`
	Title       string `json:"title"`
}

type PlanningUser struct {
	ID               int    `json:"id"`
	EmploymentNumber string `json:"employmentNumber"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
}

type PlanningProjectSummary struct {
	ProjectID                 int      `json:"projectId"`
	FromDate                  *Date    `json:"fromDate"`
	ToDate                    *Date    `json:"toDate"`
	ProjectName               string   `json:"projectName"`
	PlannedHours              float64  `json:"plannedHours"`
	ExpectedInvoiceableAmount float64  `json:"expectedInvoiceableAmount"`
	ExpectedCost              *float64 `json:"expectedCost"`
}

type PlanningUserSummary struct {
	UserID                    int      `json:"userId"`
	FromDate                  *Date    `json:"fromDate"`
	ToDate                    *Date    `json:"toDate"`
	UserName                  string   `json:"userName"`
	PlannedHours              float64  `json:"plannedHours"`
	ExpectedInvoiceableAmount float64  `json:"expectedInvoiceableAmount"`
	ExpectedCost              *float64 `json:"expectedCost"`
}

func (s *PlanningService) ProjectsWithPlanning(ctx context.Context, params PlanningProjectsParams) (*Page[PlanningProject], error) {
	return listPage[PlanningProject](ctx, s.client, projectsWithPlanningPath, params)
}

// ListAllProjectsWithPlanning walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *PlanningService) ListAllProjectsWithPlanning(ctx context.Context, params PlanningProjectsParams) iter.Seq2[PlanningProject, error] {
	return listAll[PlanningProject](ctx, s.client, projectsWithPlanningPath, params)
}

func (s *PlanningService) UsersWithPlanning(ctx context.Context, params PlanningUsersParams) (*Page[PlanningUser], error) {
	return listPage[PlanningUser](ctx, s.client, usersWithPlanningPath, params)
}

// ListAllUsersWithPlanning walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *PlanningService) ListAllUsersWithPlanning(ctx context.Context, params PlanningUsersParams) iter.Seq2[PlanningUser, error] {
	return listAll[PlanningUser](ctx, s.client, usersWithPlanningPath, params)
}

func (s *PlanningService) ProjectSummariesForUser(ctx context.Context, params PlanningProjectSummaryParams) (*Page[PlanningProjectSummary], error) {
	return listPage[PlanningProjectSummary](ctx, s.client, planningProjectSummariesPath, params)
}

// ListAllProjectSummariesForUser walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *PlanningService) ListAllProjectSummariesForUser(ctx context.Context, params PlanningProjectSummaryParams) iter.Seq2[PlanningProjectSummary, error] {
	return listAll[PlanningProjectSummary](ctx, s.client, planningProjectSummariesPath, params)
}

func (s *PlanningService) UserSummariesForProject(ctx context.Context, params PlanningUserSummaryParams) (*Page[PlanningUserSummary], error) {
	return listPage[PlanningUserSummary](ctx, s.client, planningUserSummariesPath, params)
}

// ListAllUserSummariesForProject walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *PlanningService) ListAllUserSummariesForProject(ctx context.Context, params PlanningUserSummaryParams) iter.Seq2[PlanningUserSummary, error] {
	return listAll[PlanningUserSummary](ctx, s.client, planningUserSummariesPath, params)
}
