package blikk

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectsListEncodesFiltersAndDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "projects_list.json")}
	client := newTestClient(t, handler.ServeHTTP)
	since := DateTimeOf(time.Date(2021, 1, 1, 14, 0, 0, 0, time.UTC))

	page, err := client.Projects.List(context.Background(), ProjectListParams{
		ListOptions:           ListOptions{Page: 2, PageSize: 100},
		Query:                 "roof",
		CreatedDateFrom:       since,
		CreatedDateTo:         since,
		UpdatedDateFrom:       since,
		UpdatedDateTo:         since,
		CreatedOrUpdatedSince: since,
		ProjectTagIDs:         []int{1, 2},
		ProjectCategoryIDs:    []int{3},
		ProjectStatusIDs:      []int{4, 5},
		ProjectType:           InvoiceTypeFixedPrice,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Core/Projects", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                         {"2"},
		"pageSize":                     {"100"},
		"filter.query":                 {"roof"},
		"filter.createdDateFrom":       {"2021-01-01 14:00:00"},
		"filter.createdDateTo":         {"2021-01-01 14:00:00"},
		"filter.updatedDateFrom":       {"2021-01-01 14:00:00"},
		"filter.updatedDateTo":         {"2021-01-01 14:00:00"},
		"filter.createdOrUpdatedSince": {"2021-01-01 14:00:00"},
		"filter.projectTagIds":         {"1,2"},
		"filter.projectCategoryIds":    {"3"},
		"filter.projectStatusIds":      {"4,5"},
		"filter.projectType":           {"FixedPrice"},
	}, request.Query)

	require.Len(t, page.Items, 1)
	project := page.Items[0]
	assert.Equal(t, 1, project.ID)
	assert.Equal(t, "26", project.OrderNumber)
	assert.Equal(t, "Lorem ipsum", project.Title)
	assert.Equal(t, &ProjectStatus{ID: 1, Name: "Ongoing", IsCompletedStatus: false}, project.Status)
	assert.Equal(t, &ProjectCategory{ID: 1, Name: "Lorem ipsum", Color: "#e3168a"}, project.Category)
	assert.Nil(t, project.SalesResponsible)
	assert.Equal(t, "2021-02-01", project.StartDate.String())
	assert.Equal(t, "2022-02-15", project.EndDate.String())
	assert.Equal(t, InvoiceTypeOngoing, project.InvoiceType)
	assert.Equal(t, &Location{Longitude: 21.48181104660036, Latitude: 65.3175658912035, StreetAddress: "Some street 123", PostalCode: "94100", City: "Washington", CountryName: "USA"}, project.Location)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, project.ProjectManager)
	assert.Equal(t, &NameID{ID: 1, Name: "Ipsum"}, project.Customer)
	assert.Equal(t, &ProjectRef{ID: 14, Name: "Dolar sit amet", Number: "P012"}, project.ProjectCollection)
	assert.Equal(t, []Tag{{ID: 1, Title: "Ipsum", Color: "#41b18b"}}, project.Tags)
	assert.Equal(t, &CostCenterRef{ID: 123, Name: "Cost Center 123", Code: "CC123"}, project.CostCenter)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, project.CreatedBy)
	assert.Nil(t, project.UpdatedBy)
	assert.Equal(t, "2026-09-28 09:40:57", project.Created.String())
	assert.Nil(t, project.Updated)
}

func TestProjectsGetDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "project.json")}
	client := newTestClient(t, handler.ServeHTTP)

	project, err := client.Projects.Get(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, "/v1/Core/Projects/42", handler.last().Path)
	assert.Equal(t, 42, project.ID)
	assert.True(t, project.ReadyForInvoice)
	assert.Equal(t, "Lorem ipsum", project.InvoiceText)
	assert.Equal(t, "Lorem ipsum", project.YourReference)
	assert.Equal(t, "Lorem ipsum", project.OurReference)
	assert.Equal(t, "Lorem ipsum", project.CustomerReferenceMarking)
	assert.Equal(t, "Lorem ipsum", project.Description)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, project.SalesResponsible)
	assert.Equal(t, &NameID{ID: 609, Name: "Booking Project 570"}, project.BookingProject)
	assert.Equal(t, "2026-10-12 09:40:57", project.CompletedDate.String())
	assert.Equal(t, "Lorem ipsum", project.ElectronicInvoiceReferenceCode)
	assert.Equal(t, "Lorem ipsum", project.ImportantInvoiceInformation)
	assert.False(t, project.AllIsMember)
	assert.Equal(t, []NameID{{ID: 1, Name: "John Doe"}, {ID: 2, Name: "Jane Doe"}}, project.Members)
}

func TestProjectsGetByOrderNumberUsesOrderPath(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "project.json")}
	client := newTestClient(t, handler.ServeHTTP)

	project, err := client.Projects.GetByOrderNumber(context.Background(), "100-5")

	require.NoError(t, err)
	assert.Equal(t, "/v1/Core/Projects/order/100-5", handler.last().Path)
	assert.Equal(t, 42, project.ID)
}
