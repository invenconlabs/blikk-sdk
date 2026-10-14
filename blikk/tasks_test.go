package blikk

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTasksListEncodesEveryFilter(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, "GET", "/v1/Core/Tasks", "tasks_list.json"))
	moment := DateTimeOf(time.Date(2026, 10, 1, 14, 0, 0, 0, time.UTC))

	_, err := client.Tasks.List(context.Background(), TaskListParams{
		ListOptions:           ListOptions{SortBy: TaskSortByStartDate, SortOrder: SortAscending},
		CreatedDateFrom:       moment,
		CreatedDateTo:         moment,
		UpdatedDateFrom:       moment,
		UpdatedDateTo:         moment,
		CreatedOrUpdatedSince: moment,
		From:                  DateTimeOf(time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)),
		To:                    DateTimeOf(time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC)),
		MemberIDs:             []int{1, 2},
		CreatedByUserIDs:      []int{3},
		TaskIDs:               []int{4, 5},
		TaskBoardID:           6,
		TaskBoardColumnID:     7,
		ProjectID:             8,
	})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"page":                         {"1"},
		"sortBy":                       {"startDate"},
		"sortOrder":                    {"ascending"},
		"filter.createdDateFrom":       {"2026-10-01 14:00:00"},
		"filter.createdDateTo":         {"2026-10-01 14:00:00"},
		"filter.updatedDateFrom":       {"2026-10-01 14:00:00"},
		"filter.updatedDateTo":         {"2026-10-01 14:00:00"},
		"filter.createdOrUpdatedSince": {"2026-10-01 14:00:00"},
		"filter.from":                  {"2026-10-01 00:00:00"},
		"filter.to":                    {"2026-10-31 00:00:00"},
		"filter.memberIds":             {"1,2"},
		"filter.createdByUserIds":      {"3"},
		"filter.taskIds":               {"4,5"},
		"filter.taskBoardId":           {"6"},
		"filter.taskBoardColumnId":     {"7"},
		"filter.projectId":             {"8"},
	}, handler.last().Query)
}

func TestTasksListDecodesFixture(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, "GET", "/v1/Core/Tasks", "tasks_list.json"))

	page, err := client.Tasks.List(context.Background(), TaskListParams{})

	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	task := page.Items[0]
	assert.Equal(t, 1, task.ID)
	assert.Equal(t, "Uppgift 1", task.Title)
	assert.Equal(t, "En beskrivning om vad som denna uppgift handlar om.", task.Description)
	assert.Nil(t, task.CompletedDate)
	assert.True(t, task.TimeReportingEnabled)
	assert.Equal(t, 10.5, task.Estimate)
	assert.Nil(t, task.StartDate)
	assert.Empty(t, task.StartTime)
	assert.Nil(t, task.EndDate)
	assert.Empty(t, task.EndTime)
	assert.Equal(t, 5, task.NumberOfComments)
	assert.Equal(t, 3, task.NumberOfAttachments)
	assert.Equal(t, 14, task.NumberOfChecklistPoints)
	assert.Equal(t, []TagRef{{ID: 1, Title: "2.8", Color: "#41b18b"}, {ID: 2, Title: "GUI", Color: "#41b18b"}}, task.Tags)
	assert.Equal(t, []TaskMember{{UserID: 1, UserName: "John Doe"}, {UserID: 2, UserName: "Jane Doe"}}, task.Members)
	assert.Equal(t, &NameID{ID: 1, Name: "Att göra"}, task.BoardColumn)
	assert.Equal(t, &NameID{ID: 1, Name: "Sprint 2.8"}, task.Board)
	assert.Equal(t, &ProjectRef{ID: 15, Name: "Utvecklingsprojektet", Number: "26"}, task.Project)
	assert.Nil(t, task.ProjectCollection)
	assert.Nil(t, task.Opportunity)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, task.CreatedByUser)
	assert.Equal(t, "2026-10-10 09:40:57", task.CreatedDate.String())
	assert.Nil(t, task.UpdatedDate)
	assert.Nil(t, task.ArchivedDate)
}

