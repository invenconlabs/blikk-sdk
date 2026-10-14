package blikk

import (
	"context"
	"net/http"
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

func updateResource(ctx context.Context, c *Client, input any, base string, pathParams ...string) error {
	path, err := resourcePath(base, pathParams...)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodPut, path, nil, input, nil)
}

func deleteResource(ctx context.Context, c *Client, base string, pathParams ...string) error {
	path, err := resourcePath(base, pathParams...)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
