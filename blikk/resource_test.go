package blikk

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetResourceWithQuerySendsQueryAndDecodesSingleObject(t *testing.T) {
	client, handler := newRecordingClient(t, routes(t, map[string]http.HandlerFunc{
		"GET /v1/admin/costcenters": respondWith(http.StatusOK, `{"objectName":"costcenter.list","id":42,"name":"Cost Center 42","code":"xyz"}`),
	}))

	costCenter, err := getResourceWithQuery[CostCenterRef](context.Background(), client, "v1/admin/costcenters", url.Values{"code": {"x y&z"}})

	require.NoError(t, err)
	assert.Equal(t, &CostCenterRef{ID: 42, Name: "Cost Center 42", Code: "xyz"}, costCenter)
	assert.Equal(t, map[string][]string{"code": {"x y&z"}}, handler.last().Query)
}

func TestCreateResourceIDAcceptsLowerAndUpperCaseID(t *testing.T) {
	for _, body := range []string{`{"id":42}`, `{"Id":42}`, `{"id":42,"nextProjectNumber":"P-124"}`} {
		client, handler := newRecordingClient(t, respondWith(http.StatusCreated, body))

		id, err := createResourceID(context.Background(), client, map[string]string{"title": "Good customer"}, "v1/Admin/ContactTags")

		require.NoError(t, err, body)
		assert.Equal(t, 42, id, body)
		request := handler.last()
		assert.Equal(t, http.MethodPost, request.Method)
		assert.Equal(t, "/v1/Admin/ContactTags", request.Path)
		assert.Equal(t, map[string]any{"title": "Good customer"}, request.jsonBody(t))
	}
}

func TestCreateResourceIDRejectsObjectWithoutID(t *testing.T) {
	client := newTestClient(t, respondWith(http.StatusOK, `{}`))

	_, err := createResourceID(context.Background(), client, struct{}{}, "v1/Admin/ContactTags")

	assert.ErrorContains(t, err, "has no id")
}

func TestCreateResourceIDReturnsZeroForEmptyBody(t *testing.T) {
	for _, body := range []string{``, " \n", `null`} {
		client := newTestClient(t, respondWith(http.StatusOK, body))

		id, err := createResourceID(context.Background(), client, struct{}{}, "v1/Core/MaterialReports")

		require.NoError(t, err, "body %q", body)
		assert.Zero(t, id, "body %q", body)
	}
}

func TestPatchResourceSendsQueryAndBody(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusNoContent, ""))

	err := patchResource(context.Background(), client, "v1/Admin/ContactTags/42", url.Values{"date": {"2025-01-01"}}, map[string]any{"colorId": 3})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodPatch, request.Method)
	assert.Equal(t, "/v1/Admin/ContactTags/42", request.Path)
	assert.Equal(t, map[string][]string{"date": {"2025-01-01"}}, request.Query)
	assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
	assert.Equal(t, map[string]any{"colorId": float64(3)}, request.jsonBody(t))
}

func TestPatchResourceWithoutBodySendsNoContentType(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusNoContent, ""))

	err := patchResource(context.Background(), client, "v1/core/invoices/42/setsenttoeconomysystem", url.Values{"date": {"2025-01-01"}}, nil)

	require.NoError(t, err)
	request := handler.last()
	assert.Empty(t, request.Body)
	assert.Empty(t, request.Header.Get("Content-Type"))
}

func TestUpdateResourceReturningDecodesResponse(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"id":42,"name":"Renamed","code":"xyz"}`))

	updated, err := updateResourceReturning[CostCenterRef](context.Background(), client, map[string]string{"name": "Renamed"}, "v1/admin/costcenters", "42")

	require.NoError(t, err)
	assert.Equal(t, &CostCenterRef{ID: 42, Name: "Renamed", Code: "xyz"}, updated)
	assert.Equal(t, http.MethodPut, handler.last().Method)
	assert.Equal(t, "/v1/admin/costcenters/42", handler.last().Path)
}