func TestTasksCreateSendsFullInput(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"id":42}`))

	id, err := client.Tasks.Create(context.Background(), TaskInput{
		Title:                "My task",
		Description:          "A very descriptive text about this task",
		IsCompleted:          true,
		TimeReportingEnabled: true,
		Estimate:             26.5,
		StartDate:            DateTimeOf(time.Date(2026, 10, 23, 9, 40, 57, 0, time.UTC)),
		EndDate:              DateTimeOf(time.Date(2026, 10, 24, 17, 0, 0, 0, time.UTC)),
		TagIDs:               []int{1, 5, 12},
		MemberIDs:            []int{1, 2},
		TaskBoardColumnID:    1,
		ProjectID:            15,
		IsArchived:           true,
	})

	require.NoError(t, err)
	assert.Equal(t, 42, id)
	request := handler.last()
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/v1/Core/Tasks", request.Path)
	assert.Equal(t, map[string]any{
		"title":                "My task",
		"description":          "A very descriptive text about this task",
		"isCompleted":          true,
		"timeReportingEnabled": true,
		"estimate":             26.5,
		"startDate":            "2026-10-23 09:40:57",
		"endDate":              "2026-10-24 17:00:00",
		"tagIds":               []any{1.0, 5.0, 12.0},
		"memberIds":            []any{1.0, 2.0},
		"taskBoardColumnId":    1.0,
		"projectId":            15.0,
		"isArchived":           true,
	}, request.jsonBody(t))
}

func TestTasksCreateSendsMinimalInput(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"id":42}`))

	_, err := client.Tasks.Create(context.Background(), TaskInput{Title: "My task"})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"title": "My task", "isCompleted": false, "timeReportingEnabled": false, "isArchived": false}, handler.last().jsonBody(t))
}

func TestTasksUpdateSendsPut(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	err := client.Tasks.Update(context.Background(), 7, TaskInput{Title: "Renamed", ProjectID: 15})

	require.NoError(t, err)
	request := handler.last()
	assert.Equal(t, http.MethodPut, request.Method)
	assert.Equal(t, "/v1/Core/Tasks/7", request.Path)
	assert.Equal(t, map[string]any{"title": "Renamed", "projectId": 15.0, "isCompleted": false, "timeReportingEnabled": false, "isArchived": false}, request.jsonBody(t))
}

func TestTasksDeleteSendsDelete(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	require.NoError(t, client.Tasks.Delete(context.Background(), 7))

	assertTaskRequest(t, handler.last(), http.MethodDelete, "/v1/Core/Tasks/7")
}

func TestTasksListCommentsDecodesFixture(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, "GET", "/v1/Core/Tasks/1/Comments", "task_comments_list.json"))

	page, err := client.Tasks.ListComments(context.Background(), 1, ListOptions{SortBy: TaskCommentSortByCreatedDate, SortOrder: SortDescending})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{"page": {"1"}, "sortBy": {"createdDate"}, "sortOrder": {"descending"}}, handler.last().Query)
	require.Len(t, page.Items, 1)
	comment := page.Items[0]
	assert.Equal(t, 1, comment.ID)
	assert.Equal(t, "Lorem ipsum dolar sit amet", comment.Text)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, comment.CreatedByUser)
	assert.Nil(t, comment.UpdatedByUser)
	assert.Equal(t, "2026-10-11 09:40:57", comment.CreatedDate.String())
	assert.Nil(t, comment.UpdatedDate)
}

func TestTasksListAllCommentsWalksTaskPath(t *testing.T) {
	client := newTestClient(t, serveFixtureAt(t, "GET", "/v1/Core/Tasks/1/Comments", "task_comments_list.json"))

	var ids []int
	for comment, err := range client.Tasks.ListAllComments(context.Background(), 1, ListOptions{}) {
		require.NoError(t, err)
		ids = append(ids, comment.ID)
	}

	assert.Equal(t, []int{1}, ids)
}

