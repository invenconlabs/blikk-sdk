package blikk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVatCodesListDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "vatcodes_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	page, err := client.VatCodes.List(context.Background(), VatCodeListParams{
		ListOptions: ListOptions{SortBy: VatCodeSortByVat},
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/VatCodes", request.Path)
	assert.Equal(t, map[string][]string{"page": {"1"}, "sortBy": {"vat"}}, request.Query)

	require.Len(t, page.Items, 1)
	item := page.Items[0]
	assert.Equal(t, 42, item.ID)
	assert.Equal(t, "VAT 1", item.Name)
	assert.Equal(t, "XYZ", item.Code)
	assert.Equal(t, 12.0, item.VAT)
	assert.False(t, item.IsVatFree)
	assert.Equal(t, &NameID{ID: 123, Name: "John Doe"}, item.CreatedBy)
	assert.Nil(t, item.UpdatedBy)
	assert.Equal(t, "2026-10-09 09:40:57", item.Created.String())
	assert.Nil(t, item.Updated)
}

func TestVatCodesListAllUsesListPath(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "vatcodes_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	var ids []int
	for item, err := range client.VatCodes.ListAll(context.Background(), VatCodeListParams{}) {
		require.NoError(t, err)
		ids = append(ids, item.ID)
	}

	assert.Equal(t, []int{42}, ids)
	assert.Equal(t, "/v1/Admin/VatCodes", handler.last().Path)
}

func TestVatCodesGetDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "vatcode.json")}
	client := newTestClient(t, handler.ServeHTTP)

	vatCode, err := client.VatCodes.Get(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, "/v1/Admin/VatCodes/42", handler.last().Path)
	assert.Equal(t, 42, vatCode.ID)
	assert.Equal(t, 12.0, vatCode.VAT)
	assert.Equal(t, &NameID{ID: 123, Name: "John Doe"}, vatCode.UpdatedBy)
	assert.Equal(t, "2026-10-11 09:40:57", vatCode.Updated.String())
}

func TestVatCodesGetByCodeSendsCodeQuery(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "vatcode.json")}
	client := newTestClient(t, handler.ServeHTTP)

	vatCode, err := client.VatCodes.GetByCode(context.Background(), "XYZ")

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/VatCodes", request.Path)
	assert.Equal(t, map[string][]string{"code": {"XYZ"}}, request.Query)
	assert.Equal(t, 42, vatCode.ID)
}
