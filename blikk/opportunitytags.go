package blikk

import (
	"context"
	"iter"
)

const opportunityTagsPath = "v1/Admin/OpportunityTags"

type OpportunityTagsService struct {
	client *Client
}

const (
	OpportunityTagSortByName        = "name"
	OpportunityTagSortByCreatedDate = "createdDate"
	OpportunityTagSortByUpdatedDate = "updatedDate"
)

type OpportunityTagListParams struct {
	ListOptions
}

func (s *OpportunityTagsService) List(ctx context.Context, params OpportunityTagListParams) (*Page[Tag], error) {
	return listPage[Tag](ctx, s.client, opportunityTagsPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *OpportunityTagsService) ListAll(ctx context.Context, params OpportunityTagListParams) iter.Seq2[Tag, error] {
	return listAll[Tag](ctx, s.client, opportunityTagsPath, params)
}
