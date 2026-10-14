package blikk

import (
	"context"
	"iter"
	"strconv"
)

const opportunitiesPath = "v1/Core/Opportunities"

type OpportunitiesService struct {
	client *Client
}

type OpportunityState string

const (
	OpportunityStateOpen OpportunityState = "Open"
	OpportunityStateWon  OpportunityState = "Won"
	OpportunityStateLost OpportunityState = "Lost"
)

func (s *OpportunityState) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, s, OpportunityStateOpen, OpportunityStateWon, OpportunityStateLost)
}

const (
	OpportunitySortByTitle       = "title"
	OpportunitySortByCreatedDate = "createdDate"
	OpportunitySortByUpdatedDate = "updatedDate"
)

type OpportunityListParams struct {
	ListOptions
	CreatedDateFrom          DateTime         `query:"filter.createdDateFrom"`
	CreatedDateTo            DateTime         `query:"filter.createdDateTo"`
	UpdatedDateFrom          DateTime         `query:"filter.updatedDateFrom"`
	UpdatedDateTo            DateTime         `query:"filter.updatedDateTo"`
	ClosedDateFrom           DateTime         `query:"filter.closedDateFrom"`
	ClosedDateTo             DateTime         `query:"filter.closedDateTo"`
	EstimatedClosingDateFrom DateTime         `query:"filter.estimatedClosingDateFrom"`
	EstimatedClosingDateTo   DateTime         `query:"filter.estimatedClosingDateTo"`
	HasOffers                *bool            `query:"filter.hasOffers"`
	OpportunityStatusIDs     []int            `query:"filter.opportunityStatusIds"`
	OpportunityTagIDs        []int            `query:"filter.opportunityTagIds"`
	OpportunityState         OpportunityState `query:"filter.opportunityState"`
}

type OpportunityListItem struct {
	ID                   int              `json:"id"`
	Title                string           `json:"title"`
	Turnover             float64          `json:"turnover"`
	HasOffers            bool             `json:"hasOffers"`
	User                 NameID           `json:"user"`
	Customer             *NameID          `json:"customer"`
	Contact              *NameID          `json:"contact"`
	OpportunityState     OpportunityState `json:"opportunityState"`
	ClosedDate           *Date            `json:"closedDate"`
	EstimatedClosingDate *DateTime        `json:"estimatedClosingDate"`
	ClosedByUser         *NameID          `json:"closedByUser"`
	Status               *NameID          `json:"status"`
	UpdatedDate          *DateTime        `json:"updatedDate"`
	UpdatedByUser        *NameID          `json:"updatedByUser"`
	CreatedDate          DateTime         `json:"createdDate"`
	CreatedByUser        *NameID          `json:"createdByUser"`
	Tags                 []TagRef         `json:"tags"`
}

type Opportunity struct {
	ID                   int              `json:"id"`
	Title                string           `json:"title"`
	Description          string           `json:"description"`
	Turnover             float64          `json:"turnover"`
	Probability          float64          `json:"probability"`
	HasOffers            bool             `json:"hasOffers"`
	User                 NameID           `json:"user"`
	Customer             *NameID          `json:"customer"`
	Contact              *NameID          `json:"contact"`
	EstimatedClosingDate *DateTime        `json:"estimatedClosingDate"`
	OpportunityState     OpportunityState `json:"opportunityState"`
	ClosedDate           *Date            `json:"closedDate"`
	ClosedByUser         *NameID          `json:"closedByUser"`
	Project              *NameID          `json:"project"`
	Status               *NameID          `json:"status"`
	UpdatedDate          *DateTime        `json:"updatedDate"`
	UpdatedByUser        *NameID          `json:"updatedByUser"`
	CreatedDate          DateTime         `json:"createdDate"`
	CreatedByUser        *NameID          `json:"createdByUser"`
	Tags                 []TagRef         `json:"tags"`
}

type OpportunityCreateInput struct {
	CustomerID           int      `json:"customerId"`
	Title                string   `json:"title"`
	StatusID             int      `json:"statusId"`
	UserID               int      `json:"userId"`
	Turnover             float64  `json:"turnover,omitempty"`
	EstimatedClosingDate DateTime `json:"estimatedClosingDate,omitzero"`
	Description          string   `json:"description,omitempty"`
}

type OpportunityCreateResult struct {
	ID                   int      `json:"id"`
	CustomerID           int      `json:"customerId"`
	Title                string   `json:"title"`
	StatusID             int      `json:"statusId"`
	Turnover             float64  `json:"turnover"`
	EstimatedClosingDate *Date    `json:"estimatedClosingDate"`
	Description          string   `json:"description"`
	UserID               int      `json:"userId"`
	CreatedDate          DateTime `json:"createdDate"`
	CreatedByID          int      `json:"createdById"`
}

type OpportunityUpdateInput struct {
	Title                string   `json:"title"`
	UserID               int      `json:"userId"`
	Probability          int      `json:"probability"`
	Description          string   `json:"description,omitempty"`
	CustomerID           int      `json:"customerId,omitempty"`
	ContactID            int      `json:"contactId,omitempty"`
	Turnover             float64  `json:"turnover,omitempty"`
	EstimatedClosingDate DateTime `json:"estimatedClosingDate,omitzero"`
	TagIDs               []int    `json:"tagIds,omitempty"`
}

func (s *OpportunitiesService) List(ctx context.Context, params OpportunityListParams) (*Page[OpportunityListItem], error) {
	return listPage[OpportunityListItem](ctx, s.client, opportunitiesPath, params)
}

func (s *OpportunitiesService) ListAll(ctx context.Context, params OpportunityListParams) iter.Seq2[OpportunityListItem, error] {
	return listAll[OpportunityListItem](ctx, s.client, opportunitiesPath, params)
}

func (s *OpportunitiesService) Get(ctx context.Context, id int) (*Opportunity, error) {
	return getResource[Opportunity](ctx, s.client, opportunitiesPath, strconv.Itoa(id))
}

func (s *OpportunitiesService) Create(ctx context.Context, input OpportunityCreateInput) (*OpportunityCreateResult, error) {
	return createResource[OpportunityCreateResult](ctx, s.client, input, opportunitiesPath)
}

// Update replaces the whole opportunity; any field left unset in input is cleared.
func (s *OpportunitiesService) Update(ctx context.Context, id int, input OpportunityUpdateInput) error {
	return updateResource(ctx, s.client, input, opportunitiesPath, strconv.Itoa(id))
}
