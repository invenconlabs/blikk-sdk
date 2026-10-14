package blikk

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentDecodesDocsExample(t *testing.T) {
	var comment Comment
	require.NoError(t, json.Unmarshal([]byte(`{
		"objectName": "task.comment.list",
		"id": 1,
		"text": "Lorem ipsum dolar sit amet",
		"createdByUser": {"objectName": "generic.name_id", "id": 1, "name": "John Doe"},
		"updatedByUser": null,
		"createdDate": "2026-10-11 09:40:57 +00:00",
		"updatedDate": null
	}`), &comment))

	assert.Equal(t, 1, comment.ID)
	assert.Equal(t, "Lorem ipsum dolar sit amet", comment.Text)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, comment.CreatedByUser)
	assert.Nil(t, comment.UpdatedByUser)
	assert.Equal(t, "2026-10-11 09:40:57", comment.CreatedDate.String())
	assert.Nil(t, comment.UpdatedDate)
}

func TestTagDecodesDocsExample(t *testing.T) {
	var tag Tag
	require.NoError(t, json.Unmarshal([]byte(`{
		"objectName": "tag.list",
		"id": 1,
		"title": "Good customer",
		"color": "#41b18b",
		"updatedDate": null,
		"updatedBy": null,
		"createdDate": "2026-10-09 09:40:57 +00:00",
		"createdBy": {"objectName": "generic.name_id", "id": 1, "name": "John Doe"}
	}`), &tag))

	assert.Equal(t, 1, tag.ID)
	assert.Equal(t, "Good customer", tag.Title)
	assert.Equal(t, "#41b18b", tag.Color)
	assert.Equal(t, "2026-10-09 09:40:57", tag.CreatedDate.String())
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, tag.CreatedBy)
	assert.Nil(t, tag.UpdatedDate)
	assert.Nil(t, tag.UpdatedBy)
}

func TestVatCodeDecodesDocsExample(t *testing.T) {
	var vatCode VatCode
	require.NoError(t, json.Unmarshal([]byte(`{
		"objectName": "vat.code",
		"id": 42,
		"name": "VAT 1",
		"code": "3",
		"vat": 12.00,
		"isVatFree": false,
		"createdBy": {"objectName": "generic.name_id", "id": 123, "name": "John Doe"},
		"updatedBy": {"objectName": "generic.name_id", "id": 123, "name": "John Doe"},
		"created": "2026-10-09 09:40:57 +00:00",
		"updated": "2026-10-11 09:40:57 +00:00"
	}`), &vatCode))

	assert.Equal(t, 42, vatCode.ID)
	assert.Equal(t, "VAT 1", vatCode.Name)
	assert.Equal(t, "3", vatCode.Code)
	assert.InDelta(t, 12.0, vatCode.VAT, 0.0001)
	assert.False(t, vatCode.IsVatFree)
	assert.Equal(t, &NameID{ID: 123, Name: "John Doe"}, vatCode.CreatedBy)
	assert.Equal(t, "2026-10-09 09:40:57", vatCode.Created.String())
	assert.Equal(t, "2026-10-11 09:40:57", vatCode.Updated.String())
}
