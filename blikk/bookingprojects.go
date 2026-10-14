package blikk

import (
	"context"
	"errors"
	"iter"
	"net/url"
	"strconv"
)

const bookingProjectsPath = "v1/admin/bookingprojects"

type BookingProjectsService struct {
	client *Client
}

const (
	BookingProjectSortByName          = "name"
	BookingProjectSortByBookingNumber = "bookingNumber"
	BookingProjectSortByIsActive      = "isActive"
	BookingProjectSortByStartDate     = "startDate"
	BookingProjectSortByEndDate       = "endDate"
	BookingProjectSortByCreatedDate   = "createdDate"
	BookingProjectSortByUpdatedDate   = "updatedDate"
)

type BookingProjectListParams struct {
	ListOptions
	CreatedDateFrom DateTime `query:"filter.createdDateFrom"`
	CreatedDateTo   DateTime `query:"filter.createdDateTo"`
	UpdatedDateFrom DateTime `query:"filter.updatedDateFrom"`
	UpdatedDateTo   DateTime `query:"filter.updatedDateTo"`
	StartDateFrom   DateTime `query:"filter.startDateFrom"`
	StartDateTo     DateTime `query:"filter.startDateTo"`
	EndDateFrom     DateTime `query:"filter.endDateFrom"`
	EndDateTo       DateTime `query:"filter.endDateTo"`
	IsActive        *bool    `query:"filter.isActive"`
	Query           string   `query:"filter.query"`
}

type BookingProject struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	BookingNumber string    `json:"bookingNumber"`
	IsActive      bool      `json:"isActive"`
	StartDate     *DateTime `json:"startDate"`
	EndDate       *DateTime `json:"endDate"`
	CreatedBy     NameID    `json:"createdBy"`
	UpdatedBy     *NameID   `json:"updatedBy"`
	Created       DateTime  `json:"created"`
	Updated       *DateTime `json:"updated"`
}

// BookingProjectInput.IsActive is a pointer because the API defaults it to true on create and keeps it on update when omitted.
type BookingProjectInput struct {
	Name          string   `json:"name"`
	BookingNumber string   `json:"bookingNumber"`
	IsActive      *bool    `json:"isActive,omitempty"`
	StartDate     DateTime `json:"startDate,omitzero"`
	EndDate       DateTime `json:"endDate,omitzero"`
}

func (s *BookingProjectsService) List(ctx context.Context, params BookingProjectListParams) (*Page[BookingProject], error) {
	return listPage[BookingProject](ctx, s.client, bookingProjectsPath, params)
}

func (s *BookingProjectsService) ListAll(ctx context.Context, params BookingProjectListParams) iter.Seq2[BookingProject, error] {
	return listAll[BookingProject](ctx, s.client, bookingProjectsPath, params)
}

func (s *BookingProjectsService) Get(ctx context.Context, id int) (*BookingProject, error) {
	return getResource[BookingProject](ctx, s.client, bookingProjectsPath, strconv.Itoa(id))
}

func (s *BookingProjectsService) GetByBookingNumber(ctx context.Context, bookingNumber string) (*BookingProject, error) {
	// Without a booking number the same URL is the list endpoint.
	if bookingNumber == "" {
		return nil, errors.New("blikk: booking number must not be empty")
	}
	return getResourceWithQuery[BookingProject](ctx, s.client, bookingProjectsPath, url.Values{"bookingnumber": {bookingNumber}})
}

func (s *BookingProjectsService) Create(ctx context.Context, input BookingProjectInput) (*BookingProject, error) {
	return createResource[BookingProject](ctx, s.client, input, bookingProjectsPath)
}

func (s *BookingProjectsService) Update(ctx context.Context, id int, input BookingProjectInput) (*BookingProject, error) {
	return updateResourceReturning[BookingProject](ctx, s.client, input, bookingProjectsPath, strconv.Itoa(id))
}
