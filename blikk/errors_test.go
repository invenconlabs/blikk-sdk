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
		"blikk: 400 Bad Request: The request is invalid. (filter.From: Invalid datespan.; page: The field Page must be greater than or equal to 1.)",
		apiErr.Error())
}

func TestForbiddenErrorDecodesPermissions(t *testing.T) {
	body := `{"message":"You don't have the necessary permissions to perform this action","ipAddressAllowed":true,"permissionsRequired":["Contact_Read"]}`
	apiErr := getUserError(t, respondWith(http.StatusForbidden, body))

	assert.Equal(t, http.StatusForbidden, apiErr.StatusCode)
	assert.Equal(t, []string{"Contact_Read"}, apiErr.PermissionsRequired)
	require.NotNil(t, apiErr.IPAddressAllowed)
	assert.True(t, *apiErr.IPAddressAllowed)
	assert.Equal(t, "blikk: 403 Forbidden: You don't have the necessary permissions to perform this action", apiErr.Error())
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
	assert.Equal(t, "blikk: 404 Not Found", err.Error())
}

func TestServerErrorDecodesRequestID(t *testing.T) {
	apiErr := getUserError(t, respondWith(http.StatusInternalServerError, `{"message":"An unexpected error has occurred.","requestId":"r-0fbf44e8-afb8-42be-ae69-55c8bc03311b"}`))

	assert.Equal(t, "r-0fbf44e8-afb8-42be-ae69-55c8bc03311b", apiErr.RequestID)
	assert.Equal(t, "blikk: 500 Internal Server Error: An unexpected error has occurred. (request r-0fbf44e8-afb8-42be-ae69-55c8bc03311b)", apiErr.Error())
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

func TestErrorHelpersSeeThroughWrapping(t *testing.T) {
	wrapped := fmt.Errorf("loading user: %w", &Error{StatusCode: http.StatusTooManyRequests})

	assert.True(t, IsRateLimited(wrapped))
	assert.False(t, IsNotFound(wrapped))
	assert.False(t, IsRateLimited(errors.New("plain")))
	assert.True(t, IsUnauthorized(&Error{StatusCode: http.StatusUnauthorized}))
}
