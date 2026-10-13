package blikk

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeReportsListEncodesEveryFilter(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "timereports_list.json")}
	client := newTestClient(t, handler.ServeHTTP)
	moment := DateTimeOf(time.Date(2021, 1, 1, 14, 0, 0, 0, time.UTC))

	_, err := client.TimeReports.List(context.Background(), TimeReportListParams{
		ListOptions:           ListOptions{SortBy: TimeReportSortByDate, SortOrder: SortDescending},
		CreatedDateFrom:       moment,
		CreatedDateTo:         moment,
		UpdatedDateFrom:       moment,
		UpdatedDateTo:         moment,
		CreatedOrUpdatedSince: moment,
		From:                  NewDate(2021, time.October, 1),
		To:                    NewDate(2021, time.October, 31),
		UserIDs:               []int{1, 2},
		ProjectID:             3,
		InternalProjectID:     4,
		AbsenceProjectID:      5,
		ContactID:             6,
		ActivityID:            7,
		IsSentToAttest:        Ptr(true),
		IsAttested:            Ptr(false),
		ProjectType:           TimeReportProjectTypeAbsenceProject,
		HasAdditions:          Ptr(false),
		HasEquipment:          Ptr(true),
		HasAllowanceReports:   Ptr(false),
		HasTravelReports:      Ptr(true),
		TaskID:                8,
		InvoiceStatus:         []InvoiceStatus{InvoiceStatusPending, InvoiceStatusInvoiced},
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Core/TimeReports", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                         {"1"},
		"sortBy":                       {"date"},
		"sortOrder":                    {"descending"},
		"filter.createdDateFrom":       {"2021-01-01 14:00:00"},
		"filter.createdDateTo":         {"2021-01-01 14:00:00"},
		"filter.updatedDateFrom":       {"2021-01-01 14:00:00"},
		"filter.updatedDateTo":         {"2021-01-01 14:00:00"},
		"filter.createdOrUpdatedSince": {"2021-01-01 14:00:00"},
		"filter.from":                  {"2021-10-01"},
		"filter.to":                    {"2021-10-31"},
		"filter.userIds":               {"1,2"},
		"filter.projectId":             {"3"},
		"filter.internalProjectId":     {"4"},
		"filter.absenceProjectId":      {"5"},
		"filter.contactId":             {"6"},
		"filter.activityId":            {"7"},
		"filter.isSentToAttest":        {"true"},
		"filter.isAttested":            {"false"},
		"filter.projectType":           {"AbsenceProject"},
		"filter.hasAdditions":          {"false"},
		"filter.hasEquipment":          {"true"},
		"filter.hasAllowanceReports":   {"false"},
		"filter.hasTravelReports":      {"true"},
		"filter.taskId":                {"8"},
		"filter.invoiceStatus":         {"Pending,Invoiced"},
	}, request.Query)
}

func TestTimeReportsListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixture(t, "timereports_list.json"))

	page, err := client.TimeReports.List(context.Background(), TimeReportListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	report := page.Items[0]
	assert.Equal(t, 270, report.ID)
	assert.Equal(t, NewDate(2021, time.August, 30), report.Date)
	assert.Equal(t, "07:00", report.ClockStart)
	assert.Equal(t, "17:00", report.ClockEnd)
	assert.Equal(t, 8.0, report.Hours)
	assert.Equal(t, 7.0, report.InvoiceableHours)
	assert.Equal(t, 60, report.BreakMinutes)
	assert.Equal(t, 200.0, report.Cost)
	assert.Equal(t, 500.0, report.Rate)
	assert.Zero(t, report.Discount)
	assert.Equal(t, "A comment that is visible on the invoice", report.Comment)
	assert.Equal(t, "A comment that is not visible on the invoice", report.InternalComment)
	assert.Equal(t, "2021-09-15 07:02:16", report.SentToAttestDate.String())
	assert.Equal(t, "2021-09-15 07:10:18", report.AttestedDate.String())
	assert.Equal(t, NameID{ID: 1, Name: "John Doe"}, report.User)
	assert.Equal(t, &ProjectRef{ID: 19, Name: "Maintenance", Number: "12"}, report.Project)
	assert.Nil(t, report.InternalProject)
	assert.Nil(t, report.AbsenceProject)
	assert.Equal(t, NameID{ID: 87, Name: "Microsoft"}, report.Contact)
	assert.Equal(t, NameID{ID: 34, Name: "Project management"}, report.Activity)
	assert.Equal(t, NameID{ID: 16, Name: "Normal"}, report.TimeCode)
	assert.Equal(t, NameID{ID: 49, Name: "Work hours"}, report.TimeArticle)
	assert.Equal(t, &CostCenterRef{ID: 460, Name: "Center"}, report.CostCenter)
	assert.Equal(t, 16012, report.InvoiceID)
	assert.Nil(t, report.InvoicedDate)
	assert.Zero(t, report.InvoiceDraftID)
	assert.Equal(t, 15747, report.TravelReportID)
	assert.Equal(t, 314, report.AllowanceReportID)
	assert.Equal(t, "2021-09-15 05:02:04", report.CreatedDate.String())
	assert.Equal(t, "2021-10-09 12:44:07", report.UpdatedDate.String())
	assert.False(t, report.HasAdditions)
	assert.True(t, report.HasEquipment)
	assert.Equal(t, &NameID{ID: 6, Name: "John Doe"}, report.CreatedBy)
	assert.Equal(t, &NameID{ID: 1416, Name: "Derek Johnson"}, report.UpdatedBy)
	assert.Zero(t, report.TaskID)
}

func TestTimeReportsGetDecodesMinimalFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "timereport.json")}
	client := newTestClient(t, handler.ServeHTTP)

	report, err := client.TimeReports.Get(context.Background(), 270)

	require.NoError(t, err)
	assert.Equal(t, "/v1/Core/TimeReports/270", handler.last().Path)
	assert.Equal(t, 270, report.ID)
	assert.Nil(t, report.AttestedByUser)
	assert.Nil(t, report.BookingProject)
	assert.Nil(t, report.Task)
	assert.Nil(t, report.TravelReport)
	assert.Nil(t, report.AllowanceReport)
	assert.Empty(t, report.Materials)
	assert.Empty(t, report.Equipments)
	assert.Empty(t, report.Additions)
	assert.Zero(t, report.BreakMinutes)
	assert.Empty(t, report.ClockStart)
}

func TestTimeReportsGetDecodesNestedReports(t *testing.T) {
	client := newTestClient(t, serveFixture(t, "timereport_full.json"))

	report, err := client.TimeReports.Get(context.Background(), 270)

	require.NoError(t, err)
	require.NotNil(t, report.TravelReport)
	travel := report.TravelReport
	assert.Equal(t, 270, travel.ID)
	assert.Equal(t, "piteå-luleå-piteå", travel.Place)
	assert.Equal(t, "piteå-luleå-piteå", travel.InvoiceablePlace)
	assert.Equal(t, 100.0, travel.Distance)
	assert.Equal(t, 150.0, travel.InvoiceableDistance)
	assert.True(t, travel.ToSalary)
	assert.True(t, travel.ToBeInvoiced)
	assert.Equal(t, 10.0, travel.Cost)
	assert.Equal(t, 15.0, travel.Rate)
	assert.Equal(t, &CostCenterRef{ID: 3, Name: "Cost center 3"}, travel.CostCenter)
	assert.Equal(t, &NameID{ID: 2, Name: "BP 2"}, travel.BookingProject)
	assert.Equal(t, &CompanyCar{ID: 12, Name: "Company car #1", TripStart: 15000, TripEnd: 15250, AddressStart: "My street 123, My city", AddressGoal: "My street 321, Another city", AddressEnd: "My street 123, My city"}, travel.CompanyCar)

	require.NotNil(t, report.AllowanceReport)
	allowance := report.AllowanceReport
	assert.Equal(t, 12, allowance.ID)
	assert.Equal(t, "Lunch", allowance.Description)
	assert.Equal(t, 200.0, allowance.Cost)
	assert.Equal(t, 250.0, allowance.Amount)
	assert.Equal(t, 250.0, allowance.InvoiceableAmount)
	assert.Equal(t, "Whole day", allowance.InvoiceableSpecification)
	assert.False(t, allowance.ToBeInvoiced)
	assert.Nil(t, allowance.Article)

	require.Len(t, report.Materials, 1)
	material := report.Materials[0]
	assert.Equal(t, "Breaker bar", material.Title)
	assert.Equal(t, 2.0, material.Units)
	assert.Equal(t, 100.0, material.CostPerUnit)
	assert.Equal(t, 125.0, material.UnitPrice)
	assert.Equal(t, 125.0, material.InvoiceableUnitPrice)
	assert.Equal(t, 25.0, material.TaxPercent)
	assert.Equal(t, 62.5, material.Tax)
	assert.Equal(t, 250.0, material.InvoiceableAmount)
	assert.Equal(t, &NameID{ID: 24, Name: "Breaker bar"}, material.Article)
	assert.Equal(t, "pcs", material.InvoiceableUnitText)
	assert.Equal(t, "0012", material.ArticleNumber)
	assert.False(t, material.IsRotRut)

	require.Len(t, report.Equipments, 1)
	equipment := report.Equipments[0]
	assert.Equal(t, 5.0, equipment.Hours)
	assert.Equal(t, 40.0, equipment.InvoiceableRate)
	assert.Equal(t, 5.0, equipment.InvoiceableHours)
	assert.Equal(t, NameID{ID: 45, Name: "Jackhammer"}, equipment.Article)
	assert.Equal(t, "Used the jackhammer to obliterate a concrete slab.", equipment.Comment)

	require.Len(t, report.Additions, 1)
	addition := report.Additions[0]
	assert.Equal(t, 2.0, addition.SalaryCount)
	assert.Equal(t, 2.0, addition.InvoiceableCount)
	assert.Equal(t, 400.0, addition.InvoiceablePrice)
	assert.Equal(t, NameID{ID: 78, Name: "Night time"}, addition.Article)
}

