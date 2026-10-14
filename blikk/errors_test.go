package blikk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getUserError(t *testing.T, handler http.HandlerFunc) *Error {
	t.Helper()
	client := newTestClient(t, handler)
	_, err := client.Users.Get(context.Background(), 1)
	require.Error(t, err)
	var apiErr *Error
	require.True(t, errors.As(err, &apiErr), "expected *Error, got %T", err)
	return apiErr
}

func TestBadRequestErrorDecodesModelState(t *testing.T) {
	body := `{"message":"The request is invalid.","modelState":{"page":["The field Page must be greater than or equal to 1."],"filter.From":["Invalid datespan."]}}`
	apiErr := getUserError(t, respondWith(http.StatusBadRequest, body))

	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	assert.Equal(t, "The request is invalid.", apiErr.Message)
	assert.Equal(t, map[string][]string{
		"page":        {"The field Page must be greater than or equal to 1."},
		"filter.From": {"Invalid datespan."},
	}, apiErr.ModelState)
	assert.Equal(t, []byte(body), apiErr.Body)
	assert.Equal(t,
		"blikk: GET /v1/Admin/Users/1: 400 Bad Request: The request is invalid. (filter.From: Invalid datespan.; page: The field Page must be greater than or equal to 1.)",
		apiErr.Error())
}

func TestForbiddenErrorDecodesPermissions(t *testing.T) {
	body := `{"message":"You don't have the necessary permissions to perform this action","ipAddressAllowed":true,"permissionsRequired":["Contact_Read"]}`
	apiErr := getUserError(t, respondWith(http.StatusForbidden, body))

	assert.Equal(t, http.StatusForbidden, apiErr.StatusCode)
	assert.Equal(t, []string{"Contact_Read"}, apiErr.PermissionsRequired)
	require.NotNil(t, apiErr.IPAddressAllowed)
	assert.True(t, *apiErr.IPAddressAllowed)
	assert.Equal(t, "blikk: GET /v1/Admin/Users/1: 403 Forbidden: You don't have the necessary permissions to perform this action (requires Contact_Read)", apiErr.Error())
}

func TestForbiddenErrorDecodesIPRestriction(t *testing.T) {
	apiErr := getUserError(t, respondWith(http.StatusForbidden, `{"message":"Your IP-address \"127.0.0.1\" is not allowed","ipAddressAllowed":false,"permissionsRequired":[]}`))

	require.NotNil(t, apiErr.IPAddressAllowed)
	assert.False(t, *apiErr.IPAddressAllowed)
}

func TestNotFoundErrorWithEmptyBody(t *testing.T) {
	client := newTestClient(t, respondWith(http.StatusNotFound, ""))
	_, err := client.Projects.Get(context.Background(), 999)

	require.Error(t, err)
	assert.True(t, IsNotFound(err))
	assert.False(t, IsUnauthorized(err))
	assert.Equal(t, "blikk: GET /v1/Core/Projects/999: 404 Not Found", err.Error())
}

func TestServerErrorDecodesRequestID(t *testing.T) {
	apiErr := getUserError(t, respondWith(http.StatusInternalServerError, `{"message":"An unexpected error has occurred.","requestId":"r-0fbf44e8-afb8-42be-ae69-55c8bc03311b"}`))

	assert.Equal(t, "r-0fbf44e8-afb8-42be-ae69-55c8bc03311b", apiErr.RequestID)
	assert.Equal(t, "blikk: GET /v1/Admin/Users/1: 500 Internal Server Error: An unexpected error has occurred. (request r-0fbf44e8-afb8-42be-ae69-55c8bc03311b)", apiErr.Error())
}

func TestNonJSONErrorBodyIsKeptRaw(t *testing.T) {
	apiErr := getUserError(t, respondWith(http.StatusBadGateway, "<html>bad gateway</html>"))

	assert.Equal(t, http.StatusBadGateway, apiErr.StatusCode)
	assert.Empty(t, apiErr.Message)
	assert.Equal(t, "<html>bad gateway</html>", string(apiErr.Body))
}

func TestErrorBodyIsCappedAt64KB(t *testing.T) {
	apiErr := getUserError(t, respondWith(http.StatusInternalServerError, strings.Repeat("x", 100<<10)))

	assert.Len(t, apiErr.Body, 64<<10)
}

func TestErrorMessageListsEveryRequiredPermissionBeforeRequestID(t *testing.T) {
	apiErr := &Error{StatusCode: http.StatusForbidden, Message: "Denied", PermissionsRequired: []string{"Timereport_Read", "Timereport_Write"}, RequestID: "r-abc"}

	assert.Equal(t, "blikk: 403 Forbidden: Denied (requires Timereport_Read, Timereport_Write) (request r-abc)", apiErr.Error())
}

