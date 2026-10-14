package blikk

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeletedEntitiesListPostsFilterInBody(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, http.MethodPost, "/v1/Core/DeletedEntities", "deletedentities_list.json"))

	_, err := client.DeletedEntities.List(context.Background(), DeletedEntityListParams{
		Page:     2,
		PageSize: 25,
		Filter: DeletedEntityFilter{
			DeletedFromDate: DateTimeOf(time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)),
			DeletedToDate:   DateTimeOf(time.Date(2026, 6, 30, 23, 59, 59, 0, time.UTC)),
			EntityTypes:     []EntityType{EntityTypeTimeReport, EntityTypeContact},
		},
	})

	require.NoError(t, err)
	request := handler.last()
	assert.Empty(t, request.Query)
	assert.JSONEq(t, `{
		"filter": {
			"deletedFromDate": "2026-06-01 00:00:00",
			"deletedToDate": "2026-06-30 23:59:59",
			"entityTypes": ["timereport", "contact"]
		},
		"page": 2,
		"pageSize": 25
	}`, string(request.Body))
}

func TestDeletedEntitiesListOmitsUnsetFilterFields(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixture(t, "deletedentities_list.json"))

	_, err := client.DeletedEntities.List(context.Background(), DeletedEntityListParams{})

	require.NoError(t, err)
	assert.JSONEq(t, `{"filter":{},"page":1}`, string(handler.last().Body))
}

func TestDeletedEntitiesListDecodesFixture(t *testing.T) {
	client, _ := newRecordingClient(t, serveFixture(t, "deletedentities_list.json"))

	page, err := client.DeletedEntities.List(context.Background(), DeletedEntityListParams{})

	require.NoError(t, err)
	assert.Equal(t, []DeletedEntity{
		{EntityType: EntityTypeTimeReport, EntityID: 1, DeletedDate: Ptr(NewDate(2026, time.October, 8)), DeletedByID: 1},
		{EntityType: EntityTypeContact, EntityID: 2, DeletedDate: Ptr(NewDate(2026, time.October, 9)), DeletedByID: 1},
	}, page.Items)
}

func TestDeletedEntitiesListAllWalksPagesWithSameFilter(t *testing.T) {
	client, handler := newRecordingClient(t, threeBodyPages)
	params := DeletedEntityListParams{Filter: DeletedEntityFilter{EntityTypes: []EntityType{EntityTypeProject}}}

	var ids []int
	for entity, err := range client.DeletedEntities.ListAll(context.Background(), params) {
		require.NoError(t, err)
		ids = append(ids, entity.EntityID)
	}

	assert.Len(t, ids, 6)
	for _, request := range handler.all() {
		assert.Equal(t, map[string]any{"entityTypes": []any{"project"}}, request.jsonBody(t)["filter"])
	}
}

func TestEntityTypeMatchesCapitalisedNames(t *testing.T) {
	var entityTypes []EntityType
	require.NoError(t, json.Unmarshal([]byte(`["Project","TimeReport","supplierinvoice","Unknown"]`), &entityTypes))

	assert.Equal(t, []EntityType{EntityTypeProject, EntityTypeTimeReport, EntityTypeSupplierInvoice, "Unknown"}, entityTypes)
}
