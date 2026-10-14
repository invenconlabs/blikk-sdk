package blikk

import (
	"context"
	"iter"
)

const absenceProjectsPath = "v1/Admin/AbsenceProjects"

type AbsenceProjectsService struct {
	client *Client
}

const AbsenceProjectSortByName = "name"

type AbsenceProjectListParams struct {
	ListOptions
}

type AbsenceProjectListItem struct {
	ID                               int       `json:"id"`
	Name                             string    `json:"name"`
	CommentRequiredWhenTimeReporting bool      `json:"commentRequiredWhenTimeReporting"`
	SalaryCodeForMonthlySalary       string    `json:"salaryCodeForMonthlySalary"`
	SalaryCodeForHourlySalary        string    `json:"salaryCodeForHourlySalary"`
	RefersToSickness                 bool      `json:"refersToSickness"`
	IsActive                         bool      `json:"isActive"`
	CreatedBy                        *NameID   `json:"createdBy"`
	UpdatedBy                        *NameID   `json:"updatedBy"`
	Created                          *DateTime `json:"created"`
	Updated                          *DateTime `json:"updated"`
}

func (s *AbsenceProjectsService) List(ctx context.Context, params AbsenceProjectListParams) (*Page[AbsenceProjectListItem], error) {
	return listPage[AbsenceProjectListItem](ctx, s.client, absenceProjectsPath, params)
}

func (s *AbsenceProjectsService) ListAll(ctx context.Context, params AbsenceProjectListParams) iter.Seq2[AbsenceProjectListItem, error] {
	return listAll[AbsenceProjectListItem](ctx, s.client, absenceProjectsPath, params)
}