func TestMissingUserBadRequestIsNotFound(t *testing.T) {
	apiErr := getUserError(t, respondWith(http.StatusBadRequest, `{"message":"The user could not be found."}`))

	assert.True(t, IsNotFound(apiErr))
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
}

func TestMissingTimeReportBadRequestIsNotFound(t *testing.T) {
	client := newTestClient(t, respondWith(http.StatusBadRequest, `{"message":"The request is invalid.","modelState":{"":["The timereport could not be found."]}}`))

	_, err := client.TimeReports.Get(context.Background(), 5)

	assert.True(t, IsNotFound(err))
	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	assert.Equal(t, http.MethodGet, apiErr.Method)
	assert.Equal(t, "/v1/Core/TimeReports/5", apiErr.Path)
	assert.Equal(t, "blikk: GET /v1/Core/TimeReports/5: 400 Bad Request: The request is invalid. (The timereport could not be found.)", apiErr.Error())
}

const missingProjectBody = `{"modelState":{"":["The project could not be found."]}}`

func TestMissingReferenceOnUpdateIsNotNotFound(t *testing.T) {
	client := newTestClient(t, respondWith(http.StatusBadRequest, missingProjectBody))

	err := client.TimeReports.Update(context.Background(), 5, TimeReportInput{ProjectID: 999999999})

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.MethodPut, apiErr.Method)
	assert.False(t, IsNotFound(err))
}

func TestMissingReferenceOnCreateIsNotNotFound(t *testing.T) {
	client := newTestClient(t, respondWith(http.StatusBadRequest, missingProjectBody))

	_, err := client.TimeReports.Create(context.Background(), TimeReportInput{ProjectID: 999999999})

	require.Error(t, err)
	assert.False(t, IsNotFound(err))
}

func TestCouldNotBeFoundOnGetAndDeleteIsNotFound(t *testing.T) {
	client := newTestClient(t, respondWith(http.StatusBadRequest, missingProjectBody))

	_, getErr := client.TimeReports.Get(context.Background(), 5)
	deleteErr := client.TimeReports.Delete(context.Background(), 5)

	assert.True(t, IsNotFound(getErr))
	assert.True(t, IsNotFound(deleteErr))
}

func TestRequestIDFallsBackToResponseHeader(t *testing.T) {
	apiErr := getUserError(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "r-from-header")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"The request is invalid."}`))
	})

	assert.Equal(t, "r-from-header", apiErr.RequestID)
	assert.Equal(t, "blikk: GET /v1/Admin/Users/1: 400 Bad Request: The request is invalid. (request r-from-header)", apiErr.Error())
}

func TestRequestIDInBodyWinsOverHeader(t *testing.T) {
	apiErr := getUserError(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "r-from-header")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"requestId":"r-from-body"}`))
	})

	assert.Equal(t, "r-from-body", apiErr.RequestID)
}

func TestModelStateKeysAreSortedAndEmptyKeyHasNoPrefix(t *testing.T) {
	apiErr := &Error{StatusCode: http.StatusBadRequest, ModelState: map[string][]string{"date": {"Required."}, "": {"The project could not be found."}}}

	assert.Equal(t, "blikk: 400 Bad Request (The project could not be found.; date: Required.)", apiErr.Error())
}

func TestValidationBadRequestIsNotNotFound(t *testing.T) {
	apiErr := getUserError(t, respondWith(http.StatusBadRequest, `{"message":"The request is invalid.","modelState":{"page":["The field Page must be greater than or equal to 1."]}}`))

	assert.False(t, IsNotFound(apiErr))
}

func TestNotFoundPhraseOnlyCountsForBadRequest(t *testing.T) {
	assert.False(t, IsNotFound(&Error{StatusCode: http.StatusInternalServerError, Message: "The user could not be found."}))
	assert.True(t, IsNotFound(&Error{Method: http.MethodGet, StatusCode: http.StatusBadRequest, Message: "THE USER COULD NOT BE FOUND."}))
}

func TestErrorHelpersSeeThroughWrapping(t *testing.T) {
	wrapped := fmt.Errorf("loading user: %w", &Error{StatusCode: http.StatusTooManyRequests})

	assert.True(t, IsRateLimited(wrapped))
	assert.False(t, IsNotFound(wrapped))
	assert.False(t, IsRateLimited(errors.New("plain")))
	assert.True(t, IsUnauthorized(&Error{StatusCode: http.StatusUnauthorized}))
}
