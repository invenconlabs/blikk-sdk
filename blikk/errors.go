package blikk

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type Error struct {
	Method              string
	Path                string
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

const requestIDHeader = "X-Request-Id"

func newError(method, path string, statusCode int, header http.Header, body []byte) *Error {
	var decoded errorBody
	_ = json.Unmarshal(body, &decoded)
	return &Error{
		Method:              method,
		Path:                path,
		StatusCode:          statusCode,
		Message:             decoded.Message,
		ModelState:          decoded.ModelState,
		PermissionsRequired: decoded.PermissionsRequired,
		IPAddressAllowed:    decoded.IPAddressAllowed,
		RequestID:           cmp.Or(decoded.RequestID, header.Get(requestIDHeader)),
		Body:                body,
	}
}

func (e *Error) Error() string {
	var b strings.Builder
	b.WriteString("blikk: ")
	if e.Method != "" {
		b.WriteString(e.Method + " " + e.Path + ": ")
	}
	fmt.Fprintf(&b, "%d %s", e.StatusCode, http.StatusText(e.StatusCode))
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
		parts[i] = formatModelStateEntry(field, modelState[field])
	}
	return strings.Join(parts, "; ")
}

func formatModelStateEntry(field string, messages []string) string {
	joined := strings.Join(messages, ", ")
	if field == "" {
		return joined
	}
	return field + ": " + joined
}

const notFoundPhrase = "could not be found"

// IsNotFound is also true for a GET or DELETE 400 saying "could not be found"; on writes that 400 means a referenced entity is missing.
func IsNotFound(err error) bool {
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		return false
	}
	return apiErr.StatusCode == http.StatusNotFound || isBadRequestForMissingEntity(apiErr)
}

func isBadRequestForMissingEntity(apiErr *Error) bool {
	if apiErr.StatusCode != http.StatusBadRequest || !isReadOrDelete(apiErr.Method) {
		return false
	}
	if mentionsNotFound(apiErr.Message) {
		return true
	}
	for _, messages := range apiErr.ModelState {
		for _, message := range messages {
			if mentionsNotFound(message) {
				return true
			}
		}
	}
	return false
}

func isReadOrDelete(method string) bool {
	return method == http.MethodGet || method == http.MethodDelete
}

func mentionsNotFound(message string) bool {
	return strings.Contains(strings.ToLower(message), notFoundPhrase)
}

func IsRateLimited(err error) bool { return hasStatus(err, http.StatusTooManyRequests) }

func IsUnauthorized(err error) bool { return hasStatus(err, http.StatusUnauthorized) }

func hasStatus(err error, statusCode int) bool {
	var apiErr *Error
	return errors.As(err, &apiErr) && apiErr.StatusCode == statusCode
}
