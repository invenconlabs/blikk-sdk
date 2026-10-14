package blikk

import (
	"context"
	"iter"
)

const opportunityStatusesPath = "v1/Admin/OpportunityStatuses"

type OpportunityStatusesService struct {
	client *Client
}

const (
	OpportunityStatusSortByName        = "name"
	OpportunityStatusSortBySortOrder   = "sortOrder"
	OpportunityStatusSortByCreatedDate = "createdDate"
	OpportunityStatusSortByUpdatedDate = "updatedDate"
)

type OpportunityStatusListParams struct {
	ListOptions
	CreatedDateFrom DateTime `query:"filter.createdDateFrom"`
	CreatedDateTo   DateTime `query:"filter.createdDateTo"`
	UpdatedDateFrom DateTime `query:"filter.updatedDateFrom"`
	UpdatedDateTo   DateTime `query:"filter.updatedDateTo"`
}

type OpportunityStatus struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	SortOrder   int       `json:"sortOrder"`
	UpdatedDate *DateTime `json:"updatedDate"`
	UpdatedBy   *NameID   `json:"updatedBy"`
	CreatedDate DateTime  `json:"createdDate"`
	CreatedBy   *NameID   `json:"createdBy"`
}

func (s *OpportunityStatusesService) List(ctx context.Context, params OpportunityStatusListParams) (*Page[OpportunityStatus], error) {
	return listPage[OpportunityStatus](ctx, s.client, opportunityStatusesPath, params)
}

func (s *OpportunityStatusesService) ListAll(ctx context.Context, params OpportunityStatusListParams) iter.Seq2[OpportunityStatus, error] {
	return listAll[OpportunityStatus](ctx, s.client, opportunityStatusesPath, params)
}
