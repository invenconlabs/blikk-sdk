package blikk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04:05"
)

var dateTimeLayouts = []string{
	"2006-01-02 15:04:05 -07:00",
	time.RFC3339,
	"2006-01-02T15:04:05",
	dateTimeLayout,
	dateLayout,
}

var jsonNull = []byte("null")

type Date struct {
	t time.Time
}

func NewDate(year int, month time.Month, day int) Date {
	return Date{t: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

func DateOf(t time.Time) Date {
	return NewDate(t.Date())
}

func ParseDate(s string) (Date, error) {
	if t, err := time.Parse(dateLayout, s); err == nil {
		return Date{t: t}, nil
	}
	t, err := parseInOriginalOffset(s)
	if err != nil {
		return Date{}, fmt.Errorf("blikk: parsing date %q: %w", s, err)
	}
	return DateOf(t), nil
}

func (d Date) Time() time.Time { return d.t }

func (d Date) IsZero() bool { return d.t.IsZero() }

func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return d.t.Format(dateLayout)
}

func (d Date) queryValue() string { return d.String() }

func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return jsonNull, nil
	}
	return []byte(strconv.Quote(d.String())), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	s, isNull, err := unquoteJSONString(data)
	if err != nil || isNull {
		*d = Date{}
		return err
	}
	parsed, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

type DateTime struct {
	t time.Time
}

func DateTimeOf(t time.Time) DateTime {
	return DateTime{t: t}
}

// ParseDateTime accepts every datetime format the Blikk API has been seen to return and normalises the result to UTC.
func ParseDateTime(s string) (DateTime, error) {
	t, err := parseInOriginalOffset(s)
	if err != nil {
		return DateTime{}, err
	}
	return DateTime{t: t.UTC()}, nil
}

func parseInOriginalOffset(s string) (time.Time, error) {
	for _, layout := range dateTimeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("blikk: unrecognised datetime %q", s)
}

func (dt DateTime) Time() time.Time { return dt.t }

func (dt DateTime) IsZero() bool { return dt.t.IsZero() }

func (dt DateTime) String() string {
	if dt.IsZero() {
		return ""
	}
	return dt.t.UTC().Format(dateTimeLayout)
}

func (dt DateTime) queryValue() string { return dt.String() }

func (dt DateTime) MarshalJSON() ([]byte, error) {
	if dt.IsZero() {
		return jsonNull, nil
	}
	return []byte(strconv.Quote(dt.String())), nil
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	s, isNull, err := unquoteJSONString(data)
	if err != nil || isNull {
		*dt = DateTime{}
		return err
	}
	parsed, err := ParseDateTime(s)
	if err != nil {
		return err
	}
	*dt = parsed
	return nil
}

func unquoteJSONString(data []byte) (value string, isNull bool, err error) {
	if bytes.Equal(data, jsonNull) {
		return "", true, nil
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return "", false, fmt.Errorf("blikk: expected a JSON string, got %s", data)
	}
	return value, value == "", nil
}
