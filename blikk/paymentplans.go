package blikk

import (
	"context"
	"iter"
)

const paymentPlansPath = "v1/Core/Paymentplans"

type PaymentPlansService struct {
	client *Client
}

const (
	PaymentPlanSortByDescription = "description"
	PaymentPlanSortByCreatedDate = "createdDate"
	PaymentPlanSortByUpdatedDate = "updatedDate"
)

type PaymentPlanListParams struct {
	ListOptions
	BatchID               string   `query:"filter.batchId"`
	CreatedDateFrom       DateTime `query:"filter.createdDateFrom"`
	CreatedDateTo         DateTime `query:"filter.createdDateTo"`
	UpdatedDateFrom       DateTime `query:"filter.updatedDateFrom"`
	UpdatedDateTo         DateTime `query:"filter.updatedDateTo"`
	CreatedOrUpdatedSince DateTime `query:"filter.createdOrUpdatedSince"`
	ProjectID             int      `query:"filter.projectId"`
}

type PaymentPlan struct {
	ID             int        `json:"id"`
	ArticleNumber  string     `json:"articleNumber"`
	Description    string     `json:"description"`
	Price          float64    `json:"price"`
	Units          float64    `json:"units"`
	BatchID        string     `json:"batchId"`
	Date           *DateTime  `json:"date"`
	InvoiceID      int        `json:"invoiceId"`
	InvoicedDate   *DateTime  `json:"invoicedDate"`
	InvoiceDraftID int        `json:"invoiceDraftId"`
	Article        *NameID    `json:"article"`
	Project        ProjectRef `json:"project"`
	CreatedDate    DateTime   `json:"createdDate"`
	CreatedBy      NameID     `json:"createdBy"`
	UpdatedDate    *DateTime  `json:"updatedDate"`
	UpdatedBy      *NameID    `json:"updatedBy"`
}

type PaymentPlanInput struct {
	ProjectID   int                          `json:"projectId"`
	Date        DateTime                     `json:"date,omitzero"`
	Amount      float64                      `json:"amount,omitempty"`
	Description string                       `json:"description,omitempty"`
	ArticleRows []PaymentPlanArticleRowInput `json:"articleRows,omitempty"`
}

type PaymentPlanArticleRowInput struct {
	ArticleNumber string  `json:"articleNumber,omitempty"`
	Description   string  `json:"description,omitempty"`
	Quantity      float64 `json:"quantity,omitempty"`
	UnitPrice     float64 `json:"unitPrice,omitempty"`
}

type PaymentPlanCreateResult struct {
	BatchID        string `json:"batchId"`
	PaymentPlanIDs []int  `json:"paymentplanIds"`
}

func (s *PaymentPlansService) List(ctx context.Context, params PaymentPlanListParams) (*Page[PaymentPlan], error) {
	return listPage[PaymentPlan](ctx, s.client, paymentPlansPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *PaymentPlansService) ListAll(ctx context.Context, params PaymentPlanListParams) iter.Seq2[PaymentPlan, error] {
	return listAll[PaymentPlan](ctx, s.client, paymentPlansPath, params)
}

// Create adds one invoicing occasion and returns the batch it created with one payment plan per row.
func (s *PaymentPlansService) Create(ctx context.Context, input PaymentPlanInput) (*PaymentPlanCreateResult, error) {
	return createResource[PaymentPlanCreateResult](ctx, s.client, input, paymentPlansPath)
}
