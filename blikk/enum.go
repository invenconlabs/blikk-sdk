package blikk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func unmarshalEnum[T ~string](data []byte, target *T, known ...T) error {
	raw, err := decodeEnumValue(data)
	if err != nil {
		return err
	}
	switch value := raw.(type) {
	case nil:
		*target = ""
	case string:
		*target = matchEnum(value, known)
	case json.Number:
		*target = T(value.String())
	default:
		return fmt.Errorf("blikk: enum must be a string or number, got %s", data)
	}
	return nil
}

func decodeEnumValue(data []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var raw any
	err := decoder.Decode(&raw)
	return raw, err
}

func matchEnum[T ~string](raw string, known []T) T {
	for _, constant := range known {
		if strings.EqualFold(raw, string(constant)) {
			return constant
		}
	}
	return T(raw)
}
