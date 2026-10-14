package blikk

import (
	"context"
	"iter"
)

const projectCategoriesPath = "v1/Admin/ProjectCategories"

type ProjectCategoriesService struct {
	client *Client
}

const ProjectCategorySortByName = "name"

type ProjectCategoryListItem struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (s *ProjectCategoriesService) List(ctx context.Context, params ListOptions) (*Page[ProjectCategoryListItem], error) {
	return listPage[ProjectCategoryListItem](ctx, s.client, projectCategoriesPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *ProjectCategoriesService) ListAll(ctx context.Context, params ListOptions) iter.Seq2[ProjectCategoryListItem, error] {
	return listAll[ProjectCategoryListItem](ctx, s.client, projectCategoriesPath, params)
}
