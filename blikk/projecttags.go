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

type ProjectTagListParams struct {
	ListOptions
}

func (s *ProjectTagsService) List(ctx context.Context, params ProjectTagListParams) (*Page[Tag], error) {
	return listPage[Tag](ctx, s.client, projectTagsPath, params)
}

func (s *ProjectTagsService) ListAll(ctx context.Context, params ProjectTagListParams) iter.Seq2[Tag, error] {
	return listAll[Tag](ctx, s.client, projectTagsPath, params)
}