func TestTasksCommentWritesSendText(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"id":9}`))
	ctx := context.Background()

	id, err := client.Tasks.CreateComment(ctx, 1, "Lorem ipsum dolar sit amet.")
	require.NoError(t, err)
	assert.Equal(t, 9, id)
	assertTaskRequest(t, handler.last(), http.MethodPost, "/v1/Core/Tasks/1/Comments")
	assert.Equal(t, map[string]any{"text": "Lorem ipsum dolar sit amet."}, handler.last().jsonBody(t))

	require.NoError(t, client.Tasks.UpdateComment(ctx, 1, 9, "Edited"))
	assertTaskRequest(t, handler.last(), http.MethodPut, "/v1/Core/Tasks/1/Comments/9")
	assert.Equal(t, map[string]any{"text": "Edited"}, handler.last().jsonBody(t))

	require.NoError(t, client.Tasks.DeleteComment(ctx, 1, 9))
	assertTaskRequest(t, handler.last(), http.MethodDelete, "/v1/Core/Tasks/1/Comments/9")
}

func TestTasksListChecklistsDecodesFixture(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, "GET", "/v1/Core/Tasks/1/Checklists", "task_checklists_list.json"))

	page, err := client.Tasks.ListChecklists(context.Background(), 1, ListOptions{SortBy: TaskChecklistSortByName})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{"page": {"1"}, "sortBy": {"name"}}, handler.last().Query)
	require.Len(t, page.Items, 1)
	checklist := page.Items[0]
	assert.Equal(t, 1, checklist.ID)
	assert.Equal(t, "Att göra", checklist.Name)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, checklist.CreatedByUser)
	assert.Equal(t, "2026-10-11 09:40:57", checklist.CreatedDate.String())
	assert.Nil(t, checklist.UpdatedDate)
	require.Len(t, checklist.Items, 2)
	item := checklist.Items[1]
	assert.Equal(t, 2, item.ID)
	assert.Equal(t, "Lorem ipsum 1", item.Title)
	assert.Equal(t, 1, item.SortOrder)
	assert.Equal(t, "2026-10-13 09:40:57", item.CompletedDate.String())
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, item.CreatedByUser)
	assert.Equal(t, "2026-10-11 09:40:57", item.CreatedDate.String())
	assert.Nil(t, item.UpdatedDate)
}

func TestTasksCreateChecklistSendsFullInput(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"id":3}`))

	id, err := client.Tasks.CreateChecklist(context.Background(), 1, TaskChecklistInput{
		Name: "Lorem ipsum",
		Items: []TaskChecklistItemInput{
			{Name: "Lorem ipsum dolar", IsCompleted: true},
			{Name: "Lorem ipsum dolar sit"},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, 3, id)
	assertTaskRequest(t, handler.last(), http.MethodPost, "/v1/Core/Tasks/1/Checklists")
	assert.Equal(t, map[string]any{
		"name": "Lorem ipsum",
		"items": []any{
			map[string]any{"name": "Lorem ipsum dolar", "isCompleted": true},
			map[string]any{"name": "Lorem ipsum dolar sit", "isCompleted": false},
		},
	}, handler.last().jsonBody(t))
}

func TestTasksCreateChecklistSendsMinimalInput(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"id":3}`))

	_, err := client.Tasks.CreateChecklist(context.Background(), 1, TaskChecklistInput{Name: "Empty"})

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"name": "Empty"}, handler.last().jsonBody(t))
}

func TestTasksUpdateChecklistKeepsExistingItemIDs(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	err := client.Tasks.UpdateChecklist(context.Background(), 1, 3, TaskChecklistInput{
		Name: "Lorem ipsum dolar",
		Items: []TaskChecklistItemInput{
			{ID: 1, Name: "Existing item"},
			{Name: "New item"},
		},
	})

	require.NoError(t, err)
	assertTaskRequest(t, handler.last(), http.MethodPut, "/v1/Core/Tasks/1/Checklists/3")
	assert.Equal(t, map[string]any{
		"name": "Lorem ipsum dolar",
		"items": []any{
			map[string]any{"id": 1.0, "name": "Existing item", "isCompleted": false},
			map[string]any{"name": "New item", "isCompleted": false},
		},
	}, handler.last().jsonBody(t))
}

func TestTasksDeleteChecklistSendsDelete(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))

	require.NoError(t, client.Tasks.DeleteChecklist(context.Background(), 1, 3))

	assertTaskRequest(t, handler.last(), http.MethodDelete, "/v1/Core/Tasks/1/Checklists/3")
}

func TestTasksListAttachmentsDecodesFileAndLink(t *testing.T) {
	client, handler := newRecordingClient(t, serveFixtureAt(t, "GET", "/v1/Core/Tasks/1/Attachments", "task_attachments_list.json"))

	page, err := client.Tasks.ListAttachments(context.Background(), 1, ListOptions{SortBy: TaskAttachmentSortByName})

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{"page": {"1"}, "sortBy": {"name"}}, handler.last().Query)
	require.Len(t, page.Items, 2)
	file, link := page.Items[0], page.Items[1]
	assert.Equal(t, 1, file.ID)
	assert.Equal(t, "MyImage.png", file.Name)
	assert.Equal(t, TaskAttachmentTypeFile, file.Type)
	assert.Equal(t, "image/png", file.ContentType)
	assert.Equal(t, Ptr(12542), file.SizeInBytes)
	assert.Empty(t, file.ExternalLinkURL)
	assert.Equal(t, &NameID{ID: 1, Name: "John Doe"}, file.CreatedByUser)
	assert.Equal(t, "2026-10-11 09:40:57", file.CreatedDate.String())
	assert.Nil(t, file.UpdatedDate)
	assert.Equal(t, TaskAttachmentTypeLink, link.Type)
	assert.Empty(t, link.ContentType)
	assert.Nil(t, link.SizeInBytes)
	assert.Equal(t, "https://example.com", link.ExternalLinkURL)
}

func TestTasksDownloadAttachmentStreamsFile(t *testing.T) {
	content := []byte{0x89, 'P', 'N', 'G'}
	client := newTestClient(t, routes(t, map[string]http.HandlerFunc{
		"GET /v1/Core/Tasks/1/Attachments/2/Download": serveFile("image/png", `attachment; filename="MyImage.png"`, content),
	}))

	download, err := client.Tasks.DownloadAttachment(context.Background(), 1, 2)

	require.NoError(t, err)
	defer download.Content.Close()
	body, err := io.ReadAll(download.Content)
	require.NoError(t, err)
	assert.Equal(t, content, body)
	assert.Equal(t, "image/png", download.ContentType)
	assert.Equal(t, "MyImage.png", download.Filename)
}

func TestTasksCreateAttachmentSendsLinkAsMultipartFields(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"id":5}`))

	id, err := client.Tasks.CreateAttachment(context.Background(), 1, TaskAttachmentInput{
		Name:            "Our external library",
		ExternalURL:     "https://example.com",
		ExternalURLType: ExternalURLTypeOneDrive,
	})

	require.NoError(t, err)
	assert.Equal(t, 5, id)
	request := handler.last()
	assertTaskRequest(t, request, http.MethodPost, "/v1/Core/Tasks/1/Attachments")
	assert.True(t, strings.HasPrefix(request.Header.Get("Content-Type"), "multipart/form-data; boundary="))
	form := request.multipartForm(t)
	assert.Equal(t, map[string]string{
		"name":            "Our external library",
		"externalUrl":     "https://example.com",
		"externalUrlType": "OneDrive",
	}, formValues(form))
	assert.Empty(t, form.File)
}

