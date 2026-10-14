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

type ProjectCategoryListParams struct {
	ListOptions
}

func (s *ProjectCategoriesService) List(ctx context.Context, params ProjectCategoryListParams) (*Page[ProjectCategory], error) {
	return listPage[ProjectCategory](ctx, s.client, projectCategoriesPath, params)
}

func (s *ProjectCategoriesService) ListAll(ctx context.Context, params ProjectCategoryListParams) iter.Seq2[ProjectCategory, error] {
	return listAll[ProjectCategory](ctx, s.client, projectCategoriesPath, params)
}
