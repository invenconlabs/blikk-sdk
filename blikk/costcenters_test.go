package blikk

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCostCentersListSendsFiltersAndDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "costcenters_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	page, err := client.CostCenters.List(context.Background(), CostCenterListParams{
		ListOptions:     ListOptions{SortBy: CostCenterSortByCode, SortOrder: SortAscending},
		CreatedDateFrom: DateTimeOf(time.Date(2026, time.October, 1, 8, 0, 0, 0, time.UTC)),
		CreatedDateTo:   DateTimeOf(time.Date(2026, time.October, 2, 8, 0, 0, 0, time.UTC)),
		UpdatedDateFrom: DateTimeOf(time.Date(2026, time.October, 3, 8, 0, 0, 0, time.UTC)),
		UpdatedDateTo:   DateTimeOf(time.Date(2026, time.October, 4, 8, 0, 0, 0, time.UTC)),
		Query:           "cost",
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/CostCenters", request.Path)
	assert.Equal(t, map[string][]string{
		"page":                   {"1"},
		"sortBy":                 {"code"},
		"sortOrder":              {"ascending"},
		"filter.createdDateFrom": {"2026-10-01 08:00:00"},
		"filter.createdDateTo":   {"2026-10-02 08:00:00"},
		"filter.updatedDateFrom": {"2026-10-03 08:00:00"},
		"filter.updatedDateTo":   {"2026-10-04 08:00:00"},
		"filter.query":           {"cost"},
	}, request.Query)

	require.Len(t, page.Items, 1)
	item := page.Items[0]
	assert.Equal(t, 42, item.ID)
	assert.Equal(t, "Cost Center 1", item.Name)
	assert.Equal(t, "123", item.Code)
	assert.True(t, item.IsActive)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, item.CreatedBy)
	assert.Nil(t, item.UpdatedBy)
	assert.Equal(t, "2026-10-09 09:40:57", item.Created.String())
	assert.Nil(t, item.Updated)
}

func TestCostCentersListAllUsesListPath(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "costcenters_list.json")}
	client := newTestClient(t, handler.ServeHTTP)

	var ids []int
	for item, err := range client.CostCenters.ListAll(context.Background(), CostCenterListParams{}) {
		require.NoError(t, err)
		ids = append(ids, item.ID)
	}

	assert.Equal(t, []int{42}, ids)
	assert.Equal(t, "/v1/Admin/CostCenters", handler.last().Path)
}

func TestCostCentersGetDecodesFixture(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "costcenter.json")}
	client := newTestClient(t, handler.ServeHTTP)

	costCenter, err := client.CostCenters.Get(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, "/v1/Admin/CostCenters/42", handler.last().Path)
	assert.Equal(t, 42, costCenter.ID)
	assert.Equal(t, "123", costCenter.Code)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, costCenter.UpdatedBy)
	assert.Equal(t, "2026-10-11 09:40:57", costCenter.Updated.String())
}

func TestCostCentersGetByCodeSendsCodeQuery(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "costcenter.json")}
	client := newTestClient(t, handler.ServeHTTP)

	costCenter, err := client.CostCenters.GetByCode(context.Background(), "123")

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, "/v1/Admin/CostCenters", request.Path)
	assert.Equal(t, map[string][]string{"code": {"123"}}, request.Query)
	assert.Equal(t, 42, costCenter.ID)
}

func TestCostCentersCreateSendsBodyAndReturnsID(t *testing.T) {
	handler := &recordingHandler{respond: func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(loadFixture(t, "costcenter.json"))
	}}
	client := newTestClient(t, handler.ServeHTTP)

	id, err := client.CostCenters.Create(context.Background(), CostCenterInput{Name: "Cost Center 1", Code: "123", IsActive: Ptr(false)})

	require.NoError(t, err)
	assert.Equal(t, 42, id)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/Admin/CostCenters", request.Path)
	assert.Equal(t, map[string]any{"name": "Cost Center 1", "code": "123", "isActive": false}, request.jsonBody(t))
}

func TestCostCentersCreateMinimalInputSendsOnlyName(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "costcenter.json")}
	client := newTestClient(t, handler.ServeHTTP)

	_, err := client.CostCenters.Create(context.Background(), CostCenterInput{Name: "Cost Center 1"})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"name": "Cost Center 1"}, handler.last().jsonBody(t))
}

func TestCostCentersUpdateSendsPutAndReturnsCostCenter(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "costcenter.json")}
	client := newTestClient(t, handler.ServeHTTP)

	costCenter, err := client.CostCenters.Update(context.Background(), 42, CostCenterInput{Name: "Cost Center 1", Code: "123", IsActive: Ptr(true)})

	require.NoError(t, err)
	assert.Equal(t, 42, costCenter.ID)
	request := handler.last()
	assert.Equal(t, http.MethodPut, request.Method)
	assert.Equal(t, "/v1/Admin/CostCenters/42", request.Path)
	assert.Equal(t, map[string]any{"name": "Cost Center 1", "code": "123", "isActive": true}, request.jsonBody(t))
}

func TestCostCentersUpdateMinimalInputSendsOnlyName(t *testing.T) {
	handler := &recordingHandler{respond: serveFixture(t, "costcenter.json")}
	client := newTestClient(t, handler.ServeHTTP)

	_, err := client.CostCenters.Update(context.Background(), 42, CostCenterInput{Name: "Cost Center 1"})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"name": "Cost Center 1"}, handler.last().jsonBody(t))
}
