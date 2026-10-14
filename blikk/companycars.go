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

type CompanyCarType string

const (
	CompanyCarTypeCompanyCar CompanyCarType = "CompanyCar"
	// The API's own spelling.
	CompanyCarTypeOfficalCar CompanyCarType = "OfficalCar"
)

func (t *CompanyCarType) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, t, CompanyCarTypeCompanyCar, CompanyCarTypeOfficalCar)
}

type CompanyCar struct {
	ID                 int            `json:"id"`
	Name               string         `json:"name"`
	IsActive           bool           `json:"isActive"`
	MakeAndModel       string         `json:"makeAndModel"`
	RegistrationNumber string         `json:"registrationNumber"`
	Type               CompanyCarType `json:"type"`
	CreatedBy          *NameID        `json:"createdBy"`
	UpdatedBy          *NameID        `json:"updatedBy"`
	Created            *DateTime      `json:"created"`
	Updated            *DateTime      `json:"updated"`
}

func (s *CompanyCarsService) List(ctx context.Context, params CompanyCarListParams) (*Page[CompanyCar], error) {
	return listPage[CompanyCar](ctx, s.client, companyCarsPath, params)
}

func (s *CompanyCarsService) ListAll(ctx context.Context, params CompanyCarListParams) iter.Seq2[CompanyCar, error] {
	return listAll[CompanyCar](ctx, s.client, companyCarsPath, params)
}
