package blikk

import (
	"context"
	"iter"
	"net/url"
	"strconv"
)

const costCentersPath = "v1/Admin/CostCenters"

type CostCentersService struct {
	client *Client
}

const (
	CostCenterSortByCreatedDate = "createdDate"
	CostCenterSortByUpdatedDate = "updatedDate"
	CostCenterSortByName        = "name"
	CostCenterSortByCode        = "code"
)

type CostCenterListParams struct {
	ListOptions
	CreatedDateFrom DateTime `query:"filter.createdDateFrom"`
	CreatedDateTo   DateTime `query:"filter.createdDateTo"`
	UpdatedDateFrom DateTime `query:"filter.updatedDateFrom"`
	UpdatedDateTo   DateTime `query:"filter.updatedDateTo"`
	Query           string   `query:"filter.query"`
}

type CostCenter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	IsActive  bool      `json:"isActive"`
	CreatedBy *NameID   `json:"createdBy"`
	UpdatedBy *NameID   `json:"updatedBy"`
	Created   *DateTime `json:"created"`
	Updated   *DateTime `json:"updated"`
}

type CostCenterInput struct {
	Name     string `json:"name"`
	Code     string `json:"code,omitempty"`
	IsActive *bool  `json:"isActive,omitempty"`
}

func (s *CostCentersService) List(ctx context.Context, params CostCenterListParams) (*Page[CostCenter], error) {
	return listPage[CostCenter](ctx, s.client, costCentersPath, params)
}

func (s *CostCentersService) ListAll(ctx context.Context, params CostCenterListParams) iter.Seq2[CostCenter, error] {
	return listAll[CostCenter](ctx, s.client, costCentersPath, params)
}

func (s *CostCentersService) Get(ctx context.Context, id int) (*CostCenter, error) {
	return getResource[CostCenter](ctx, s.client, costCentersPath, strconv.Itoa(id))
}

func (s *CostCentersService) GetByCode(ctx context.Context, code string) (*CostCenter, error) {
	return getResourceWithQuery[CostCenter](ctx, s.client, costCentersPath, url.Values{"code": {code}})
}

func (s *CostCentersService) Create(ctx context.Context, input CostCenterInput) (*CostCenter, error) {
	return createResource[CostCenter](ctx, s.client, input, costCentersPath)
}

func (s *CostCentersService) Update(ctx context.Context, id int, input CostCenterInput) (*CostCenter, error) {
	return updateResourceReturning[CostCenter](ctx, s.client, input, costCentersPath, strconv.Itoa(id))
}
