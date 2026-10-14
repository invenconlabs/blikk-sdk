package blikk

import (
	"context"
	"iter"
	"strconv"
)

const offersPath = "v1/Core/Offers"

type OffersService struct {
	client *Client
}

type OfferState string

const (
	OfferStateDraft            OfferState = "Draft"
	OfferStateLocked           OfferState = "Locked"
	OfferStateSentToCustomer   OfferState = "SentToCustomer"
	OfferStateOpenedByCustomer OfferState = "OpenedByCustomer"
	OfferStateAccepted         OfferState = "Accepted"
	OfferStateAcceptedSigned   OfferState = "AcceptedSigned"
	OfferStateDenied           OfferState = "Denied"
)

func (s *OfferState) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, s, OfferStateDraft, OfferStateLocked, OfferStateSentToCustomer, OfferStateOpenedByCustomer, OfferStateAccepted, OfferStateAcceptedSigned, OfferStateDenied)
}

const (
	OfferSortByTitle       = "title"
	OfferSortByOfferNumber = "offerNumber"
	OfferSortByCreatedDate = "createdDate"
	OfferSortByUpdatedDate = "updatedDate"
)

type OfferListParams struct {
	ListOptions
	OfferNumber     int        `query:"filter.offerNumber"`
	CreatedDateFrom DateTime   `query:"filter.createdDateFrom"`
	CreatedDateTo   DateTime   `query:"filter.createdDateTo"`
	UpdatedDateFrom DateTime   `query:"filter.updatedDateFrom"`
	UpdatedDateTo   DateTime   `query:"filter.updatedDateTo"`
	OfferState      OfferState `query:"filter.offerState"`
	OpportunityID   int        `query:"filter.opportunityId"`
}

type Offer struct {
	ID            int        `json:"id"`
	OfferNumber   int        `json:"offerNumber"`
	Version       int        `json:"version"`
	Title         string     `json:"title"`
	OfferValue    float64    `json:"offerValue"`
	HasPdf        bool       `json:"hasPdf"`
	User          NameID     `json:"user"`
	Customer      *NameID    `json:"customer"`
	Opportunity   *NameID    `json:"opportunity"`
	OfferState    OfferState `json:"offerState"`
	UpdatedDate   *DateTime  `json:"updatedDate"`
	UpdatedByUser *NameID    `json:"updatedByUser"`
	CreatedDate   DateTime   `json:"createdDate"`
	CreatedByUser *NameID    `json:"createdByUser"`
}

func (s *OffersService) List(ctx context.Context, params OfferListParams) (*Page[Offer], error) {
	return listPage[Offer](ctx, s.client, offersPath, params)
}

func (s *OffersService) ListAll(ctx context.Context, params OfferListParams) iter.Seq2[Offer, error] {
	return listAll[Offer](ctx, s.client, offersPath, params)
}

// Download streams the offer PDF; the caller must close Content.
func (s *OffersService) Download(ctx context.Context, offerID int) (*Download, error) {
	path, err := resourcePath(offersPath, strconv.Itoa(offerID), "Download")
	if err != nil {
		return nil, err
	}
	return s.client.download(ctx, path, nil)
}
