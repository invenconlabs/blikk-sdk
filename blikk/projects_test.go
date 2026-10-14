package blikk

import (
	"context"
	"net/http"
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

func TestProjectsCreateSendsFullBodyAndDecodesResult(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusCreated, `{"id":1,"nextProjectNumber":"P-124"}`))

	result, err := client.Projects.Create(context.Background(), ProjectCreateInput{
		ProjectNumber:            "P-123",
		Title:                    "Lorem projectum",
		Description:              "A nice description",
		ContactID:                1,
		ResponsibleID:            2,
		StatusID:                 3,
		ProjectCollectionID:      4,
		StartDate:                DateTimeOf(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:                  DateTimeOf(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)),
		WorkSiteAddress:          &ProjectWorkSiteAddressInput{StreetAddress: "Street 1", PostalCode: "12345", City: "Town", Longitude: 19.04, Latitude: 47.49},
		CategoryID:               5,
		CostCenterID:             6,
		InvoiceText:              "Text here",
		YourReference:            "Your ref",
		OurReference:             "Our ref",
		CustomerReferenceMarking: "Marking",
		AllIsMember:              true,
	})

	require.NoError(t, err)
	assert.Equal(t, &ProjectCreateResult{ID: 1, NextProjectNumber: "P-124"}, result)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/Core/Projects", request.Path)
	assert.JSONEq(t, `{
		"projectNumber": "P-123",
		"title": "Lorem projectum",
		"description": "A nice description",
		"contactId": 1,
		"responsibleId": 2,
		"statusId": 3,
		"projectCollectionId": 4,
		"startDate": "2023-01-01 00:00:00",
		"endDate": "2023-12-31 00:00:00",
		"workSiteAddress": {"streetAddress": "Street 1", "postalCode": "12345", "city": "Town", "longitude": 19.04, "latitude": 47.49},
		"categoryId": 5,
		"costCenterId": 6,
		"invoiceText": "Text here",
		"yourReference": "Your ref",
		"ourReference": "Our ref",
		"customerReferenceMarking": "Marking",
		"allIsMember": true
	}`, string(request.Body))
}

func TestProjectsCreateMinimalInputSendsRequiredFieldsAndBooleans(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusCreated, `{"id":1,"nextProjectNumber":"P-124"}`))

	_, err := client.Projects.Create(context.Background(), ProjectCreateInput{Title: "Lorem projectum", ContactID: 1, ResponsibleID: 2})

	require.NoError(t, err)
	assert.JSONEq(t, `{"title":"Lorem projectum","contactId":1,"responsibleId":2,"allIsMember":false}`, string(handler.last().Body))
}

func TestProjectsUpdateSendsFullBodyWithID(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusNoContent, ""))

	err := client.Projects.Update(context.Background(), 42, ProjectUpdateInput{
		Title:                    "Our Project",
		ProjectCollectionID:      19,
		StartDate:                DateTimeOf(time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC)),
		EndDate:                  DateTimeOf(time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC)),
		TagIDs:                   []int{12, 14, 16},
		ContactID:                53,
		WorkSiteAddress:          &ProjectWorkSiteAddressInput{StreetAddress: "Váci utca 1-3.", PostalCode: "1051", City: "Budapest", Longitude: 19.040235, Latitude: 47.497913},
		ResponsibleID:            32,
		SalesResponsibleID:       43,
		Description:              "Our Common Project",
		StatusID:                 12,
		ReadyForInvoice:          true,
		CostCenterID:             21,
		BookingProjectID:         42,
		CategoryID:               1,
		InvoiceText:              "Invoice #17",
		YourReference:            "Mr. White",
		OurReference:             "Mr. Black",
		CustomerReferenceMarking: "Our Project 2026",
		AllIsMember:              true,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodPut, request.Method)
	assert.Equal(t, "/v1/Core/Projects/42", request.Path)
	assert.JSONEq(t, `{
		"id": 42,
		"title": "Our Project",
		"projectCollectionId": 19,
		"startDate": "2026-02-13 00:00:00",
		"endDate": "2026-04-13 00:00:00",
		"tagIds": [12, 14, 16],
		"contactId": 53,
		"workSiteAddress": {"streetAddress": "Váci utca 1-3.", "postalCode": "1051", "city": "Budapest", "longitude": 19.040235, "latitude": 47.497913},
		"responsibleId": 32,
		"salesResponsibleId": 43,
		"description": "Our Common Project",
		"statusId": 12,
		"readyForInvoice": true,
		"costCenterId": 21,
		"bookingProjectId": 42,
		"categoryId": 1,
		"invoiceText": "Invoice #17",
		"yourReference": "Mr. White",
		"ourReference": "Mr. Black",
		"customerReferenceMarking": "Our Project 2026",
		"allIsMember": true
	}`, string(request.Body))
}

func TestProjectsUpdateMinimalInputSendsIDFromArgumentAndBooleans(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusNoContent, ""))

	err := client.Projects.Update(context.Background(), 42, ProjectUpdateInput{ID: 7})

	require.NoError(t, err)
	assert.JSONEq(t, `{"id":42,"readyForInvoice":false,"allIsMember":false}`, string(handler.last().Body))
}

