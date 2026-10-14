//go:build integration

package blikk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAbsenceProjects(t *testing.T) {
	client := integrationClient(t)

	page, err := client.AbsenceProjects.List(integrationContext(t), AbsenceProjectListParams{ListOptions: ListOptions{PageSize: 5, SortBy: AbsenceProjectSortByName}})

	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	assert.NotZero(t, page.Items[0].ID)
	assert.NotEmpty(t, page.Items[0].Name)
	assert.NotNil(t, page.Items[0].Created)
}

func TestIntegrationActivities(t *testing.T) {
	client := integrationClient(t)

	page, err := client.Activities.List(integrationContext(t), ActivityListParams{ListOptions: ListOptions{PageSize: 5, SortBy: ActivitySortByName}})

	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	assert.NotZero(t, page.Items[0].ID)
	assert.NotEmpty(t, page.Items[0].Name)
	assert.NotNil(t, page.Items[0].Created)
}

func TestIntegrationCompanyCars(t *testing.T) {
	client := integrationClient(t)

	page, err := client.CompanyCars.List(integrationContext(t), CompanyCarListParams{ListOptions: ListOptions{PageSize: 5, SortBy: CompanyCarSortByName}})

	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	assert.NotZero(t, page.Items[0].ID)
	assert.NotEmpty(t, page.Items[0].Name)
	assert.NotEmpty(t, page.Items[0].Type)
	assert.NotNil(t, page.Items[0].Created)
}

func TestIntegrationCostCenters(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.CostCenters.List(ctx, CostCenterListParams{ListOptions: ListOptions{PageSize: 5, SortBy: CostCenterSortByName}})
	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)
	assert.NotEmpty(t, listed.Name)
	assert.NotNil(t, listed.Created)

	byID, err := client.CostCenters.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, byID.ID)
	assert.Equal(t, listed.Name, byID.Name)

	if listed.Code == "" {
		t.Skip("first cost center has no code")
	}
	byCode, err := client.CostCenters.GetByCode(ctx, listed.Code)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, byCode.ID)
}

func TestIntegrationCountries(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.Countries.List(ctx, CountryListParams{ListOptions: ListOptions{PageSize: 5, SortBy: CountrySortByName}})
	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	assert.NotZero(t, page.Items[0].ID)
	assert.NotEmpty(t, page.Items[0].Name)
	assert.NotEmpty(t, page.Items[0].Alpha3Code)

	for _, code := range []string{"SWE", "SE"} {
		country, err := client.Countries.GetByCode(ctx, code)
		require.NoError(t, err, code)
		assert.Equal(t, "SWE", country.Alpha3Code)
		assert.Equal(t, "SE", country.Alpha2Code)
	}
}

func TestIntegrationDepartments(t *testing.T) {
	client := integrationClient(t)

	page, err := client.Departments.List(integrationContext(t), DepartmentListParams{ListOptions: ListOptions{PageSize: 5, SortBy: DepartmentSortByName}})

	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	assert.NotZero(t, page.Items[0].ID)
	assert.NotEmpty(t, page.Items[0].Name)
	assert.NotNil(t, page.Items[0].Created)
}

func TestIntegrationInternalProjects(t *testing.T) {
	client := integrationClient(t)

	page, err := client.InternalProjects.List(integrationContext(t), InternalProjectListParams{ListOptions: ListOptions{PageSize: 5, SortBy: InternalProjectSortByName}})

	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	assert.NotZero(t, page.Items[0].ID)
	assert.NotEmpty(t, page.Items[0].Name)
	assert.NotNil(t, page.Items[0].Created)
}

func TestIntegrationTimeCodes(t *testing.T) {
	client := integrationClient(t)

	page, err := client.TimeCodes.List(integrationContext(t), TimeCodeListParams{ListOptions: ListOptions{PageSize: 5, SortBy: TimeCodeSortByName}})

	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	assert.NotZero(t, page.Items[0].ID)
	assert.NotEmpty(t, page.Items[0].Name)
	assert.NotNil(t, page.Items[0].Created)
}

func TestIntegrationUnits(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.Units.List(ctx, UnitListParams{ListOptions: ListOptions{PageSize: 5, SortBy: UnitSortByName}})
	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)
	assert.NotEmpty(t, listed.Name)
	assert.NotNil(t, listed.CreatedDate)

	byID, err := client.Units.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, byID.ID)
	assert.Equal(t, listed.Name, byID.Name)

	byName, err := client.Units.GetByName(ctx, listed.Name)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, byName.ID)
}

func TestIntegrationVatCodes(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationContext(t)

	page, err := client.VatCodes.List(ctx, VatCodeListParams{ListOptions: ListOptions{PageSize: 5, SortBy: VatCodeSortByCode}})
	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	listed := page.Items[0]
	assert.NotZero(t, listed.ID)
	assert.NotEmpty(t, listed.Name)
	assert.NotEmpty(t, listed.Code)
	assert.NotNil(t, listed.Created)

	byID, err := client.VatCodes.Get(ctx, listed.ID)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, byID.ID)
	assert.Equal(t, listed.VAT, byID.VAT)

	byCode, err := client.VatCodes.GetByCode(ctx, listed.Code)
	require.NoError(t, err)
	assert.Equal(t, listed.ID, byCode.ID)
}
