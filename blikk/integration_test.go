//go:build integration

package blikk

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupIntegrationTest initializes a client for integration tests.
// It skips the test if the required environment variables are not set.
func setupIntegrationTest(t *testing.T) *Client {
	appID := os.Getenv("BLIKK_APP_ID")
	appSecret := os.Getenv("BLIKK_APP_SECRET")

	if appID == "" || appSecret == "" {
		t.Skip("BLIKK_APP_ID and BLIKK_APP_SECRET must be set for integration tests")
	}

	token, err := GetAccessToken()
	require.NoError(t, err, "Failed to get access token for integration test")
	require.NotEmpty(t, token)

	return NewClient(token)
}

func TestIntegration_ListUsers(t *testing.T) {
	client := setupIntegrationTest(t)
	options := NewListOptions()

	users, err := List[Users](client, options)
	require.NoError(t, err)
	assert.NotEmpty(t, users, "Expected to find at least one user")
	fmt.Printf("Found %d users\n", len(users))
}

func TestIntegration_ListTimeReportsForPreviousWeek(t *testing.T) {
	client := setupIntegrationTest(t)
	options := NewListOptions()
	from, to := PreviousWeek()
	options.FromDate = &from
	options.ToDate = &to

	reports, err := List[TimeReports](client, options)
	require.NoError(t, err)
	// It's okay if there are no reports, so we don't assert NotEmpty
	fmt.Printf("Found %d time reports for previous week (%s to %s)\n", len(reports), from.Format(time.DateOnly), to.Format(time.DateOnly))
}

func TestIntegration_GetUser(t *testing.T) {
	client := setupIntegrationTest(t)
	options := NewListOptions()

	// First, list users to get a valid ID
	users, err := List[Users](client, options)
	require.NoError(t, err)
	require.NotEmpty(t, users, "Cannot test GetUser without at least one user to fetch")

	userID := users[0].ID
	userFirstName := users[0].FirstName
	userLastName := users[0].LastName

	// Now, get the specific user
	user, err := Get[User](client, fmt.Sprintf("%d", userID))
	require.NoError(t, err)

	assert.Equal(t, userID, user.ID)
	assert.Equal(t, userFirstName, user.FirstName)
	assert.Equal(t, userLastName, user.LastName)
	fmt.Printf("Successfully fetched user %s %s\n", user.FirstName, user.LastName)
}
