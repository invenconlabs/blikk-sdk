package blikk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type Error struct {
	StatusCode          int
	Message             string
	ModelState          map[string][]string
	PermissionsRequired []string
	IPAddressAllowed    *bool
	RequestID           string
	Body                []byte
}

type errorBody struct {
	Message             string              `json:"message"`
	ModelState          map[string][]string `json:"modelState"`
	PermissionsRequired []string            `json:"permissionsRequired"`
	IPAddressAllowed    *bool               `json:"ipAddressAllowed"`
	RequestID           string              `json:"requestId"`
}

func newError(statusCode int, body []byte) *Error {
	var decoded errorBody
	_ = json.Unmarshal(body, &decoded)
	return &Error{
		StatusCode:          statusCode,
		Message:             decoded.Message,
		ModelState:          decoded.ModelState,
		PermissionsRequired: decoded.PermissionsRequired,
		IPAddressAllowed:    decoded.IPAddressAllowed,
		RequestID:           decoded.RequestID,
		Body:                body,
	}
}

func (e *Error) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "blikk: %d %s", e.StatusCode, http.StatusText(e.StatusCode))
	if e.Message != "" {
		b.WriteString(": " + e.Message)
	}
	if len(e.ModelState) > 0 {
		b.WriteString(" (" + formatModelState(e.ModelState) + ")")
	}
	if len(e.PermissionsRequired) > 0 {
		b.WriteString(" (requires " + strings.Join(e.PermissionsRequired, ", ") + ")")
	}
	if e.RequestID != "" {
		b.WriteString(" (request " + e.RequestID + ")")
	}
	return b.String()
}

func formatModelState(modelState map[string][]string) string {
	fields := make([]string, 0, len(modelState))
	for field := range modelState {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	parts := make([]string, len(fields))
	for i, field := range fields {
		parts[i] = field + ": " + strings.Join(modelState[field], ", ")
	}
	return strings.Join(parts, "; ")
}

func IsNotFound(err error) bool { return hasStatus(err, http.StatusNotFound) }

func IsRateLimited(err error) bool { return hasStatus(err, http.StatusTooManyRequests) }

func IsUnauthorized(err error) bool { return hasStatus(err, http.StatusUnauthorized) }

func hasStatus(err error, statusCode int) bool {
	var apiErr *Error
	return errors.As(err, &apiErr) && apiErr.StatusCode == statusCode
}
