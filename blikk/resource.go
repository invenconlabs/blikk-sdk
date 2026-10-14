package blikk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func getResource[T any](ctx context.Context, c *Client, base string, pathParams ...string) (*T, error) {
	path, err := resourcePath(base, pathParams...)
	if err != nil {
		return nil, err
	}
	var resource T
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

func getResourceWithQuery[T any](ctx context.Context, c *Client, path string, query url.Values) (*T, error) {
	var resource T
	if err := c.do(ctx, http.MethodGet, path, query, nil, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

func createResource[T any](ctx context.Context, c *Client, input any, base string, pathParams ...string) (*T, error) {
	path, err := resourcePath(base, pathParams...)
	if err != nil {
		return nil, err
	}
	var created T
	if err := c.do(ctx, http.MethodPost, path, nil, input, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

type createdResource struct {
	ID int `json:"id"`
}

// encoding/json matches keys case-insensitively, so this also reads the {"Id": 42} some endpoints return.
func createResourceID(ctx context.Context, c *Client, input any, base string, pathParams ...string) (int, error) {
	created, err := createResource[*createdResource](ctx, c, input, base, pathParams...)
	if err != nil {
		return 0, err
	}
	return createdID(*created, base)
}

// A nil response means the API answered the create with an empty body; the create still succeeded.
func createdID(created *createdResource, base string) (int, error) {
	switch {
	case created == nil:
		return 0, nil
	case created.ID == 0:
		return 0, fmt.Errorf("blikk: create response of %s has no id", base)
	default:
		return created.ID, nil
	}
}

func uploadResourceID(ctx context.Context, c *Client, fields map[string]string, files []UploadFile, base string, pathParams ...string) (int, error) {
	path, err := resourcePath(base, pathParams...)
	if err != nil {
		return 0, err
	}
	var created *createdResource
	if err := c.upload(ctx, path, fields, files, &created); err != nil {
		return 0, err
	}
	return createdID(created, base)
}

func updateResource(ctx context.Context, c *Client, input any, base string, pathParams ...string) error {
	path, err := resourcePath(base, pathParams...)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodPut, path, nil, input, nil)
}

func updateResourceReturning[T any](ctx context.Context, c *Client, input any, base string, pathParams ...string) (*T, error) {
	path, err := resourcePath(base, pathParams...)
	if err != nil {
		return nil, err
	}
	var updated T
	if err := c.do(ctx, http.MethodPut, path, nil, input, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func patchResource(ctx context.Context, c *Client, path string, query url.Values, body any) error {
	return c.do(ctx, http.MethodPatch, path, query, body, nil)
}

func deleteResource(ctx context.Context, c *Client, base string, pathParams ...string) error {
	path, err := resourcePath(base, pathParams...)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
