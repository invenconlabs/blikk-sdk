package blikk

import (
	"context"
	"iter"
)

const projectTagsPath = "v1/Admin/ProjectTags"

type ProjectTagsService struct {
	client *Client
}

const (
	ProjectTagSortByName        = "name"
	ProjectTagSortByCreatedDate = "createdDate"
)

func (s *ProjectTagsService) List(ctx context.Context, params ListOptions) (*Page[TagListItem], error) {
	return listPage[TagListItem](ctx, s.client, projectTagsPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *ProjectTagsService) ListAll(ctx context.Context, params ListOptions) iter.Seq2[TagListItem, error] {
	return listAll[TagListItem](ctx, s.client, projectTagsPath, params)
}
