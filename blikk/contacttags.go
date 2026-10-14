package blikk

import (
	"context"
	"iter"
	"strconv"
)

const contactTagsPath = "v1/Admin/ContactTags"

type ContactTagsService struct {
	client *Client
}

const (
	ContactTagSortByName        = "name"
	ContactTagSortByCreatedDate = "createdDate"
)

type ContactTagListParams struct {
	ListOptions
}

type ContactTagCreateInput struct {
	Title   string `json:"title,omitempty"`
	ColorID int    `json:"colorId"`
}

type ContactTagUpdateInput struct {
	Title   string `json:"title,omitempty"`
	ColorID int    `json:"colorId,omitempty"`
}

func (s *ContactTagsService) List(ctx context.Context, params ContactTagListParams) (*Page[TagListItem], error) {
	return listPage[TagListItem](ctx, s.client, contactTagsPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *ContactTagsService) ListAll(ctx context.Context, params ContactTagListParams) iter.Seq2[TagListItem, error] {
	return listAll[TagListItem](ctx, s.client, contactTagsPath, params)
}

// Create returns the ID of the new contact tag.
func (s *ContactTagsService) Create(ctx context.Context, input ContactTagCreateInput) (int, error) {
	return createResourceID(ctx, s.client, input, contactTagsPath)
}

// Update changes only the fields set in input.
func (s *ContactTagsService) Update(ctx context.Context, id int, input ContactTagUpdateInput) error {
	return patchContactTag(ctx, s.client, id, input)
}

// Delete fails while the tag is still used by a contact.
func (s *ContactTagsService) Delete(ctx context.Context, id int) error {
	return deleteResource(ctx, s.client, contactTagsPath, strconv.Itoa(id))
}

func patchContactTag(ctx context.Context, c *Client, id int, input ContactTagUpdateInput) error {
	path, err := resourcePath(contactTagsPath, strconv.Itoa(id))
	if err != nil {
		return err
	}
	return patchResource(ctx, c, path, nil, input)
}
