package blikk

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type queryValuer interface {
	queryValue() string
}

var queryValuerType = reflect.TypeFor[queryValuer]()

func encodeQuery(params any) (url.Values, error) {
	query := url.Values{}
	if params == nil {
		return query, nil
	}
	value := reflect.ValueOf(params)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("blikk: query params must be a struct, got %s", value.Type())
	}
	if err := encodeStructFields(query, value); err != nil {
		return nil, err
	}
	return query, nil
}

func encodeStructFields(query url.Values, structValue reflect.Value) error {
	structType := structValue.Type()
	for i := range structType.NumField() {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)
		name, tagged := field.Tag.Lookup("query")
		if field.Anonymous && !tagged && fieldValue.Kind() == reflect.Struct {
			if err := encodeStructFields(query, fieldValue); err != nil {
				return err
			}
			continue
		}
		if !tagged || isUnsetQueryValue(fieldValue) {
			continue
		}
		encoded, err := encodeQueryValue(fieldValue)
		if err != nil {
			return fmt.Errorf("blikk: encoding query field %s: %w", field.Name, err)
		}
		query.Set(name, encoded)
	}
	return nil
}

func isUnsetQueryValue(value reflect.Value) bool {
	return value.IsZero() || (value.Kind() == reflect.Slice && value.Len() == 0)
}

func encodeQueryValue(value reflect.Value) (string, error) {
	if value.Type().Implements(queryValuerType) {
		return value.Interface().(queryValuer).queryValue(), nil
	}
	switch value.Kind() {
	case reflect.String:
		return value.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), nil
	case reflect.Bool:
		return strconv.FormatBool(value.Bool()), nil
	case reflect.Pointer:
		if value.Elem().Kind() != reflect.Bool {
			return "", fmt.Errorf("unsupported pointer type %s", value.Type())
		}
		return strconv.FormatBool(value.Elem().Bool()), nil
	case reflect.Slice:
		return joinSliceValues(value)
	default:
		return "", fmt.Errorf("unsupported type %s", value.Type())
	}
}

func joinSliceValues(slice reflect.Value) (string, error) {
	parts := make([]string, slice.Len())
	for i := range slice.Len() {
		element := slice.Index(i)
		switch element.Kind() {
		case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			encoded, err := encodeQueryValue(element)
			if err != nil {
				return "", err
			}
			parts[i] = encoded
		default:
			return "", fmt.Errorf("unsupported slice element type %s", element.Type())
		}
	}
	return strings.Join(parts, ","), nil
}
