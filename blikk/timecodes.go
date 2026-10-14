package blikk

import (
	"context"
	"iter"
)

const timeCodesPath = "v1/Admin/TimeCodes"

type TimeCodesService struct {
	client *Client
}

const (
	TimeCodeSortByName        = "name"
	TimeCodeSortByCreatedDate = "createdDate"
	TimeCodeSortByUpdatedDate = "updatedDate"
)

type TimeCodeListParams struct {
	ListOptions
}

type TimeCodeListItem struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isActive"`
	IsStandard  bool      `json:"isStandard"`
	CreatedBy   *NameID   `json:"createdBy"`
	UpdatedBy   *NameID   `json:"updatedBy"`
	Created     *DateTime `json:"created"`
	Updated     *DateTime `json:"updated"`
}

func (s *TimeCodesService) List(ctx context.Context, params TimeCodeListParams) (*Page[TimeCodeListItem], error) {
	return listPage[TimeCodeListItem](ctx, s.client, timeCodesPath, params)
}

func (s *TimeCodesService) ListAll(ctx context.Context, params TimeCodeListParams) iter.Seq2[TimeCodeListItem, error] {
	return listAll[TimeCodeListItem](ctx, s.client, timeCodesPath, params)
}
