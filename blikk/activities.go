package blikk

import (
	"context"
	"iter"
)

const activitiesPath = "v1/Admin/Activities"

type ActivitiesService struct {
	client *Client
}

const (
	ActivitySortByName        = "name"
	ActivitySortByCreatedDate = "createdDate"
	ActivitySortByUpdatedDate = "updatedDate"
)

type ActivityListParams struct {
	ListOptions
}

type Activity struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	IsActive             bool      `json:"isActive"`
	PreferredTimeCode    *NameID   `json:"preferredTimeCode"`
	PreferredTimeArticle *NameID   `json:"preferredTimeArticle"`
	CreatedBy            *NameID   `json:"createdBy"`
	UpdatedBy            *NameID   `json:"updatedBy"`
	Created              *DateTime `json:"created"`
	Updated              *DateTime `json:"updated"`
}

func (s *ActivitiesService) List(ctx context.Context, params ActivityListParams) (*Page[Activity], error) {
	return listPage[Activity](ctx, s.client, activitiesPath, params)
}

func (s *ActivitiesService) ListAll(ctx context.Context, params ActivityListParams) iter.Seq2[Activity, error] {
	return listAll[Activity](ctx, s.client, activitiesPath, params)
}
