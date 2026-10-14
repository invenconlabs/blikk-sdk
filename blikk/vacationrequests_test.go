package blikk

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVacationRequestsListEncodesEveryFilter(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "vacationrequests_list.json"))
	moment := DateTimeOf(time.Date(2026, 10, 1, 8, 0, 0, 0, time.UTC))

	_, err := client.VacationRequests.List(context.Background(), VacationRequestListParams{
		ListOptions:     ListOptions{SortBy: VacationRequestSortByFromDate, SortOrder: SortAscending},
		CreatedDateFrom: moment,
		CreatedDateTo:   moment,
		UpdatedDateFrom: moment,
		UpdatedDateTo:   moment,
		FromDate:        moment,
		ToDate:          moment,
		UserIDs:         []int{1, 2},
		DepartmentIDs:   []int{3},
		ManagerIDs:      []int{4, 5},
		AbsenceIDs:      []int{6},
		VacationStatus:  VacationStatusAwaitingAttest,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Core/VacationRequests", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                   {"1"},
		"sortBy":                 {"fromDate"},
		"sortOrder":              {"ascending"},
		"filter.createdDateFrom": {"2026-10-01 08:00:00"},
		"filter.createdDateTo":   {"2026-10-01 08:00:00"},
		"filter.updatedDateFrom": {"2026-10-01 08:00:00"},
		"filter.updatedDateTo":   {"2026-10-01 08:00:00"},
		"filter.fromDate":        {"2026-10-01 08:00:00"},
		"filter.toDate":          {"2026-10-01 08:00:00"},
		"filter.userIds":         {"1,2"},
		"filter.departmentIds":   {"3"},
		"filter.managerIds":      {"4,5"},
		"filter.absenceIds":      {"6"},
		"filter.vacationStatus":  {"AwaitingAttest"},
	}, request.Query)
}

func TestVacationRequestsListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, "GET", "/v1/Core/VacationRequests", "vacationrequests_list.json"))

	page, err := client.VacationRequests.List(context.Background(), VacationRequestListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	request := page.Items[0]
	assert.Equal(t, 1, request.ID)
	assert.Equal(t, NameID{ID: 1, Name: "John Doe"}, request.User)
	assert.Equal(t, NewDate(2026, time.October, 13), request.FromDate)
	assert.Equal(t, NewDate(2026, time.October, 14), *request.ToDate)
	assert.Equal(t, 8.0, request.Hours)
	assert.Empty(t, request.Description)
	assert.Equal(t, NameID{ID: 1, Name: "Semester"}, request.AbsenceProject)
	assert.Equal(t, "2026-10-13 09:40:57", request.AttestedDate.String())
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, request.AttestedByUser)
	assert.Nil(t, request.DeniedDate)
	assert.Nil(t, request.DeniedByUser)
	assert.Equal(t, "2026-10-13 09:40:57", request.UpdatedDate.String())
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, request.UpdatedByUser)
	assert.Equal(t, "2026-10-13 09:40:57", request.CreatedDate.String())
}

func TestVacationRequestsListAllWalksPages(t *testing.T) {
	client := newTestClient(t, serveFixture(t, "vacationrequests_list.json"))

	var ids []int
	for request, err := range client.VacationRequests.ListAll(context.Background(), VacationRequestListParams{}) {
		require.NoError(t, err)
		ids = append(ids, request.ID)
	}

	assert.Equal(t, []int{1}, ids)
}

func TestVacationStatusDecodesCaseInsensitively(t *testing.T) {
	var statuses []VacationStatus
	require.NoError(t, json.Unmarshal([]byte(`["attested","awaitingAttest","Denied","other"]`), &statuses))

	assert.Equal(t, []VacationStatus{VacationStatusAttested, VacationStatusAwaitingAttest, VacationStatusDenied, "other"}, statuses)
}
