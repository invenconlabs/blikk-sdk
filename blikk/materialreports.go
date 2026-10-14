package blikk

import (
	"context"
	"iter"
	"strconv"
)

const materialReportsPath = "v1/Core/MaterialReports"

type MaterialReportsService struct {
	client *Client
}

const (
	MaterialReportSortBySortOrder     = "sortOrder"
	MaterialReportSortByArticleNumber = "articleNumber"
	MaterialReportSortByTitle         = "title"
	MaterialReportSortByCreatedDate   = "createdDate"
	MaterialReportSortByUpdatedDate   = "updatedDate"
)

type MaterialReportListParams struct {
	ListOptions
	CreatedDateFrom       DateTime `query:"filter.createdDateFrom"`
	CreatedDateTo         DateTime `query:"filter.createdDateTo"`
	UpdatedDateFrom       DateTime `query:"filter.updatedDateFrom"`
	UpdatedDateTo         DateTime `query:"filter.updatedDateTo"`
	CreatedOrUpdatedSince DateTime `query:"filter.createdOrUpdatedSince"`
	Query                 string   `query:"filter.query"`
	OnlyText              bool     `query:"filter.onlyText"`
	OnlyInvoiced          bool     `query:"filter.onlyInvoiced"`
	ProjectID             int      `query:"filter.projectId"`
	TimeReportID          int      `query:"filter.timeReportId"`
}

type MaterialReport struct {
	ID                 int            `json:"id"`
	SortOrder          int            `json:"sortOrder"`
	ArticleNumber      string         `json:"articleNumber"`
	Title              string         `json:"title"`
	ConsumedUnits      float64        `json:"consumedUnits"`
	InvoiceableUnits   float64        `json:"invoiceableUnits"`
	UnitText           string         `json:"unitText"`
	CostPerUnit        float64        `json:"costPerUnit"`
	PricePerUnit       float64        `json:"pricePerUnit"`
	DiscountPercentage float64        `json:"discountPercentage"`
	VatPercentage      float64        `json:"vatPercentage"`
	TotalAmount        float64        `json:"totalAmount"`
	ToBeInvoiced       bool           `json:"toBeInvoiced"`
	InvoiceID          int            `json:"invoiceId"`
	InvoicedDate       *DateTime      `json:"invoicedDate"`
	InvoiceDraftID     int            `json:"invoiceDraftId"`
	TimeReportID       int            `json:"timeReportId"`
	IsTextOnly         bool           `json:"isTextOnly"`
	Article            *NameID        `json:"article"`
	CostCenter         *CostCenterRef `json:"costCenter"`
	BookingProject     *NameID        `json:"bookingProject"`
	Project            ProjectRef     `json:"project"`
	CreatedDate        DateTime       `json:"createdDate"`
	CreatedBy          NameID         `json:"createdBy"`
	UpdatedDate        *DateTime      `json:"updatedDate"`
	UpdatedBy          *NameID        `json:"updatedBy"`
}

type MaterialReportCreateInput struct {
	ArticleID int `json:"articleId"`
	// The API calls projects work orders.
	ProjectID        int      `json:"workOrderId"`
	Units            float64  `json:"units"`
	InvoiceableUnits float64  `json:"invoicableUnits"`
	Title            string   `json:"title,omitempty"`
	UnitText         string   `json:"unitText,omitempty"`
	Cost             float64  `json:"cost,omitempty"`
	UnitPrice        float64  `json:"unitPrice,omitempty"`
	Discount         float64  `json:"discount,omitempty"`
	Created          DateTime `json:"created,omitzero"`
}

type MaterialReportUpdateInput struct {
	ArticleNumber    string  `json:"articleNumber,omitempty"`
	Units            float64 `json:"units,omitempty"`
	InvoiceableUnits float64 `json:"invoicableUnits,omitempty"`
	Title            string  `json:"title,omitempty"`
	UnitText         string  `json:"unitText,omitempty"`
	Cost             float64 `json:"cost,omitempty"`
	UnitPrice        float64 `json:"unitPrice,omitempty"`
	Discount         float64 `json:"discount,omitempty"`
}

func (s *MaterialReportsService) List(ctx context.Context, params MaterialReportListParams) (*Page[MaterialReport], error) {
	return listPage[MaterialReport](ctx, s.client, materialReportsPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *MaterialReportsService) ListAll(ctx context.Context, params MaterialReportListParams) iter.Seq2[MaterialReport, error] {
	return listAll[MaterialReport](ctx, s.client, materialReportsPath, params)
}

func (s *MaterialReportsService) Get(ctx context.Context, id int) (*MaterialReport, error) {
	return getResource[MaterialReport](ctx, s.client, materialReportsPath, strconv.Itoa(id))
}

// Create returns the ID of the new material report.
func (s *MaterialReportsService) Create(ctx context.Context, input MaterialReportCreateInput) (int, error) {
	return createResourceID(ctx, s.client, input, materialReportsPath)
}

// Update changes only the fields set in input; zero values are not sent.
func (s *MaterialReportsService) Update(ctx context.Context, id int, input MaterialReportUpdateInput) error {
	return updateResource(ctx, s.client, input, materialReportsPath, strconv.Itoa(id))
}

func (s *MaterialReportsService) Delete(ctx context.Context, id int) error {
	return deleteResource(ctx, s.client, materialReportsPath, strconv.Itoa(id))
}
