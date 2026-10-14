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

type PlanningProjectListParams struct {
	ListOptions
	FromDate        Date `query:"fromDate"`
	ToDate          Date `query:"toDate"`
	ExcludeProjects bool `query:"excludeProjects"`
	ExcludeAbsence  bool `query:"excludeAbsence"`
	ExcludeInternal bool `query:"excludeInternal"`
}

// PlanningUserListParams uses *bool because the API excludes deleted and restricted users unless told false.
type PlanningUserListParams struct {
	ListOptions
	FromDate          Date  `query:"fromDate"`
	ToDate            Date  `query:"toDate"`
	ExcludeDeleted    *bool `query:"excludeDeleted"`
	ExcludeRestricted *bool `query:"excludeRestricted"`
}

type PlanningProjectSummaryListParams struct {
	ListOptions
	UserID          int    `query:"userId"`
	FromDate        Date   `query:"fromDate"`
	ToDate          Date   `query:"toDate"`
	ExcludeProjects bool   `query:"excludeProjects"`
	ExcludeAbsence  bool   `query:"excludeAbsence"`
	ExcludeInternal bool   `query:"excludeInternal"`
	GroupBy         string `query:"groupBy"`
}

type PlanningUserSummaryListParams struct {
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

func (s *PlanningService) ListProjectsWithPlanning(ctx context.Context, params PlanningProjectListParams) (*Page[PlanningProject], error) {
	return listPage[PlanningProject](ctx, s.client, projectsWithPlanningPath, params)
}

func (s *PlanningService) ListAllProjectsWithPlanning(ctx context.Context, params PlanningProjectListParams) iter.Seq2[PlanningProject, error] {
	return listAll[PlanningProject](ctx, s.client, projectsWithPlanningPath, params)
}

func (s *PlanningService) ListUsersWithPlanning(ctx context.Context, params PlanningUserListParams) (*Page[PlanningUser], error) {
	return listPage[PlanningUser](ctx, s.client, usersWithPlanningPath, params)
}

func (s *PlanningService) ListAllUsersWithPlanning(ctx context.Context, params PlanningUserListParams) iter.Seq2[PlanningUser, error] {
	return listAll[PlanningUser](ctx, s.client, usersWithPlanningPath, params)
}

func (s *PlanningService) ListProjectSummariesForUser(ctx context.Context, params PlanningProjectSummaryListParams) (*Page[PlanningProjectSummary], error) {
	return listPage[PlanningProjectSummary](ctx, s.client, planningProjectSummariesPath, params)
}

func (s *PlanningService) ListAllProjectSummariesForUser(ctx context.Context, params PlanningProjectSummaryListParams) iter.Seq2[PlanningProjectSummary, error] {
	return listAll[PlanningProjectSummary](ctx, s.client, planningProjectSummariesPath, params)
}

func (s *PlanningService) ListUserSummariesForProject(ctx context.Context, params PlanningUserSummaryListParams) (*Page[PlanningUserSummary], error) {
	return listPage[PlanningUserSummary](ctx, s.client, planningUserSummariesPath, params)
}

func (s *PlanningService) ListAllUserSummariesForProject(ctx context.Context, params PlanningUserSummaryListParams) iter.Seq2[PlanningUserSummary, error] {
	return listAll[PlanningUserSummary](ctx, s.client, planningUserSummariesPath, params)
}
