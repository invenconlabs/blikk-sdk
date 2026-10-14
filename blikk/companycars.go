package blikk

import (
	"context"
	"iter"
)

const companyCarsPath = "v1/Admin/CompanyCars"

type CompanyCarsService struct {
	client *Client
}

const (
	CompanyCarSortByName        = "name"
	CompanyCarSortByCreatedDate = "createdDate"
	CompanyCarSortByUpdatedDate = "updatedDate"
)

type CompanyCarListParams struct {
	ListOptions
}

type CompanyCarListItem struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	IsActive           bool      `json:"isActive"`
	MakeAndModel       string    `json:"makeAndModel"`
	RegistrationNumber string    `json:"registrationNumber"`
	Type               string    `json:"type"`
	CreatedBy          *NameID   `json:"createdBy"`
	UpdatedBy          *NameID   `json:"updatedBy"`
	Created            *DateTime `json:"created"`
	Updated            *DateTime `json:"updated"`
}

func (s *CompanyCarsService) List(ctx context.Context, params CompanyCarListParams) (*Page[CompanyCarListItem], error) {
	return listPage[CompanyCarListItem](ctx, s.client, companyCarsPath, params)
}

func (s *CompanyCarsService) ListAll(ctx context.Context, params CompanyCarListParams) iter.Seq2[CompanyCarListItem, error] {
	return listAll[CompanyCarListItem](ctx, s.client, companyCarsPath, params)
}
