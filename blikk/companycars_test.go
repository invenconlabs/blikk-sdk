package blikk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompanyCarsListDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "companycars_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	page, err := client.CompanyCars.List(context.Background(), CompanyCarListParams{
		ListOptions: ListOptions{SortBy: CompanyCarSortByCreatedDate},
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/CompanyCars", request.Path)
	assert.Equal(t, map[string][]string{"page": {"1"}, "sortBy": {"createdDate"}}, request.Query)

	require.Len(t, page.Items, 1)
	item := page.Items[0]
	assert.Equal(t, 1, item.ID)
	assert.Equal(t, "Lorem ipsum", item.Name)
	assert.True(t, item.IsActive)
	assert.Equal(t, "BMW 320 GT", item.MakeAndModel)
	assert.Equal(t, "ABC123", item.RegistrationNumber)
	assert.Equal(t, "CompanyCar", item.Type)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, item.CreatedBy)
	assert.Nil(t, item.UpdatedBy)
	assert.Equal(t, "2026-10-03 09:40:57", item.Created.String())
	assert.Nil(t, item.Updated)
}

func TestCompanyCarsListAllUsesListPath(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "companycars_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	var ids []int
	for item, err := range client.CompanyCars.ListAll(context.Background(), CompanyCarListParams{}) {
		require.NoError(t, err)
		ids = append(ids, item.ID)
	}

	assert.Equal(t, []int{1}, ids)
	assert.Equal(t, "/v1/Admin/CompanyCars", handler.last().Path)
}
