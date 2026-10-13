package blikk

import (
	"encoding/json"
	"strings"
)

func unmarshalEnum[T ~string](data []byte, target *T, known ...T) error {
	var raw *string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if raw == nil {
		*target = ""
		return nil
	}
	*target = matchEnum(*raw, known)
	return nil
}

func matchEnum[T ~string](raw string, known []T) T {
	for _, constant := range known {
		if strings.EqualFold(raw, string(constant)) {
			return constant
		}
	}
	return T(raw)
}
