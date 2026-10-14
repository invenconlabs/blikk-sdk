package blikk

import (
	"context"
	"iter"
	"net/url"
	"strconv"
)

const unitsPath = "v1/Admin/Units"

type UnitsService struct {
	client *Client
}

const (
	UnitSortByName       = "name"
	UnitSortBySortOrder  = "sortOrder"
	UnitSortByIsStandard = "isStandard"
)

type UnitListParams struct {
	ListOptions
}

type Unit struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	SortOrder   int       `json:"sortOrder"`
	IsStandard  bool      `json:"isStandard"`
	CreatedBy   *NameID   `json:"createdBy"`
	UpdatedBy   *NameID   `json:"updatedBy"`
	CreatedDate *DateTime `json:"createdDate"`
	UpdatedDate *DateTime `json:"updatedDate"`
}

type UnitInput struct {
	Name       string `json:"name"`
	SortOrder  int    `json:"sortOrder,omitempty"`
	IsStandard bool   `json:"isStandard"`
}

func (s *UnitsService) List(ctx context.Context, params UnitListParams) (*Page[Unit], error) {
	return listPage[Unit](ctx, s.client, unitsPath, params)
}

func (s *UnitsService) ListAll(ctx context.Context, params UnitListParams) iter.Seq2[Unit, error] {
	return listAll[Unit](ctx, s.client, unitsPath, params)
}

func (s *UnitsService) Get(ctx context.Context, id int) (*Unit, error) {
	return getResource[Unit](ctx, s.client, unitsPath, strconv.Itoa(id))
}

func (s *UnitsService) GetByName(ctx context.Context, name string) (*Unit, error) {
	return getResourceWithQuery[Unit](ctx, s.client, unitsPath, url.Values{"name": {name}})
}

func (s *UnitsService) Create(ctx context.Context, input UnitInput) (*Unit, error) {
	return createResource[Unit](ctx, s.client, input, unitsPath)
}

func (s *UnitsService) Update(ctx context.Context, id int, input UnitInput) (*Unit, error) {
	return updateResourceReturning[Unit](ctx, s.client, input, unitsPath, strconv.Itoa(id))
}
