package blikk

import (
	"context"
	"iter"
)

const internalProjectsPath = "v1/Admin/InternalProjects"

type InternalProjectsService struct {
	client *Client
}

const InternalProjectSortByName = "name"

type InternalProjectListParams struct {
	ListOptions
}

type InternalProject struct {
	ID                               int       `json:"id"`
	Name                             string    `json:"name"`
	CommentRequiredWhenTimeReporting bool      `json:"commentRequiredWhenTimeReporting"`
	IsActive                         bool      `json:"isActive"`
	CreatedBy                        *NameID   `json:"createdBy"`
	UpdatedBy                        *NameID   `json:"updatedBy"`
	Created                          *DateTime `json:"created"`
	Updated                          *DateTime `json:"updated"`
}

func (s *InternalProjectsService) List(ctx context.Context, params InternalProjectListParams) (*Page[InternalProject], error) {
	return listPage[InternalProject](ctx, s.client, internalProjectsPath, params)
}

func (s *InternalProjectsService) ListAll(ctx context.Context, params InternalProjectListParams) iter.Seq2[InternalProject, error] {
	return listAll[InternalProject](ctx, s.client, internalProjectsPath, params)
}
