package blikk

import (
	"cmp"
	"context"
	"iter"
	"maps"
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
	query, err := encodeQuery(params)
	if err != nil {
		return failedSeq[T](err)
	}
	setIfAbsent(query, "page", "1")
	setIfAbsent(query, "pageSize", strconv.Itoa(MaxPageSize))
	firstPage, _ := strconv.Atoi(query.Get("page"))
	return walkPages(firstPage, func(pageNumber int) (*Page[T], error) {
		pageQuery := maps.Clone(query)
		pageQuery.Set("page", strconv.Itoa(pageNumber))
		return fetchPage[T](ctx, c, path, pageQuery)
	})
}

type pagedBody struct {
	Filter   any `json:"filter,omitempty"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize,omitempty"`
}

func listPageWithBody[T any](ctx context.Context, c *Client, path string, body pagedBody) (*Page[T], error) {
	body.Page = cmp.Or(body.Page, 1)
	return fetchPageWithBody[T](ctx, c, path, body)
}

func listAllWithBody[T any](ctx context.Context, c *Client, path string, body pagedBody) iter.Seq2[T, error] {
	body.PageSize = cmp.Or(body.PageSize, MaxPageSize)
	return walkPages(cmp.Or(body.Page, 1), func(pageNumber int) (*Page[T], error) {
		pageBody := body
		pageBody.Page = pageNumber
		return fetchPageWithBody[T](ctx, c, path, pageBody)
	})
}

func walkPages[T any](firstPage int, fetch func(pageNumber int) (*Page[T], error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for pageNumber := firstPage; ; pageNumber++ {
			page, err := fetch(pageNumber)
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
		}
	}
}

func failedSeq[T any](err error) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) { yieldError(yield, err) }
}

func fetchPage[T any](ctx context.Context, c *Client, path string, query url.Values) (*Page[T], error) {
	var page Page[T]
	if err := c.do(ctx, http.MethodGet, path, query, nil, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

func fetchPageWithBody[T any](ctx context.Context, c *Client, path string, body pagedBody) (*Page[T], error) {
	var page Page[T]
	if err := c.do(ctx, http.MethodPost, path, nil, body, &page); err != nil {
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
