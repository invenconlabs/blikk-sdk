package blikk

import (
	"context"
	"iter"
)

const projectStatusesPath = "v1/Admin/ProjectStatuses"

type ProjectStatusesService struct {
	client *Client
}

const (
	ProjectStatusSortByName        = "name"
	ProjectStatusSortBySortOrder   = "sortOrder"
	ProjectStatusSortByCreatedDate = "createdDate"
	ProjectStatusSortByUpdatedDate = "updatedDate"
)

type ProjectStatusListParams struct {
	ListOptions
	CreatedDateFrom DateTime `query:"filter.createdDateFrom"`
	CreatedDateTo   DateTime `query:"filter.createdDateTo"`
	UpdatedDateFrom DateTime `query:"filter.updatedDateFrom"`
	UpdatedDateTo   DateTime `query:"filter.updatedDateTo"`
}

type ProjectStatusListItem struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	SortOrder   int       `json:"sortOrder"`
	IsCompleted bool      `json:"isCompleted"`
	IsStandard  bool      `json:"isStandard"`
	UpdatedDate *DateTime `json:"updatedDate"`
	UpdatedBy   *NameID   `json:"updatedBy"`
	CreatedDate DateTime  `json:"createdDate"`
	CreatedBy   *NameID   `json:"createdBy"`
}

func (s *ProjectStatusesService) List(ctx context.Context, params ProjectStatusListParams) (*Page[ProjectStatusListItem], error) {
	return listPage[ProjectStatusListItem](ctx, s.client, projectStatusesPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *ProjectStatusesService) ListAll(ctx context.Context, params ProjectStatusListParams) iter.Seq2[ProjectStatusListItem, error] {
	return listAll[ProjectStatusListItem](ctx, s.client, projectStatusesPath, params)
}
