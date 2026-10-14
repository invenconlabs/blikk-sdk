package blikk

import (
	"context"
	"iter"
	"net/url"
)

const countriesPath = "v1/Admin/Countries"

type CountriesService struct {
	client *Client
}

const CountrySortByName = "name"

type CountryListParams struct {
	ListOptions
}

type Country struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Alpha3Code string `json:"alpha3Code"`
	Alpha2Code string `json:"alpha2Code"`
}

func (s *CountriesService) List(ctx context.Context, params CountryListParams) (*Page[Country], error) {
	return listPage[Country](ctx, s.client, countriesPath, params)
}

func (s *CountriesService) ListAll(ctx context.Context, params CountryListParams) iter.Seq2[Country, error] {
	return listAll[Country](ctx, s.client, countriesPath, params)
}

func (s *CountriesService) GetByCode(ctx context.Context, code string) (*Country, error) {
	return getResourceWithQuery[Country](ctx, s.client, countriesPath, url.Values{"code": {code}})
}