func TestTimeReportsCreateSendsBodyAndReturnsID(t *testing.T) {
	handler := &recordingHandler{respond: respondWith(http.StatusCreated, `{"id":707}`)}
	client := newTestClient(t, handler.ServeHTTP)

	id, err := client.TimeReports.Create(context.Background(), TimeReportInput{
		Date:             NewDate(2021, time.November, 10),
		ClockStart:       "10:00",
		Hours:            7.75,
		InvoiceableHours: 9,
		UserID:           1,
		ProjectID:        1,
		ActivityID:       1,
		TimeCodeID:       1,
		TimeArticleID:    1,
		Additions:        []AdditionInput{{ArticleID: 15, YesNo: true}},
	})

	require.NoError(t, err)
	assert.Equal(t, 707, id)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/Core/TimeReports", request.Path)
	assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"date": "2021-11-10",
		"clockStart": "10:00",
		"hours": 7.75,
		"invoiceableHours": 9,
		"userId": 1,
		"projectId": 1,
		"activityId": 1,
		"timeCodeId": 1,
		"timeArticleId": 1,
		"additions": [{"articleId": 15, "yesNo": true}]
	}`, string(request.Body))
}

func TestTimeReportsCreateWithoutIDInResponseIsAnError(t *testing.T) {
	for _, body := range []string{"", `{}`, `{"id":0}`} {
		client := newTestClient(t, respondWith(http.StatusOK, body))

		id, err := client.TimeReports.Create(context.Background(), TimeReportInput{Date: NewDate(2021, time.November, 10)})

		assert.ErrorContains(t, err, "no id", "body %q", body)
		assert.Zero(t, id)
	}
}

func TestTimeReportInputMatchesDocumentedRequest(t *testing.T) {
	var documented map[string]any
	require.NoError(t, json.Unmarshal(loadFixture(t, "timereport_input.json"), &documented))
	var input TimeReportInput
	require.NoError(t, json.Unmarshal(loadFixture(t, "timereport_input.json"), &input))

	encoded, err := json.Marshal(input)
	require.NoError(t, err)

	var reencoded map[string]any
	require.NoError(t, json.Unmarshal(encoded, &reencoded))
	documented["date"] = "2021-11-10"
	assert.Equal(t, withoutEmptyValues(documented), withoutEmptyValues(reencoded))
}

func withoutEmptyValues(value any) any {
	switch v := value.(type) {
	case map[string]any:
		cleaned := map[string]any{}
		for key, element := range v {
			if !isEmptyJSONValue(element) {
				cleaned[key] = withoutEmptyValues(element)
			}
		}
		return cleaned
	case []any:
		cleaned := make([]any, len(v))
		for i, element := range v {
			cleaned[i] = withoutEmptyValues(element)
		}
		return cleaned
	default:
		return value
	}
}

func isEmptyJSONValue(value any) bool {
	return value == nil || value == false || value == 0.0 || value == ""
}

func TestTimeReportInputOmitsUnsetOptionalFieldsAndKeepsRequiredZeros(t *testing.T) {
	encoded, err := json.Marshal(TimeReportInput{
		Date:         NewDate(2021, time.November, 10),
		TravelReport: &TravelReportInput{Place: "Piteå"},
		Materials:    []MaterialInput{{Title: "Text row", IsTextOnly: true}},
	})

	require.NoError(t, err)
	assert.JSONEq(t, `{
		"date": "2021-11-10",
		"hours": 0,
		"invoiceableHours": 0,
		"userId": 0,
		"activityId": 0,
		"timeCodeId": 0,
		"timeArticleId": 0,
		"travelReport": {"place": "Piteå", "distance": 0, "toSalary": false},
		"materials": [{"title": "Text row", "unitPrice": 0, "units": 0, "isTextOnly": true}]
	}`, string(encoded))
}

func TestProjectNamePicksWhicheverProjectIsSet(t *testing.T) {
	tests := []struct {
		name   string
		report TimeReportListItem
		want   string
	}{
		{"project", TimeReportListItem{Project: &ProjectRef{Name: "Roof"}}, "Roof"},
		{"internal project", TimeReportListItem{InternalProject: &NameID{Name: "Admin"}}, "Admin"},
		{"absence project", TimeReportListItem{AbsenceProject: &NameID{Name: "Vacation"}}, "Vacation"},
		{"none", TimeReportListItem{}, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.report.ProjectName())
			detail := TimeReport{Project: tc.report.Project, InternalProject: tc.report.InternalProject, AbsenceProject: tc.report.AbsenceProject}
			assert.Equal(t, tc.want, detail.ProjectName())
		})
	}
}

func TestTimeReportsUpdateSendsPut(t *testing.T) {
	handler := &recordingHandler{respond: respondWith(http.StatusOK, `{"id":707}`)}
	client := newTestClient(t, handler.ServeHTTP)

	err := client.TimeReports.Update(context.Background(), 707, TimeReportInput{Date: NewDate(2021, time.November, 10), Hours: 8})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodPut, request.Method)
	assert.Equal(t, "/v1/Core/TimeReports/707", request.Path)
	assert.JSONEq(t, `{"date":"2021-11-10","hours":8,"invoiceableHours":0,"userId":0,"activityId":0,"timeCodeId":0,"timeArticleId":0}`, string(request.Body))
}

func TestTimeReportsDeleteSendsDelete(t *testing.T) {
	handler := &recordingHandler{respond: respondWith(http.StatusNoContent, "")}
	client := newTestClient(t, handler.ServeHTTP)

	err := client.TimeReports.Delete(context.Background(), 1)

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodDelete, request.Method)
	assert.Equal(t, "/v1/Core/TimeReports/1", request.Path)
	assert.Empty(t, request.Body)
}

func TestUserDayStatisticsEncodesFiltersAndDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "userdaystatistics.json")}
	client := newTestClient(t, handler.ServeHTTP)

	page, err := client.TimeReports.ListUserDayStatistics(context.Background(), UserDayStatisticsParams{
		From:          NewDate(2021, time.October, 1),
		To:            NewDate(2021, time.October, 2),
		UserIDs:       []int{1},
		DepartmentIDs: []int{2, 3},
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Core/TimeReports/UserDayStatistics", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                 {"1"},
		"filter.from":          {"2021-10-01"},
		"filter.to":            {"2021-10-02"},
		"filter.userIds":       {"1"},
		"filter.departmentIds": {"2,3"},
	}, request.Query)

	require.Len(t, page.Items, 1)
	statistics := page.Items[0]
	assert.Equal(t, 1, statistics.UserID)
	assert.Equal(t, "John Doe", statistics.Name)
	assert.Equal(t, &NameID{ID: 1, Name: "Developers"}, statistics.Department)
	require.Len(t, statistics.Dates, 2)
	assert.Equal(t, NewDate(2021, time.October, 1), statistics.Dates[0].Date)
	assert.Equal(t, 8.5, statistics.Dates[0].ReportedHours)
	assert.Equal(t, 8.0, statistics.Dates[0].ScheduledHours)
	assert.Equal(t, "2026-10-13 09:40:57", statistics.Dates[0].LockedDate.String())
	assert.Nil(t, statistics.Dates[0].AttestedDate)
	assert.Nil(t, statistics.Dates[1].LockedDate)
}

func TestListAllUserDayStatisticsWalksPages(t *testing.T) {
	client := newTestClient(t, serveFixture(t, "userdaystatistics.json"))

	var names []string
	for statistics, err := range client.TimeReports.ListAllUserDayStatistics(context.Background(), UserDayStatisticsParams{}) {
		require.NoError(t, err)
		names = append(names, statistics.Name)
	}

	assert.Equal(t, []string{"John Doe"}, names)
}
