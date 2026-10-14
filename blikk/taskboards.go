package blikk

import (
	"context"
	"iter"
)

const taskBoardsPath = "v1/Core/TaskBoards"

type TaskBoardsService struct {
	client *Client
}

const (
	TaskBoardSortByCreatedDate = "createdDate"
	TaskBoardSortByUpdatedDate = "updatedDate"
	TaskBoardSortByName        = "name"
)

type TaskBoardListParams struct {
	ListOptions
	Query      string `query:"filter.query"`
	ProjectIDs []int  `query:"filter.projectIds"`
}

type TaskBoard struct {
	ID            int               `json:"id"`
	Name          string            `json:"name"`
	Project       *ProjectRef       `json:"project"`
	Tags          []TagRef          `json:"tags"`
	Columns       []TaskBoardColumn `json:"columns"`
	Members       []NameID          `json:"members"`
	CreatedByUser *NameID           `json:"createdByUser"`
	UpdatedByUser *NameID           `json:"updatedByUser"`
	CreatedDate   *DateTime         `json:"createdDate"`
	UpdatedDate   *DateTime         `json:"updatedDate"`
}

type TaskBoardColumn struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	TaskCount     int       `json:"taskCount"`
	CreatedByUser *NameID   `json:"createdByUser"`
	UpdatedByUser *NameID   `json:"updatedByUser"`
	CreatedDate   *DateTime `json:"createdDate"`
	UpdatedDate   *DateTime `json:"updatedDate"`
}

func (s *TaskBoardsService) List(ctx context.Context, params TaskBoardListParams) (*Page[TaskBoard], error) {
	return listPage[TaskBoard](ctx, s.client, taskBoardsPath, params)
}

func (s *TaskBoardsService) ListAll(ctx context.Context, params TaskBoardListParams) iter.Seq2[TaskBoard, error] {
	return listAll[TaskBoard](ctx, s.client, taskBoardsPath, params)
}