func TestProjectsNextProjectNumberUnwrapsObject(t *testing.T) {
	client, _ := newRecordingClient(t, routes(t, map[string]http.HandlerFunc{
		"GET /v1/Core/Projects/nextProjectNumber": respondWith(http.StatusOK, `{"nextProjectNumber":"001"}`),
	}))

	next, err := client.Projects.NextProjectNumber(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "001", next)
}

func TestProjectsListCommentsEncodesFiltersAndDecodesFixture(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Core/Projects/1/Comments", "projects_comments_list.json"))
	from := DateTimeOf(time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC))
	to := DateTimeOf(time.Date(2026, 2, 1, 8, 0, 0, 0, time.UTC))

	page, err := client.Projects.ListComments(context.Background(), 1, ProjectCommentListParams{
		ListOptions:      ListOptions{SortBy: ProjectCommentSortByCreatedDate, SortOrder: SortDescending},
		CreatedFrom:      from,
		CreatedTo:        to,
		CreatedByUserIDs: []int{1, 2},
	})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"page":                    {"1"},
		"sortBy":                  {"createdDate"},
		"sortOrder":               {"descending"},
		"filter.createdFrom":      {"2026-01-01 08:00:00"},
		"filter.createdTo":        {"2026-02-01 08:00:00"},
		"filter.createdByUserIds": {"1,2"},
	}, handler.last().Query)
	require.Len(t, page.Items, 1)
	comment := page.Items[0]
	assert.Equal(t, 1, comment.ID)
	assert.Equal(t, "Lorem ipsum dolar sit amet", comment.Text)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, comment.CreatedByUser)
	assert.Nil(t, comment.UpdatedByUser)
	assert.Equal(t, "2026-10-11 09:40:57", comment.CreatedDate.String())
	assert.Nil(t, comment.UpdatedDate)
}

func TestProjectsListAllCommentsUsesProjectPath(t *testing.T) {
	client, _ := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Core/Projects/1/Comments", "projects_comments_list.json"))

	var ids []int
	for comment, err := range client.Projects.ListAllComments(context.Background(), 1, ProjectCommentListParams{}) {
		require.NoError(t, err)
		ids = append(ids, comment.ID)
	}

	assert.Equal(t, []int{1}, ids)
}

func TestProjectsCreateCommentSendsTextAndReturnsID(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"id":9}`))

	id, err := client.Projects.CreateComment(context.Background(), 1, "Lorem ipsum dolar sit amet.")

	require.NoError(t, err)
	assert.Equal(t, 9, id)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/Core/Projects/1/Comments", request.Path)
	assert.JSONEq(t, `{"text":"Lorem ipsum dolar sit amet."}`, string(request.Body))
}

func TestProjectsUpdateCommentSendsPut(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	err := client.Projects.UpdateComment(context.Background(), 1, 9, "Edited")

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodPut, request.Method)
	assert.Equal(t, "/v1/Core/Projects/1/Comments/9", request.Path)
	assert.JSONEq(t, `{"text":"Edited"}`, string(request.Body))
}

func TestProjectsDeleteCommentSendsDelete(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	err := client.Projects.DeleteComment(context.Background(), 1, 9)

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodDelete, request.Method)
	assert.Equal(t, "/v1/Core/Projects/1/Comments/9", request.Path)
	assert.Empty(t, request.Body)
}

func TestProjectsTimeCalculationDecodesFixture(t *testing.T) {
	client, _ := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Core/Projects/1/TimeCalculation", "projects_timecalculation.json"))

	calculation, err := client.Projects.TimeCalculation(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, &ProjectTimeCalculation{
		Total:        50,
		IsByActivity: true,
		ActivityCalculations: []ProjectActivityCalculation{
			{Activity: &NameID{ID: 1, Name: "Activity 1"}, Calculation: 8},
			{Activity: &NameID{ID: 2, Name: "Activity 2"}, Calculation: 26},
			{Activity: &NameID{ID: 3, Name: "Activity 3"}, Calculation: 16},
		},
	}, calculation)
}

func TestProjectsAddMemberSendsPostWithoutBody(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusCreated, ""))

	err := client.Projects.AddMember(context.Background(), 42, 15)

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/Core/Projects/42/members/15", request.Path)
	assert.Empty(t, request.Body)
}

func TestProjectsAddMemberConflictIsAnError(t *testing.T) {
	client, _ := newRecordingClient(t, respondWith(http.StatusConflict, ""))

	err := client.Projects.AddMember(context.Background(), 42, 15)

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusConflict, apiErr.StatusCode)
}

func TestProjectsRemoveMemberSendsDelete(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusNoContent, ""))

	err := client.Projects.RemoveMember(context.Background(), 42, 15)

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodDelete, request.Method)
	assert.Equal(t, "/v1/Core/Projects/42/members/15", request.Path)
}
