package blikk

import (
	"context"
	"iter"
	"strconv"
)

const tasksPath = "v1/Core/Tasks"

type TasksService struct {
	client *Client
}

const (
	TaskSortByCreatedDate = "createdDate"
	TaskSortByTitle       = "title"
	TaskSortByStartDate   = "startDate"
	TaskSortByEndDate     = "endDate"
)

type TaskListParams struct {
	ListOptions
	CreatedDateFrom       DateTime `query:"filter.createdDateFrom"`
	CreatedDateTo         DateTime `query:"filter.createdDateTo"`
	UpdatedDateFrom       DateTime `query:"filter.updatedDateFrom"`
	UpdatedDateTo         DateTime `query:"filter.updatedDateTo"`
	CreatedOrUpdatedSince DateTime `query:"filter.createdOrUpdatedSince"`
	From                  DateTime `query:"filter.from"`
	To                    DateTime `query:"filter.to"`
	MemberIDs             []int    `query:"filter.memberIds"`
	CreatedByUserIDs      []int    `query:"filter.createdByUserIds"`
	TaskIDs               []int    `query:"filter.taskIds"`
	TaskBoardID           int      `query:"filter.taskBoardId"`
	TaskBoardColumnID     int      `query:"filter.taskBoardColumnId"`
	ProjectID             int      `query:"filter.projectId"`
}

type TaskListItem struct {
	ID                      int          `json:"id"`
	Title                   string       `json:"title"`
	Description             string       `json:"description"`
	CompletedDate           *DateTime    `json:"completedDate"`
	TimeReportingEnabled    bool         `json:"timeReportingEnabled"`
	Estimate                float64      `json:"estimate"`
	StartDate               *Date        `json:"startDate"`
	StartTime               string       `json:"startTime"`
	EndDate                 *Date        `json:"endDate"`
	EndTime                 string       `json:"endTime"`
	NumberOfComments        int          `json:"numberOfComments"`
	NumberOfAttachments     int          `json:"numberOfAttachments"`
	NumberOfChecklistPoints int          `json:"numberOfChecklistPoints"`
	Tags                    []Tag        `json:"tags"`
	Members                 []TaskMember `json:"members"`
	BoardColumn             *NameID      `json:"boardColumn"`
	Board                   *NameID      `json:"board"`
	Project                 *ProjectRef  `json:"project"`
	ProjectCollection       *NameID      `json:"projectCollection"`
	Opportunity             *NameID      `json:"opportunity"`
	CreatedByUser           *NameID      `json:"createdByUser"`
	CreatedDate             *DateTime    `json:"createdDate"`
	UpdatedDate             *DateTime    `json:"updatedDate"`
	ArchivedDate            *DateTime    `json:"archivedDate"`
}

type TaskMember struct {
	UserID   int    `json:"userId"`
	UserName string `json:"userName"`
}

type TaskInput struct {
	Title                string   `json:"title"`
	Description          string   `json:"description,omitempty"`
	IsCompleted          bool     `json:"isCompleted"`
	TimeReportingEnabled bool     `json:"timeReportingEnabled"`
	Estimate             float64  `json:"estimate,omitempty"`
	StartDate            DateTime `json:"startDate,omitzero"`
	EndDate              DateTime `json:"endDate,omitzero"`
	TagIDs               []int    `json:"tagIds,omitempty"`
	MemberIDs            []int    `json:"memberIds,omitempty"`
	TaskBoardColumnID    int      `json:"taskBoardColumnId,omitempty"`
	ProjectID            int      `json:"projectId,omitempty"`
	IsArchived           bool     `json:"isArchived"`
}

const TaskCommentSortByCreatedDate = "createdDate"

type taskCommentInput struct {
	Text string `json:"text"`
}

const TaskChecklistSortByName = "name"

type TaskChecklist struct {
	ID            int                 `json:"id"`
	Name          string              `json:"name"`
	Items         []TaskChecklistItem `json:"items"`
	CreatedByUser *NameID             `json:"createdByUser"`
	CreatedDate   *DateTime           `json:"createdDate"`
	UpdatedDate   *DateTime           `json:"updatedDate"`
}

