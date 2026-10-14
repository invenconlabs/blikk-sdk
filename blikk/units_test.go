package blikk

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitsListDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "units_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	page, err := client.Units.List(context.Background(), UnitListParams{
		ListOptions: ListOptions{SortBy: UnitSortBySortOrder},
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/Units", request.Path)
	assert.Equal(t, map[string][]string{"page": {"1"}, "sortBy": {"sortOrder"}}, request.Query)

	require.Len(t, page.Items, 1)
	item := page.Items[0]
	assert.Equal(t, 5, item.ID)
	assert.Equal(t, "h", item.Name)
	assert.Zero(t, item.SortOrder)
	assert.False(t, item.IsStandard)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, item.CreatedBy)
	assert.Nil(t, item.UpdatedBy)
	assert.Equal(t, "2026-10-09 09:40:57", item.CreatedDate.String())
	assert.Nil(t, item.UpdatedDate)
}

func TestUnitsListAllUsesListPath(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "units_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	var ids []int
	for item, err := range client.Units.ListAll(context.Background(), UnitListParams{}) {
		require.NoError(t, err)
		ids = append(ids, item.ID)
	}

	assert.Equal(t, []int{5}, ids)
	assert.Equal(t, "/v1/Admin/Units", handler.last().Path)
}

func TestUnitsGetDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "unit.json")}
	client := newTestClient(t, handler.ServeHTTP)

	unit, err := client.Units.Get(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, "/v1/Admin/Units/42", handler.last().Path)
	assert.Equal(t, 42, unit.ID)
	assert.Equal(t, "Unit 1", unit.Name)
	assert.Equal(t, 1, unit.SortOrder)
	assert.True(t, unit.IsStandard)
	assert.Equal(t, &NameID{ID: 123, Name: "John Doe"}, unit.UpdatedBy)
	assert.Equal(t, "2026-10-11 09:40:57", unit.UpdatedDate.String())
}

func TestUnitsGetByNameSendsNameQuery(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "unit.json")}
	client := newTestClient(t, handler.ServeHTTP)

	unit, err := client.Units.GetByName(context.Background(), "Unit 1")

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/Units", request.Path)
	assert.Equal(t, map[string][]string{"name": {"Unit 1"}}, request.Query)
	assert.Equal(t, 42, unit.ID)
}

func TestUnitsCreateSendsBodyAndReturnsUnit(t *testing.T) {
	handler := &recordingHandler{respond: func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(loadFixture(t, "unit.json"))
	}}
	client := newTestClient(t, handler.ServeHTTP)

	unit, err := client.Units.Create(context.Background(), UnitInput{Name: "Unit 1", SortOrder: 1, IsStandard: true})

	require.NoError(t, err)
	assert.Equal(t, 42, unit.ID)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/Admin/Units", request.Path)
	assert.Equal(t, map[string]any{"name": "Unit 1", "sortOrder": float64(1), "isStandard": true}, request.jsonBody(t))
}

func TestUnitsCreateMinimalInputKeepsRequiredFields(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "unit.json")}
	client := newTestClient(t, handler.ServeHTTP)

	_, err := client.Units.Create(context.Background(), UnitInput{Name: "Unit 1"})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"name": "Unit 1", "isStandard": false}, handler.last().jsonBody(t))
}

func TestUnitsUpdateSendsPutAndReturnsUnit(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "unit.json")}
	client := newTestClient(t, handler.ServeHTTP)

	unit, err := client.Units.Update(context.Background(), 42, UnitInput{Name: "Unit 1", SortOrder: 1, IsStandard: true})

	require.NoError(t, err)
	assert.Equal(t, 42, unit.ID)
	request := handler.last()
	assert.Equal(t, http.MethodPut, request.Method)
	assert.Equal(t, "/v1/Admin/Units/42", request.Path)
	assert.Equal(t, map[string]any{"name": "Unit 1", "sortOrder": float64(1), "isStandard": true}, request.jsonBody(t))
}

func TestUnitsUpdateMinimalInputKeepsRequiredFields(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "unit.json")}
	client := newTestClient(t, handler.ServeHTTP)

	_, err := client.Units.Update(context.Background(), 42, UnitInput{Name: "Unit 1"})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"name": "Unit 1", "isStandard": false}, handler.last().jsonBody(t))
}
