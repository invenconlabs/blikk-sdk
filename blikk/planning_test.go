package blikk

import (
	"context"
	"iter"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanningProjectsWithPlanningEncodesParamsAndDecodesFixture(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "planning_projects.json"))

	page, err := client.Planning.ProjectsWithPlanning(context.Background(), PlanningProjectsParams{
		ListOptions:     ListOptions{PageSize: 50},
		FromDate:        NewDate(2026, time.January, 1),
		ToDate:          NewDate(2026, time.December, 31),
		ExcludeProjects: true,
		ExcludeAbsence:  true,
		ExcludeInternal: true,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Core/Planning/HasResourcePlanning/Projects", request.Path)
	assert.Equal(t, map[string][]string{
		"page":            {"1"},
		"pageSize":        {"50"},
		"fromDate":        {"2026-01-01"},
		"toDate":          {"2026-12-31"},
		"excludeProjects": {"true"},
		"excludeAbsence":  {"true"},
		"excludeInternal": {"true"},
	}, request.Query)
	require.Len(t, page.Items, 3)
	assert.Equal(t, PlanningProject{ID: 101, OrderNumber: "P001", Title: "Website Redesign"}, page.Items[0])
}

func TestPlanningUsersWithPlanningSendsExplicitFalse(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "planning_users.json"))

	page, err := client.Planning.UsersWithPlanning(context.Background(), PlanningUsersParams{
		FromDate:          NewDate(2026, time.January, 1),
		ToDate:            NewDate(2026, time.January, 7),
		ExcludeDeleted:    Ptr(false),
		ExcludeRestricted: Ptr(true),
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Core/Planning/HasResourcePlanning/Users", request.Path)
	assert.Equal(t, map[string][]string{
		"page":              {"1"},
		"fromDate":          {"2026-01-01"},
		"toDate":            {"2026-01-07"},
		"excludeDeleted":    {"false"},
		"excludeRestricted": {"true"},
	}, request.Query)
	require.Len(t, page.Items, 3)
	assert.Equal(t, PlanningUser{ID: 3, EmploymentNumber: "1003", FirstName: "Thorbjörn", LastName: "Fälldin"}, page.Items[2])
}

func TestPlanningProjectSummariesForUserEncodesParamsAndDecodesFixture(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "planning_project_summaries.json"))

	page, err := client.Planning.ProjectSummariesForUser(context.Background(), PlanningProjectSummaryParams{
		UserID:          1,
		FromDate:        NewDate(2026, time.January, 1),
		ToDate:          NewDate(2026, time.December, 31),
		ExcludeProjects: true,
		ExcludeAbsence:  true,
		ExcludeInternal: true,
		GroupBy:         PlanningGroupByWeek,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Core/Planning/GetPlanningSummaries/Projects", request.Path)
	assert.Equal(t, map[string][]string{
		"page":            {"1"},
		"userId":          {"1"},
		"fromDate":        {"2026-01-01"},
		"toDate":          {"2026-12-31"},
		"excludeProjects": {"true"},
		"excludeAbsence":  {"true"},
		"excludeInternal": {"true"},
		"groupBy":         {"week"},
	}, request.Query)
	require.Len(t, page.Items, 2)
	summary := page.Items[0]
	assert.Equal(t, 101, summary.ProjectID)
	assert.Equal(t, "2026-01-01", summary.FromDate.String())
	assert.Equal(t, "2026-12-31", summary.ToDate.String())
	assert.Equal(t, "Website Redesign", summary.ProjectName)
	assert.Equal(t, 160.0, summary.PlannedHours)
	assert.Equal(t, 152000.0, summary.ExpectedInvoiceableAmount)
	assert.Equal(t, Ptr(80000.0), summary.ExpectedCost)
	assert.Nil(t, page.Items[1].ExpectedCost)
}

func TestPlanningUserSummariesForProjectEncodesParamsAndDecodesFixture(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "planning_user_summaries.json"))

	page, err := client.Planning.UserSummariesForProject(context.Background(), PlanningUserSummaryParams{
		ProjectID: 42,
		FromDate:  NewDate(2026, time.January, 1),
		ToDate:    NewDate(2026, time.December, 31),
		GroupBy:   PlanningGroupByMonth,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Core/Planning/GetPlanningSummaries/Users", request.Path)
	assert.Equal(t, map[string][]string{
		"page":      {"1"},
		"projectId": {"42"},
		"fromDate":  {"2026-01-01"},
		"toDate":    {"2026-12-31"},
		"groupBy":   {"month"},
	}, request.Query)
	require.Len(t, page.Items, 2)
	summary := page.Items[0]
	assert.Equal(t, 1, summary.UserID)
	assert.Equal(t, "2026-01-01", summary.FromDate.String())
	assert.Equal(t, "2026-12-31", summary.ToDate.String())
	assert.Equal(t, "Ingvar Carlsson", summary.UserName)
	assert.Equal(t, 100.1, summary.PlannedHours)
	assert.Equal(t, 90090.0, summary.ExpectedInvoiceableAmount)
	assert.Equal(t, Ptr(45045.0), summary.ExpectedCost)
	assert.Nil(t, page.Items[1].ExpectedCost)
}

func TestPlanningListAllVariantsWalkPages(t *testing.T) {
	ctx := context.Background()
	client := newTestClient(t, routes(t, map[string]http.HandlerFunc{
		"GET /v1/Core/Planning/HasResourcePlanning/Projects":  serveFixture(t, "planning_projects.json"),
		"GET /v1/Core/Planning/HasResourcePlanning/Users":     serveFixture(t, "planning_users.json"),
		"GET /v1/Core/Planning/GetPlanningSummaries/Projects": serveFixture(t, "planning_project_summaries.json"),
		"GET /v1/Core/Planning/GetPlanningSummaries/Users":    serveFixture(t, "planning_user_summaries.json"),
	}))

	assert.Len(t, collectPlanningItems(t, client.Planning.ListAllProjectsWithPlanning(ctx, PlanningProjectsParams{})), 3)
	assert.Len(t, collectPlanningItems(t, client.Planning.ListAllUsersWithPlanning(ctx, PlanningUsersParams{})), 3)
	assert.Len(t, collectPlanningItems(t, client.Planning.ListAllProjectSummariesForUser(ctx, PlanningProjectSummaryParams{})), 2)
	assert.Len(t, collectPlanningItems(t, client.Planning.ListAllUserSummariesForProject(ctx, PlanningUserSummaryParams{})), 2)
}

func collectPlanningItems[T any](t *testing.T, seq iter.Seq2[T, error]) []T {
	t.Helper()
	var items []T
	for item, err := range seq {
		require.NoError(t, err)
		items = append(items, item)
	}
	return items
}
