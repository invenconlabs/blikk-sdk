package blikk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func threeUserPages(failOnPage int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == failOnPage {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"boom","requestId":"r-1"}`))
			return
		}
		first, second := page*10+1, page*10+2
		fmt.Fprintf(w, `{"objectName":"generic.pagedlist","page":%d,"pageSize":2,"itemCount":2,"totalItemCount":6,"totalPages":3,"items":[{"id":%d},{"id":%d}]}`, page, first, second)
	}
}

func collectUserIDs(t *testing.T, client *Client, params UserListParams) ([]int, []error) {
	t.Helper()
	var ids []int
	var errs []error
	for user, err := range client.Users.ListAll(context.Background(), params) {
		if err != nil {
			errs = append(errs, err)
			continue
		}
		ids = append(ids, user.ID)
	}
	return ids, errs
}

func TestListAllWalksEveryPage(t *testing.T) {
	handler := &recordingHandler{respond: threeUserPages(0)}
	client := newTestClient(t, handler.ServeHTTP)

	ids, errs := collectUserIDs(t, client, UserListParams{Query: "anders"})

	assert.Empty(t, errs)
	assert.Equal(t, []int{11, 12, 21, 22, 31, 32}, ids)
	require.Len(t, handler.requests, 3)
	for i, request := range handler.requests {
		assert.Equal(t, strconv.Itoa(i+1), request.Query["page"][0])
		assert.Equal(t, strconv.Itoa(MaxPageSize), request.Query["pageSize"][0])
		assert.Equal(t, "anders", request.Query["filter.query"][0])
	}
}

func TestListAllKeepsCallerPageSizeAndStartPage(t *testing.T) {
	handler := &recordingHandler{respond: threeUserPages(0)}
	client := newTestClient(t, handler.ServeHTTP)

	ids, _ := collectUserIDs(t, client, UserListParams{ListOptions: ListOptions{Page: 2, PageSize: 2}})

	assert.Equal(t, []int{21, 22, 31, 32}, ids)
	assert.Equal(t, "2", handler.requests[0].Query["pageSize"][0])
}

func TestListAllStopsFetchingWhenCallerBreaks(t *testing.T) {
	handler := &recordingHandler{respond: threeUserPages(0)}
	client := newTestClient(t, handler.ServeHTTP)

	var ids []int
	for user, err := range client.Users.ListAll(context.Background(), UserListParams{}) {
		require.NoError(t, err)
		ids = append(ids, user.ID)
		if len(ids) == 3 {
			break
		}
	}

	assert.Equal(t, []int{11, 12, 21}, ids)
	assert.Len(t, handler.requests, 2)
}

func TestListAllYieldsErrorOnceAndStops(t *testing.T) {
	handler := &recordingHandler{respond: threeUserPages(2)}
	client := newTestClient(t, handler.ServeHTTP)

	ids, errs := collectUserIDs(t, client, UserListParams{})

	assert.Equal(t, []int{11, 12}, ids)
	require.Len(t, errs, 1)
	var apiErr *Error
	require.ErrorAs(t, errs[0], &apiErr)
	assert.Equal(t, "r-1", apiErr.RequestID)
	assert.Len(t, handler.requests, 2)
}

func TestListAllHandlesEmptyResult(t *testing.T) {
	client := newTestClient(t, respondWith(http.StatusOK, `{"page":1,"pageSize":100,"itemCount":0,"totalItemCount":0,"totalPages":0,"items":[]}`))

	ids, errs := collectUserIDs(t, client, UserListParams{})

	assert.Empty(t, ids)
	assert.Empty(t, errs)
}

func TestListAllStopsWhenContextIsCancelled(t *testing.T) {
	handler := &recordingHandler{respond: threeUserPages(0)}
	client := newTestClient(t, handler.ServeHTTP)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var errs []error
	count := 0
	for _, err := range client.Users.ListAll(ctx, UserListParams{}) {
		if err != nil {
			errs = append(errs, err)
			continue
		}
		count++
		cancel()
	}

	assert.Equal(t, 2, count)
	require.Len(t, errs, 1)
	assert.ErrorIs(t, errs[0], context.Canceled)
}

func TestListDefaultsPageToOneAndOmitsUnsetOptions(t *testing.T) {
	handler := &recordingHandler{respond: threeUserPages(0)}
	client := newTestClient(t, handler.ServeHTTP)

	page, err := client.Users.List(context.Background(), UserListParams{})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{"page": {"1"}}, handler.last().Query)
	assert.Equal(t, 1, page.Page)
	assert.Equal(t, 3, page.TotalPages)
	assert.Equal(t, 6, page.TotalItemCount)
	assert.Equal(t, 2, page.ItemCount)
	assert.Equal(t, 2, page.PageSize)
	assert.Len(t, page.Items, 2)
}

type deletedEntityFilter struct {
	EntityTypes []string `json:"entityTypes"`
}

func threeBodyPages(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Page int `json:"page"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	fmt.Fprintf(w, `{"objectName":"generic.pagedlist","page":%d,"pageSize":2,"itemCount":2,"totalItemCount":6,"totalPages":3,"items":[{"id":%d},{"id":%d}]}`, body.Page, body.Page*10+1, body.Page*10+2)
}

