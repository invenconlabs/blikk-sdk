package blikk

import (
	"context"
	"iter"
)

const deletedEntitiesPath = "v1/Core/DeletedEntities"

type DeletedEntitiesService struct {
	client *Client
}

type EntityType string

// The API sends these lowercase but matches filter values case-insensitively.
const (
	EntityTypeTimeReport        EntityType = "timereport"
	EntityTypeArticle           EntityType = "article"
	EntityTypeContact           EntityType = "contact"
	EntityTypeInvoice           EntityType = "invoice"
	EntityTypeOffer             EntityType = "offer"
	EntityTypeOpportunity       EntityType = "opportunity"
	EntityTypeProjectCollection EntityType = "projectcollection"
	EntityTypeProject           EntityType = "project"
	EntityTypeTaskBoard         EntityType = "taskboard"
	EntityTypeTask              EntityType = "task"
	EntityTypeVacationRequest   EntityType = "vacationrequest"
	EntityTypeMaterialReport    EntityType = "materialreport"
	EntityTypePaymentPlan       EntityType = "paymentplan"
	EntityTypeSupplierInvoice   EntityType = "supplierinvoice"
)

func (t *EntityType) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, t,
		EntityTypeTimeReport, EntityTypeArticle, EntityTypeContact, EntityTypeInvoice, EntityTypeOffer,
		EntityTypeOpportunity, EntityTypeProjectCollection, EntityTypeProject, EntityTypeTaskBoard, EntityTypeTask,
		EntityTypeVacationRequest, EntityTypeMaterialReport, EntityTypePaymentPlan, EntityTypeSupplierInvoice)
}

type DeletedEntityListParams struct {
	Page     int
	PageSize int
	Filter   DeletedEntityFilter
}

type DeletedEntityFilter struct {
	DeletedFromDate DateTime     `json:"deletedFromDate,omitzero"`
	DeletedToDate   DateTime     `json:"deletedToDate,omitzero"`
	EntityTypes     []EntityType `json:"entityTypes,omitempty"`
}

type DeletedEntity struct {
	EntityType  EntityType `json:"entityType"`
	EntityID    int        `json:"entityId"`
	DeletedDate *Date      `json:"deletedDate"`
	DeletedByID int        `json:"deletedById"`
}

func (s *DeletedEntitiesService) List(ctx context.Context, params DeletedEntityListParams) (*Page[DeletedEntity], error) {
	return listPageWithBody[DeletedEntity](ctx, s.client, deletedEntitiesPath, params.pagedBody())
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *DeletedEntitiesService) ListAll(ctx context.Context, params DeletedEntityListParams) iter.Seq2[DeletedEntity, error] {
	return listAllWithBody[DeletedEntity](ctx, s.client, deletedEntitiesPath, params.pagedBody())
}

func (p DeletedEntityListParams) pagedBody() pagedBody {
	return pagedBody{Filter: p.Filter, Page: p.Page, PageSize: p.PageSize}
}
