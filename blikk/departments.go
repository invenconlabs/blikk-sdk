package blikk

import (
	"context"
	"iter"
)

const departmentsPath = "v1/Admin/Departments"

type DepartmentsService struct {
	client *Client
}

const (
	DepartmentSortByName        = "name"
	DepartmentSortByCreatedDate = "createdDate"
	DepartmentSortByUpdatedDate = "updatedDate"
)

type DepartmentListParams struct {
	ListOptions
}

type DepartmentListItem struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedBy *NameID   `json:"createdBy"`
	Created   *DateTime `json:"created"`
	UpdatedBy *NameID   `json:"updatedBy"`
	Updated   *DateTime `json:"updated"`
}

func (s *DepartmentsService) List(ctx context.Context, params DepartmentListParams) (*Page[DepartmentListItem], error) {
	return listPage[DepartmentListItem](ctx, s.client, departmentsPath, params)
}

func (s *DepartmentsService) ListAll(ctx context.Context, params DepartmentListParams) iter.Seq2[DepartmentListItem, error] {
	return listAll[DepartmentListItem](ctx, s.client, departmentsPath, params)
}
