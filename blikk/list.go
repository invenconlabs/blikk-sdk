package blikk

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
)

const MaxPageSize = 100

type SortOrder string

const (
	SortAscending  SortOrder = "ascending"
	SortDescending SortOrder = "descending"
)

type ListOptions struct {
	Page      int       `query:"page"`
	PageSize  int       `query:"pageSize"`
	SortBy    string    `query:"sortBy"`
	SortOrder SortOrder `query:"sortOrder"`
}

type Page[T any] struct {
	Page           int `json:"page"`
	PageSize       int `json:"pageSize"`
	ItemCount      int `json:"itemCount"`
	TotalItemCount int `json:"totalItemCount"`
	TotalPages     int `json:"totalPages"`
	Items          []T `json:"items"`
}

func listPage[T any](ctx context.Context, c *Client, path string, params any) (*Page[T], error) {
	query, err := encodeQuery(params)
	if err != nil {
		return nil, err
	}
	setIfAbsent(query, "page", "1")
	return fetchPage[T](ctx, c, path, query)
}

func listAll[T any](ctx context.Context, c *Client, path string, params any) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		query, err := encodeQuery(params)
		if err != nil {
			yieldError(yield, err)
			return
		}
		setIfAbsent(query, "page", "1")
		setIfAbsent(query, "pageSize", strconv.Itoa(MaxPageSize))
		pageNumber, _ := strconv.Atoi(query.Get("page"))
		for {
			query.Set("page", strconv.Itoa(pageNumber))
			page, err := fetchPage[T](ctx, c, path, query)
			if err != nil {
				yieldError(yield, err)
				return
			}
			for _, item := range page.Items {
				if !yield(item, nil) {
					return
				}
			}
			if pageNumber >= page.TotalPages {
				return
			}
			pageNumber++
		}
	}
}

func fetchPage[T any](ctx context.Context, c *Client, path string, query url.Values) (*Page[T], error) {
	var page Page[T]
	if err := c.do(ctx, http.MethodGet, path, query, nil, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

func setIfAbsent(query url.Values, key, value string) {
	if query.Get(key) == "" {
		query.Set(key, value)
	}
}

func yieldError[T any](yield func(T, error) bool, err error) {
	var zero T
	yield(zero, err)
}
