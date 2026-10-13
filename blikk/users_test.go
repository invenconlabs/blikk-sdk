package blikk

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersListDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "users_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	page, err := client.Users.List(context.Background(), UserListParams{
		ListOptions:       ListOptions{SortBy: "name", SortOrder: SortAscending},
		Query:             "john",
		IncludeDeleted:    true,
		IncludeRestricted: true,
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/Users", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                     {"1"},
		"sortBy":                   {"name"},
		"sortOrder":                {"ascending"},
		"filter.query":             {"john"},
		"filter.includeDeleted":    {"true"},
		"filter.includeRestricted": {"true"},
	}, request.Query)

	require.Len(t, page.Items, 1)
	user := page.Items[0]
	assert.Equal(t, 2, user.ID)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Empty(t, user.Email)
	assert.Equal(t, NewDate(2021, time.February, 1), *user.StartDate)
	assert.Equal(t, NewDate(2022, time.February, 15), *user.EndDate)
	assert.Equal(t, "123456", user.EmployeeNumber)
	assert.Equal(t, &NameID{ID: 5, Name: "Default profile"}, user.TimeReportingProfile)
	assert.Equal(t, &NameID{ID: 1, Name: "Developers"}, user.Department)
	assert.Equal(t, LicenseTotal, user.License)
	assert.Equal(t, []string{"administrator", "sales"}, user.Permissions)
	assert.Equal(t, &NameID{ID: 1, Name: "Jane Doe"}, user.CreatedBy)
	assert.Nil(t, user.UpdatedBy)
	assert.Equal(t, "2026-09-23 09:40:57", user.CreatedDate.String())
	assert.Nil(t, user.UpdatedDate)
	assert.Equal(t, SalaryType(""), user.SalaryType)
	assert.Nil(t, user.CostCenter)
	assert.False(t, user.IsRestricted)
}

func TestUsersGetDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "user.json")}
	client := newTestClient(t, handler.ServeHTTP)

	user, err := client.Users.Get(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, "/v1/Admin/Users/1", handler.last().Path)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, LicenseTotal, user.License)
	assert.Equal(t, "john.doe@test.com", user.Email)
	assert.Empty(t, user.SocialSecurityNumber)
	assert.Equal(t, "+46701234567", user.MobilePhoneNumber)
	assert.Equal(t, "This is a note about this user", user.Note)
	assert.Equal(t, "2021-12-01", user.StartDate.String())
	assert.Nil(t, user.EndDate)
	assert.Equal(t, &NameID{ID: 1, Name: "Developers"}, user.Department)
	assert.Nil(t, user.CostCenter)
	assert.Equal(t, SalaryTypeMonthly, user.SalaryType)
	assert.Equal(t, 450.0, user.CostPerHour)
	assert.Equal(t, "12", user.EmployeeNumber)
	assert.Nil(t, user.Address)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, user.Manager)
	assert.Equal(t, &NameID{ID: 1, Name: "Full time"}, user.Schedule)
	assert.Equal(t, &NameID{ID: 1, Name: "Consulting"}, user.StandardTimeArticle)
	assert.Equal(t, &NameID{ID: 1, Name: "Project management"}, user.StandardActivity)
	assert.Equal(t, &NameID{ID: 1, Name: "Base reporting profile"}, user.TimeReportingProfile)
	assert.False(t, user.TimeBankEnabled)
	assert.Zero(t, user.CurrentTimeBank)
	assert.Nil(t, user.Tags)
	assert.Equal(t, []string{"administrator"}, user.Permissions)
	assert.Equal(t, "2021-12-15 09:24:13", user.CreatedDate.String())
	assert.Nil(t, user.UpdatedDate)
}

func TestUserAddressDecodes(t *testing.T) {
	var user User
	err := json.Unmarshal([]byte(`{"address":{"objectName":"generic.address","streetAddress":"Storgatan 1","additionalAddress":"c/o Doe","postalCode":"94131","city":"Piteå","state":null,"countryId":205,"countryName":"Sverige"},"tags":[{"id":3,"title":"Senior","color":"#fff"}]}`), &user)

	require.NoError(t, err)
	assert.Equal(t, &Address{StreetAddress: "Storgatan 1", AdditionalAddress: "c/o Doe", PostalCode: "94131", City: "Piteå", CountryID: 205, CountryName: "Sverige"}, user.Address)
	assert.Equal(t, []Tag{{ID: 3, Title: "Senior", Color: "#fff"}}, user.Tags)
}

func TestSalaryTypeDecodesNumbersFromListAndNamesFromDetail(t *testing.T) {
	tests := []struct {
		input string
		want  SalaryType
	}{
		{`0`, ""},
		{`1`, SalaryTypeNone},
		{`2`, SalaryTypeMonthly},
		{`3`, SalaryTypeHourly},
		{`9`, SalaryType("9")},
		{`"Monthly"`, SalaryTypeMonthly},
		{`"hourly"`, SalaryTypeHourly},
		{`"Commission"`, SalaryType("Commission")},
		{`null`, ""},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := SalaryType("previous")
			require.NoError(t, json.Unmarshal([]byte(tc.input), &got))
			assert.Equal(t, tc.want, got)
		})
	}
}
