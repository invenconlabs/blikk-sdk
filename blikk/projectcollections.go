package blikk

import (
	"context"
	"iter"
	"strconv"
)

const projectCollectionsPath = "v1/Core/ProjectCollections"

type ProjectCollectionsService struct {
	client *Client
}

type ProjectCollectionType string

const (
	ProjectCollectionTypeNormal   ProjectCollectionType = "Normal"
	ProjectCollectionTypeAbsence  ProjectCollectionType = "Absence"
	ProjectCollectionTypeInternal ProjectCollectionType = "Internal"
)

func (t *ProjectCollectionType) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, t, ProjectCollectionTypeNormal, ProjectCollectionTypeAbsence, ProjectCollectionTypeInternal)
}

const (
	ProjectCollectionSortByName          = "name"
	ProjectCollectionSortByProjectNumber = "projectNumber"
)

type ProjectCollectionListParams struct {
	ListOptions
}

type ProjectCollection struct {
	ID            int                        `json:"id"`
	Name          string                     `json:"name"`
	Title         string                     `json:"title"`
	ProjectNumber string                     `json:"projectNumber"`
	Description   string                     `json:"description"`
	Status        *ProjectStatus             `json:"status"`
	Type          ProjectCollectionType      `json:"type"`
	Category      *ProjectCollectionCategory `json:"category"`
	Manager       *NameID                    `json:"manager"`
	Customer      *NameID                    `json:"customer"`
	FinishedDate  *DateTime                  `json:"finishedDate"`
}

type ProjectCollectionCategory struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type projectCollectionNextNumber struct {
	NextProjectNumberInCollection string `json:"nextProjectNumberInCollection"`
}

func (s *ProjectCollectionsService) List(ctx context.Context, params ProjectCollectionListParams) (*Page[ProjectCollection], error) {
	return listPage[ProjectCollection](ctx, s.client, projectCollectionsPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *ProjectCollectionsService) ListAll(ctx context.Context, params ProjectCollectionListParams) iter.Seq2[ProjectCollection, error] {
	return listAll[ProjectCollection](ctx, s.client, projectCollectionsPath, params)
}

func (s *ProjectCollectionsService) Get(ctx context.Context, id int) (*ProjectCollection, error) {
	return getResource[ProjectCollection](ctx, s.client, projectCollectionsPath, strconv.Itoa(id))
}

func (s *ProjectCollectionsService) NextProjectNumber(ctx context.Context, id int) (string, error) {
	next, err := getResource[projectCollectionNextNumber](ctx, s.client, projectCollectionsPath, strconv.Itoa(id), "nextProjectNumber")
	if err != nil {
		return "", err
	}
	return next.NextProjectNumberInCollection, nil
}
