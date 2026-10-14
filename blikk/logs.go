package blikk

import (
	"context"
	"iter"
)

const logsPath = "v1/Core/Logs"

type LogsService struct {
	client *Client
}

type LogType string

const (
	LogTypeProjectAll              LogType = "Project_All"
	LogTypeProjectCreated          LogType = "Project_Created"
	LogTypeProjectUpdated          LogType = "Project_Updated"
	LogTypeProjectAddedNote        LogType = "Project_Added_Note"
	LogTypeProjectDeleted          LogType = "Project_Deleted"
	LogTypeProjectMentioned        LogType = "Project_Mentioned"
	LogTypeProjectChangedStatus    LogType = "Project_Changed_Status"
	LogTypeProjectUpdatedStartDate LogType = "Project_Updated_StartDate"
	LogTypeProjectUpdatedEndDate   LogType = "Project_Updated_EndDate"
	LogTypeProjectRemovedStartDate LogType = "Project_Removed_StartDate"
	LogTypeProjectRemovedEndDate   LogType = "Project_Removed_EndDate"
)

func (t *LogType) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, t,
		LogTypeProjectAll, LogTypeProjectCreated, LogTypeProjectUpdated, LogTypeProjectAddedNote,
		LogTypeProjectDeleted, LogTypeProjectMentioned, LogTypeProjectChangedStatus, LogTypeProjectUpdatedStartDate,
		LogTypeProjectUpdatedEndDate, LogTypeProjectRemovedStartDate, LogTypeProjectRemovedEndDate)
}

const LogSortByCreatedDate = "createdDate"

type LogListParams struct {
	ListOptions
	LogType         LogType  `query:"filter.logType"`
	CreatedDateFrom DateTime `query:"filter.createdDateFrom"`
	CreatedDateTo   DateTime `query:"filter.createdDateTo"`
	ProjectID       int      `query:"filter.projectId"`
}

type Log struct {
	ID          int       `json:"id"`
	Text        string    `json:"text"`
	Note        string    `json:"note"`
	LogType     LogType   `json:"logType"`
	Entity      LogEntity `json:"entity"`
	CreatedBy   *NameID   `json:"createdBy"`
	CreatedDate DateTime  `json:"createdDate"`
}

type LogEntity struct {
	EntityType EntityType `json:"entityType"`
	ID         int        `json:"id"`
	Name       string     `json:"name"`
}

func (s *LogsService) List(ctx context.Context, params LogListParams) (*Page[Log], error) {
	return listPage[Log](ctx, s.client, logsPath, params)
}

// ListAll walks every page from params.Page onwards; on failure it yields a single error and stops.
func (s *LogsService) ListAll(ctx context.Context, params LogListParams) iter.Seq2[Log, error] {
	return listAll[Log](ctx, s.client, logsPath, params)
}