func TestTasksCreateAttachmentLinkOmitsEmptyFields(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"id":5}`))

	_, err := client.Tasks.CreateAttachment(context.Background(), 1, TaskAttachmentInput{ExternalURL: "https://example.com"})

	require.NoError(t, err)
	assert.Equal(t, map[string]string{"externalUrl": "https://example.com"}, formValues(handler.last().multipartForm(t)))
}

func TestTasksCreateAttachmentUploadsFileWithName(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, `{"id":6}`))

	id, err := client.Tasks.CreateAttachment(context.Background(), 1, TaskAttachmentInput{
		Name:        "My image",
		ExternalURL: "https://example.com/ignored",
		File:        &UploadFile{Filename: "myimage.png", ContentType: "image/png", Content: strings.NewReader("\x89PNG bytes")},
	})

	require.NoError(t, err)
	assert.Equal(t, 6, id)
	form := handler.last().multipartForm(t)
	assert.Equal(t, map[string]string{"name": "My image"}, formValues(form))
	assert.Equal(t, []uploadedPart{{Filename: "myimage.png", ContentType: "image/png", Content: "\x89PNG bytes"}}, uploadedParts(t, form, "file"))
	assert.Len(t, form.File, 1)
}

func TestTasksCreateAttachmentAcceptsEmptyResponse(t *testing.T) {
	client, _ := newRecordingClient(t, respondWith(http.StatusOK, ""))

	id, err := client.Tasks.CreateAttachment(context.Background(), 1, TaskAttachmentInput{
		File: &UploadFile{Filename: "myimage.png", Content: strings.NewReader("png")},
	})

	require.NoError(t, err)
	assert.Zero(t, id)
}

func TestTasksAttachmentWritesSendName(t *testing.T) {
	client, handler := newRecordingClient(t, respondWith(http.StatusOK, ""))
	ctx := context.Background()

	require.NoError(t, client.Tasks.UpdateAttachment(ctx, 1, 5, "My updated attachment"))
	assertTaskRequest(t, handler.last(), http.MethodPut, "/v1/Core/Tasks/1/Attachments/5")
	assert.Equal(t, map[string]any{"name": "My updated attachment"}, handler.last().jsonBody(t))

	require.NoError(t, client.Tasks.DeleteAttachment(ctx, 1, 5))
	assertTaskRequest(t, handler.last(), http.MethodDelete, "/v1/Core/Tasks/1/Attachments/5")
}

func TestExternalURLTypeDecodesDocumentedLowercase(t *testing.T) {
	var decoded struct {
		Type ExternalURLType `json:"type"`
	}

	require.NoError(t, json.Unmarshal([]byte(`{"type":"googleDrive"}`), &decoded))

	assert.Equal(t, ExternalURLTypeGoogleDrive, decoded.Type)
}

func assertTaskRequest(t *testing.T, request recordedRequest, method, path string) {
	t.Helper()
	assert.Equal(t, method, request.Method)
	assert.Equal(t, path, request.Path)
}
