# Blikk SDK for Go

This is an unofficial Go SDK for the Blikk Public API.

## Table of Contents

- [Installation](#installation)
- [Authentication](#authentication)
- [Usage](#usage)
  - [Creating a Client](#creating-a-client)
  - [Listing Resources](#listing-resources)
  - [Getting a Single Resource](#getting-a-single-resource)
- [Available Resources](#available-resources)
  - [Listable Resources](#listable-resources)
  - [Gettable Resources](#gettable-resources)
- [Filtering and Pagination](#filtering-and-pagination)
  - [Filtering](#filtering)
  - [Pagination](#pagination)
- [Configuration](#configuration)
  - [Custom Base URL](#custom-base-url)
  - [Custom HTTP Client](#custom-http-client)
- [Error Handling](#error-handling)

## Installation

```bash
go get github.com/invenconlabs/blikk-sdk
```

## Authentication

To use the Blikk API, you need an access token. The SDK provides a helper function to retrieve an access token using your App ID and App Secret. These credentials should be stored as environment variables:

- `BLIKK_APP_ID`: Your Blikk App ID.
- `BLIKK_APP_SECRET`: Your Blikk App Secret.

You can then use the `GetAccessToken` function to get a token:

```go
package main

import (
	"fmt"
	"log"

	"github.com/invenconlabs/blikk-sdk/blikk"
)

func main() {
	token, err := blikk.GetAccessToken()
	if err != nil {
		log.Fatalf("failed to get access token: %v", err)
	}

	fmt.Println("Successfully retrieved access token:", token)
}
```

## Usage

### Creating a Client

Once you have an access token, you can create a new `Client`:

```go
client := blikk.NewClient(token)
```

### Listing Resources

The `List` function allows you to retrieve a collection of resources. It handles pagination automatically, so you get all the results in a single call.

Here's an example of how to list all users:

```go
users, err := blikk.List[blikk.Users](client, blikk.NewListOptions())
if err != nil {
	log.Fatalf("failed to list users: %v", err)
}

for _, user := range users {
	fmt.Printf("User: %s %s\n", user.FirstName, user.LastName)
}
```

### Getting a Single Resource

The `Get` function allows you to retrieve a single resource by its ID.

Here's an example of how to get a single user:

```go
user, err := blikk.Get[blikk.User](client, "12345") // Get user with ID 12345
if err != nil {
	log.Fatalf("failed to get user: %v", err)
}

fmt.Printf("User: %s %s\n", user.FirstName, user.LastName)
```

## Available Resources

The following resources are available through the SDK:

### Listable Resources
- `blikk.Users`: List of users in the system.
- `blikk.Projects`: List of projects.
- `blikk.TimeReports`: List of time reports.
- `blikk.UserDayStatistics`: Daily time statistics for users.

### Gettable Resources
- `blikk.User`: Detailed information for a single user.

## Filtering and Pagination

The `List` function accepts `ListOptions` to filter and paginate the results.

### Filtering

You can filter resources based on certain criteria. For example, to get time reports for a specific user within a date range:

```go
import "github.com/invenconlabs/blikk-sdk/dateutils"

options := blikk.NewListOptions()
options.UserIDs = []uint16{123} // Filter by user ID
options.FromDate = dateutils.FirstDayOfMonth(time.Now())
options.ToDate = dateutils.LastDayOfMonth(time.Now())

timeReports, err := blikk.List[blikk.TimeReports](client, options)
// ...
```

**Note:** Not all filters are applicable to all resources. Please refer to the `validFilter` method for each resource in `models.go` to see which filters are supported.

### Pagination

Pagination is handled automatically by the `List` function. You can, however, set the page size in `ListOptions`:

```go
options := blikk.NewListOptions()
options.PageSize = 50 // Retrieve 50 items per page

// The List function will still fetch all pages and return a complete slice.
```

## Configuration

### Custom Base URL

For testing or to use a different API environment, you can override the default base URL:

```go
client := blikk.NewClient(token, blikk.WithBaseURL("https://your-custom-api-url.com/"))
```

### Custom HTTP Client

You can also provide your own `http.Client`, for example, to set custom timeouts or use a mock client for testing:

```go
import "net/http"
import "time"

httpClient := &http.Client{Timeout: 10 * time.Second}
client := blikk.NewClient(token, blikk.WithHTTPClient(httpClient))
```

## Error Handling

The SDK functions return an error if the API request fails or if there's an issue with processing the request or response. The client also has built-in retry logic for `429 Too Many Requests` errors, respecting the `Retry-After` header sent by the API.

It's important to check for errors on every call:

```go
users, err := blikk.List[blikk.Users](client, blikk.NewListOptions())
if err != nil {
	log.Fatalf("API error: %v", err)
}
```