type TaskChecklistItem struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	SortOrder     int       `json:"sortOrder"`
	CompletedDate *DateTime `json:"completedDate"`
	CreatedByUser *NameID   `json:"createdByUser"`
	CreatedDate   *DateTime `json:"createdDate"`
	UpdatedDate   *DateTime `json:"updatedDate"`
}

type TaskChecklistInput struct {
	Name  string                   `json:"name"`
	Items []TaskChecklistItemInput `json:"items,omitempty"`
}

// On UpdateChecklist an item with an ID keeps that existing item and one without is created.
type TaskChecklistItemInput struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	IsCompleted bool   `json:"isCompleted"`
}

const TaskAttachmentSortByName = "name"

type TaskAttachmentType string

const (
	TaskAttachmentTypeFile TaskAttachmentType = "File"
	TaskAttachmentTypeLink TaskAttachmentType = "Link"
)

func (t *TaskAttachmentType) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, t, TaskAttachmentTypeFile, TaskAttachmentTypeLink)
}

type ExternalURLType string

const (
	ExternalURLTypeGoogleDrive ExternalURLType = "GoogleDrive"
	ExternalURLTypeDropbox     ExternalURLType = "Dropbox"
	ExternalURLTypeOneDrive    ExternalURLType = "OneDrive"
)

func (t *ExternalURLType) UnmarshalJSON(data []byte) error {
	return unmarshalEnum(data, t, ExternalURLTypeGoogleDrive, ExternalURLTypeDropbox, ExternalURLTypeOneDrive)
}

type TaskAttachment struct {
	ID              int                `json:"id"`
	Name            string             `json:"name"`
	Type            TaskAttachmentType `json:"type"`
	ContentType     string             `json:"contentType"`
	SizeInBytes     *int               `json:"sizeInBytes"`
	ExternalLinkURL string             `json:"externalLinkUrl"`
	CreatedByUser   *NameID            `json:"createdByUser"`
	CreatedDate     *DateTime          `json:"createdDate"`
	UpdatedDate     *DateTime          `json:"updatedDate"`
}

type TaskAttachmentInput struct {
	Name            string          `json:"name,omitempty"`
	ExternalURL     string          `json:"externalUrl,omitempty"`
	ExternalURLType ExternalURLType `json:"externalUrlType,omitempty"`
}

type taskAttachmentNameInput struct {
	Name string `json:"name"`
}

func (s *TasksService) List(ctx context.Context, params TaskListParams) (*Page[TaskListItem], error) {
	return listPage[TaskListItem](ctx, s.client, tasksPath, params)
}

func (s *TasksService) ListAll(ctx context.Context, params TaskListParams) iter.Seq2[TaskListItem, error] {
	return listAll[TaskListItem](ctx, s.client, tasksPath, params)
}

func (s *TasksService) Create(ctx context.Context, input TaskInput) (int, error) {
	return createResourceID(ctx, s.client, input, tasksPath)
}

// Update replaces the whole task; any field left unset in input is cleared.
func (s *TasksService) Update(ctx context.Context, taskID int, input TaskInput) error {
	return updateResource(ctx, s.client, input, tasksPath, strconv.Itoa(taskID))
}

func (s *TasksService) Delete(ctx context.Context, id int) error {
	return deleteResource(ctx, s.client, tasksPath, strconv.Itoa(id))
}

func (s *TasksService) ListComments(ctx context.Context, taskID int, params ListOptions) (*Page[Comment], error) {
	return listTaskSubresourcePage[Comment](ctx, s.client, taskID, "Comments", params)
}

func (s *TasksService) ListAllComments(ctx context.Context, taskID int, params ListOptions) iter.Seq2[Comment, error] {
	return listAllTaskSubresource[Comment](ctx, s.client, taskID, "Comments", params)
}

func (s *TasksService) CreateComment(ctx context.Context, taskID int, text string) (int, error) {
	return createResourceID(ctx, s.client, taskCommentInput{Text: text}, tasksPath, strconv.Itoa(taskID), "Comments")
}

func (s *TasksService) UpdateComment(ctx context.Context, taskID, commentID int, text string) error {
	return updateResource(ctx, s.client, taskCommentInput{Text: text}, tasksPath, strconv.Itoa(taskID), "Comments", strconv.Itoa(commentID))
}

