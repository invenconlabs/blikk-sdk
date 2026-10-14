package blikk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountriesListDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "countries_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	page, err := client.Countries.List(context.Background(), CountryListParams{
		ListOptions: ListOptions{SortBy: CountrySortByName},
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/Countries", request.Path)
	assert.Equal(t, map[string][]string{"page": {"1"}, "sortBy": {"name"}}, request.Query)

	require.Len(t, page.Items, 2)
	assert.Equal(t, Country{ID: 1, Name: "Barbados", Alpha3Code: "BRB", Alpha2Code: "BB"}, page.Items[0])
	assert.Equal(t, Country{ID: 2, Name: "Belgien", Alpha3Code: "BEL", Alpha2Code: "BE"}, page.Items[1])
}

func TestCountriesListAllUsesListPath(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "countries_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	var ids []int
	for item, err := range client.Countries.ListAll(context.Background(), CountryListParams{}) {
		require.NoError(t, err)
		ids = append(ids, item.ID)
	}

	assert.Equal(t, []int{1, 2}, ids)
	assert.Equal(t, "/v1/Admin/Countries", handler.last().Path)
}

func TestCountriesGetByCodeDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "country.json")}
	client := newTestClient(t, handler.ServeHTTP)

	country, err := client.Countries.GetByCode(context.Background(), "SWE")

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "GET", request.Method)
	assert.Equal(t, "/v1/Admin/Countries", request.Path)
	assert.Equal(t, map[string][]string{"code": {"SWE"}}, request.Query)
	assert.Equal(t, &Country{ID: 1, Name: "Sverige", Alpha3Code: "SWE", Alpha2Code: "SE"}, country)
}
