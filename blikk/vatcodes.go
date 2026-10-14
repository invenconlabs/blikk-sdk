package blikk

import (
	"context"
	"iter"
	"net/url"
	"strconv"
)

const vatCodesPath = "v1/Admin/VatCodes"

type VatCodesService struct {
	client *Client
}

const (
	VatCodeSortByName      = "name"
	VatCodeSortByCode      = "code"
	VatCodeSortByVat       = "vat"
	VatCodeSortByIsVatFree = "isVatFree"
)

type VatCodeListParams struct {
	ListOptions
}

func (s *VatCodesService) List(ctx context.Context, params VatCodeListParams) (*Page[VatCode], error) {
	return listPage[VatCode](ctx, s.client, vatCodesPath, params)
}

func (s *VatCodesService) ListAll(ctx context.Context, params VatCodeListParams) iter.Seq2[VatCode, error] {
	return listAll[VatCode](ctx, s.client, vatCodesPath, params)
}

func (s *VatCodesService) Get(ctx context.Context, id int) (*VatCode, error) {
	return getResource[VatCode](ctx, s.client, vatCodesPath, strconv.Itoa(id))
}

func (s *VatCodesService) GetByCode(ctx context.Context, code string) (*VatCode, error) {
	return getResourceWithQuery[VatCode](ctx, s.client, vatCodesPath, url.Values{"code": {code}})
}