func TestListPageWithBodyPostsPagingInBody(t *testing.T) {
	client, handler := newRecordingClient(t, routes(t, map[string]http.HandlerFunc{"POST /v1/Core/DeletedEntities": threeBodyPages}))

	page, err := listPageWithBody[NameID](context.Background(), client, "v1/Core/DeletedEntities", pagedBody{Filter: deletedEntityFilter{EntityTypes: []string{"contact"}}})

	require.NoError(t, err)
	assert.Equal(t, []NameID{{ID: 11}, {ID: 12}}, page.Items)
	assert.Equal(t, 3, page.TotalPages)
	request := handler.last()
	assert.Empty(t, request.Query)
	assert.JSONEq(t, `{"filter":{"entityTypes":["contact"]},"page":1}`, string(request.Body))
}

func TestListPageWithBodyOmitsNilFilter(t *testing.T) {
	client, handler := newRecordingClient(t, threeBodyPages)

	_, err := listPageWithBody[NameID](context.Background(), client, "v1/Core/DeletedEntities", pagedBody{Page: 2, PageSize: 10})

	require.NoError(t, err)
	assert.JSONEq(t, `{"page":2,"pageSize":10}`, string(handler.last().Body))
}

func TestListAllWithBodyWalksEveryPage(t *testing.T) {
	client, handler := newRecordingClient(t, threeBodyPages)
	filter := deletedEntityFilter{EntityTypes: []string{"timereport"}}

	var ids []int
	for item, err := range listAllWithBody[NameID](context.Background(), client, "v1/Core/DeletedEntities", pagedBody{Filter: filter, Page: 2}) {
		require.NoError(t, err)
		ids = append(ids, item.ID)
	}

	assert.Equal(t, []int{21, 22, 31, 32}, ids)
	requests := handler.all()
	require.Len(t, requests, 2)
	for i, request := range requests {
		assert.Equal(t, http.MethodPost, request.Method)
		assert.JSONEq(t, fmt.Sprintf(`{"filter":{"entityTypes":["timereport"]},"page":%d,"pageSize":%d}`, i+2, MaxPageSize), string(request.Body))
	}
}

func TestListAllWithBodyStopsAtFirstError(t *testing.T) {
	client := newTestClient(t, respondWith(http.StatusBadRequest, `{"message":"bad filter"}`))

	var errs []error
	for _, err := range listAllWithBody[NameID](context.Background(), client, "v1/Core/DeletedEntities", pagedBody{}) {
		errs = append(errs, err)
	}

	require.Len(t, errs, 1)
	assert.ErrorContains(t, errs[0], "bad filter")
}
