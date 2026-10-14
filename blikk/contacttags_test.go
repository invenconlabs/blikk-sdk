package blikk

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContactTagsListEncodesParamsAndDecodesFixture(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, http.MethodGet, "/v1/Admin/ContactTags", "contacttags_list.json"))

	page, err := client.ContactTags.List(context.Background(), ContactTagListParams{
		ListOptions: ListOptions{Page: 2, PageSize: 5, SortBy: ContactTagSortByCreatedDate, SortOrder: SortDescending},
	})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"page":      {"2"},
		"pageSize":  {"5"},
		"sortBy":    {"createdDate"},
		"sortOrder": {"descending"},
	}, handler.last().Query)
	require.Len(t, page.Items, 1)
	tag := page.Items[0]
	assert.Equal(t, 1, tag.ID)
	assert.Equal(t, "Good customer", tag.Title)
	assert.Equal(t, "#41b18b", tag.Color)
	assert.Nil(t, tag.UpdatedDate)
	assert.Nil(t, tag.UpdatedBy)
	assert.Equal(t, "2026-09-28 07:16:57", tag.CreatedDate.String())
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, tag.CreatedBy)
}

func TestContactTagsCreateSendsBodyAndReturnsID(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"Id":42}`))

	id, err := client.ContactTags.Create(context.Background(), ContactTagCreateInput{Title: "Good customer", ColorID: 1})

	require.NoError(t, err)
	assert.Equal(t, 42, id)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/Admin/ContactTags", request.Path)
	assert.JSONEq(t, `{"title":"Good customer","colorId":1}`, string(request.Body))
}

func TestContactTagsCreateMinimalBodyKeepsRequiredColor(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"id":42}`))

	_, err := client.ContactTags.Create(context.Background(), ContactTagCreateInput{})

	require.NoError(t, err)
	assert.JSONEq(t, `{"colorId":0}`, string(handler.last().Body))
}

func TestContactTagsUpdateSendsPatch(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusNoContent, ""))

	err := client.ContactTags.Update(context.Background(), 42, ContactTagUpdateInput{Title: "Very good customer", ColorID: 2})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodPatch, request.Method)
	assert.Equal(t, "/v1/Admin/ContactTags/42", request.Path)
	assert.JSONEq(t, `{"title":"Very good customer","colorId":2}`, string(request.Body))
}

func TestContactTagsUpdateMinimalBodyOmitsUnsetFields(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusNoContent, ""))

	err := client.ContactTags.Update(context.Background(), 42, ContactTagUpdateInput{ColorID: 3})

	require.NoError(t, err)
	assert.JSONEq(t, `{"colorId":3}`, string(handler.last().Body))
}

func TestContactTagsDeleteSendsDelete(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusNoContent, ""))

	err := client.ContactTags.Delete(context.Background(), 42)

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodDelete, request.Method)
	assert.Equal(t, "/v1/Admin/ContactTags/42", request.Path)
}