func (s *TasksService) DeleteComment(ctx context.Context, taskID, commentID int) error {
	return deleteResource(ctx, s.client, tasksPath, strconv.Itoa(taskID), "Comments", strconv.Itoa(commentID))
}

func (s *TasksService) ListChecklists(ctx context.Context, taskID int, params ListOptions) (*Page[TaskChecklist], error) {
	return listTaskSubresourcePage[TaskChecklist](ctx, s.client, taskID, "Checklists", params)
}

func (s *TasksService) ListAllChecklists(ctx context.Context, taskID int, params ListOptions) iter.Seq2[TaskChecklist, error] {
	return listAllTaskSubresource[TaskChecklist](ctx, s.client, taskID, "Checklists", params)
}

func (s *TasksService) CreateChecklist(ctx context.Context, taskID int, input TaskChecklistInput) (int, error) {
	return createResourceID(ctx, s.client, input, tasksPath, strconv.Itoa(taskID), "Checklists")
}

// UpdateChecklist replaces the checklist and its items.
func (s *TasksService) UpdateChecklist(ctx context.Context, taskID, checklistID int, input TaskChecklistInput) error {
	return updateResource(ctx, s.client, input, tasksPath, strconv.Itoa(taskID), "Checklists", strconv.Itoa(checklistID))
}

func (s *TasksService) DeleteChecklist(ctx context.Context, taskID, checklistID int) error {
	return deleteResource(ctx, s.client, tasksPath, strconv.Itoa(taskID), "Checklists", strconv.Itoa(checklistID))
}

func (s *TasksService) ListAttachments(ctx context.Context, taskID int, params ListOptions) (*Page[TaskAttachment], error) {
	return listTaskSubresourcePage[TaskAttachment](ctx, s.client, taskID, "Attachments", params)
}

func (s *TasksService) ListAllAttachments(ctx context.Context, taskID int, params ListOptions) iter.Seq2[TaskAttachment, error] {
	return listAllTaskSubresource[TaskAttachment](ctx, s.client, taskID, "Attachments", params)
}

// DownloadAttachment streams the file; the caller must close Content.
func (s *TasksService) DownloadAttachment(ctx context.Context, taskID, attachmentID int) (*Download, error) {
	return downloadTaskAttachment(ctx, s.client, taskID, attachmentID)
}

// CreateAttachment adds a link attachment; uploading a file needs a multipart request, which the client does not send.
func (s *TasksService) CreateAttachment(ctx context.Context, taskID int, input TaskAttachmentInput) (int, error) {
	return createResourceID(ctx, s.client, input, tasksPath, strconv.Itoa(taskID), "Attachments")
}

func (s *TasksService) UpdateAttachment(ctx context.Context, taskID, attachmentID int, name string) error {
	return updateResource(ctx, s.client, taskAttachmentNameInput{Name: name}, tasksPath, strconv.Itoa(taskID), "Attachments", strconv.Itoa(attachmentID))
}

func (s *TasksService) DeleteAttachment(ctx context.Context, taskID, attachmentID int) error {
	return deleteResource(ctx, s.client, tasksPath, strconv.Itoa(taskID), "Attachments", strconv.Itoa(attachmentID))
}

func listTaskSubresourcePage[T any](ctx context.Context, c *Client, taskID int, subresource string, params ListOptions) (*Page[T], error) {
	path, err := resourcePath(tasksPath, strconv.Itoa(taskID), subresource)
	if err != nil {
		return nil, err
	}
	return listPage[T](ctx, c, path, params)
}

func listAllTaskSubresource[T any](ctx context.Context, c *Client, taskID int, subresource string, params ListOptions) iter.Seq2[T, error] {
	path, err := resourcePath(tasksPath, strconv.Itoa(taskID), subresource)
	if err != nil {
		return failedSeq[T](err)
	}
	return listAll[T](ctx, c, path, params)
}

func downloadTaskAttachment(ctx context.Context, c *Client, taskID, attachmentID int) (*Download, error) {
	path, err := resourcePath(tasksPath, strconv.Itoa(taskID), "Attachments", strconv.Itoa(attachmentID), "Download")
	if err != nil {
		return nil, err
	}
	return c.download(ctx, path, nil)
}
