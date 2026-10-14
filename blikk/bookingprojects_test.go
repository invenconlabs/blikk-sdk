package blikk

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookingProjectsListEncodesEveryFilter(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "bookingprojects_list.json"))
	moment := DateTimeOf(time.Date(2026, 10, 1, 8, 0, 0, 0, time.UTC))

	_, err := client.BookingProjects.List(context.Background(), BookingProjectListParams{
		ListOptions:     ListOptions{SortBy: BookingProjectSortByBookingNumber, SortOrder: SortDescending},
		CreatedDateFrom: moment,
		CreatedDateTo:   moment,
		UpdatedDateFrom: moment,
		UpdatedDateTo:   moment,
		StartDateFrom:   moment,
		StartDateTo:     moment,
		EndDateFrom:     moment,
		EndDateTo:       moment,
		IsActive:        Ptr(false),
		Query:           "xyz",
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/admin/bookingprojects", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                   {"1"},
		"sortBy":                 {"bookingNumber"},
		"sortOrder":              {"descending"},
		"filter.createdDateFrom": {"2026-10-01 08:00:00"},
		"filter.createdDateTo":   {"2026-10-01 08:00:00"},
		"filter.updatedDateFrom": {"2026-10-01 08:00:00"},
		"filter.updatedDateTo":   {"2026-10-01 08:00:00"},
		"filter.startDateFrom":   {"2026-10-01 08:00:00"},
		"filter.startDateTo":     {"2026-10-01 08:00:00"},
		"filter.endDateFrom":     {"2026-10-01 08:00:00"},
		"filter.endDateTo":       {"2026-10-01 08:00:00"},
		"filter.isActive":        {"false"},
		"filter.query":           {"xyz"},
	}, request.Query)
}

func TestBookingProjectsListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, "GET", "/v1/admin/bookingprojects", "bookingprojects_list.json"))

	page, err := client.BookingProjects.List(context.Background(), BookingProjectListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	project := page.Items[0]
	assert.Equal(t, 42, project.ID)
	assert.Equal(t, "Booking Project 1", project.Name)
	assert.Equal(t, "XYZ", project.BookingNumber)
	assert.True(t, project.IsActive)
	assert.Nil(t, project.StartDate)
	assert.Nil(t, project.EndDate)
	assert.Equal(t, NameID{ID: 1, Name: "John Doe"}, project.CreatedBy)
	assert.Equal(t, &NameID{ID: 2, Name: "Jane Doe"}, project.UpdatedBy)
	assert.Equal(t, "2026-09-28 09:40:57", project.Created.String())
	assert.Equal(t, "2026-10-02 09:40:57", project.Updated.String())
}

func TestBookingProjectsListAllWalksPages(t *testing.T) {
	client := newTestClient(t, serveFixture(t, "bookingprojects_list.json"))

	var ids []int
	for project, err := range client.BookingProjects.ListAll(context.Background(), BookingProjectListParams{}) {
		require.NoError(t, err)
		ids = append(ids, project.ID)
	}

	assert.Equal(t, []int{42}, ids)
}

func TestBookingProjectsGetDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, "GET", "/v1/admin/bookingprojects/42", "bookingprojects_detail.json"))

	project, err := client.BookingProjects.Get(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, 42, project.ID)
	assert.Equal(t, "123", project.BookingNumber)
	assert.Equal(t, "2026-10-13 09:40:57", project.StartDate.String())
	assert.Equal(t, "2026-10-18 09:40:57", project.EndDate.String())
	assert.Nil(t, project.UpdatedBy)
	assert.Nil(t, project.Updated)
}

func TestBookingProjectsGetByBookingNumberSendsQuery(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "bookingprojects_detail.json"))

	project, err := client.BookingProjects.GetByBookingNumber(context.Background(), "12 3")

	require.NoError(t, err)
	assert.Equal(t, 42, project.ID)
	request := handler.last()
	assert.Equal(t, http.MethodGet, request.Method)
	assert.Equal(t, "/v1/admin/bookingprojects", request.Path)
	assert.Equal(t, map[string][]string{"bookingnumber": {"12 3"}}, request.Query)
}

func TestBookingProjectsGetByEmptyBookingNumberIsAnError(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "bookingprojects_detail.json"))

	_, err := client.BookingProjects.GetByBookingNumber(context.Background(), "")

	assert.ErrorContains(t, err, "booking number")
	assert.Empty(t, handler.all())
}

func TestBookingProjectsCreateSendsFullInput(t *testing.T) {
	client, handler := newRecordingClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(loadFixture(t, "bookingprojects_detail.json"))
	})
	start := DateTimeOf(time.Date(2026, 10, 13, 9, 40, 57, 0, time.UTC))
	end := DateTimeOf(time.Date(2026, 10, 18, 9, 40, 57, 0, time.UTC))

	created, err := client.BookingProjects.Create(context.Background(), BookingProjectInput{
		Name:          "Booking Project 1",
		BookingNumber: "123",
		IsActive:      Ptr(false),
		StartDate:     start,
		EndDate:       end,
	})

	require.NoError(t, err)
	assert.Equal(t, 42, created.ID)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/admin/bookingprojects", request.Path)
	assert.JSONEq(t, `{
		"name": "Booking Project 1",
		"bookingNumber": "123",
		"isActive": false,
		"startDate": "2026-10-13 09:40:57",
		"endDate": "2026-10-18 09:40:57"
	}`, string(request.Body))
}

func TestBookingProjectsUpdateSendsMinimalInputAndReturnsProject(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "bookingprojects_detail.json"))

	updated, err := client.BookingProjects.Update(context.Background(), 42, BookingProjectInput{Name: "Booking Project 2", BookingNumber: "321"})

	require.NoError(t, err)
	assert.Equal(t, 42, updated.ID)
	request := handler.last()
	assert.Equal(t, http.MethodPut, request.Method)
	assert.Equal(t, "/v1/admin/bookingprojects/42", request.Path)
	assert.JSONEq(t, `{"name": "Booking Project 2", "bookingNumber": "321"}`, string(request.Body))
}
