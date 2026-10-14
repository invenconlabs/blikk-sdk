package blikk

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDateMarshalJSON(t *testing.T) {
	data, err := json.Marshal(NewDate(2024, time.January, 2))
	require.NoError(t, err)
	assert.Equal(t, `"2024-01-02"`, string(data))

	data, err = json.Marshal(Date{})
	require.NoError(t, err)
	assert.Equal(t, `null`, string(data))
}

func TestDateUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Date
	}{
		{"date only", `"2024-01-02"`, NewDate(2024, time.January, 2)},
		{"datetime with offset", `"2021-08-30 00:00:00 +00:00"`, NewDate(2021, time.August, 30)},
		{"iso datetime", `"2021-08-30T00:00:00"`, NewDate(2021, time.August, 30)},
		{"midnight at plus two keeps its own day", `"2026-10-13T00:00:00+02:00"`, NewDate(2026, time.October, 13)},
		{"late evening at minus five keeps its own day", `"2026-10-13 23:30:00 -05:00"`, NewDate(2026, time.October, 13)},
		{"null", `null`, Date{}},
		{"empty string", `""`, Date{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got Date
			require.NoError(t, json.Unmarshal([]byte(tc.input), &got))
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDateUnmarshalJSONRejectsGarbage(t *testing.T) {
	var d Date
	assert.Error(t, json.Unmarshal([]byte(`"not a date"`), &d))
	assert.Error(t, json.Unmarshal([]byte(`20240102`), &d))
}

func TestDateNullClearsPreviousValue(t *testing.T) {
	d := NewDate(2024, time.January, 2)
	require.NoError(t, json.Unmarshal([]byte(`null`), &d))
	assert.True(t, d.IsZero())
}

func TestDateOfKeepsCalendarDayOfLocation(t *testing.T) {
	stockholm := time.FixedZone("CEST", 2*60*60)
	lateEvening := time.Date(2026, time.October, 13, 23, 30, 0, 0, stockholm)

	assert.Equal(t, "2026-10-13", DateOf(lateEvening).String())
	assert.Equal(t, time.Date(2026, time.October, 13, 0, 0, 0, 0, time.UTC), DateOf(lateEvening).Time())
}

func TestDateNormalisesOverflow(t *testing.T) {
	assert.Equal(t, "2024-02-29", NewDate(2024, time.March, 0).String())
	assert.Equal(t, "2023-02-28", NewDate(2023, time.March, 0).String())
	assert.Equal(t, "2025-01-01", NewDate(2024, time.December, 32).String())
}

func TestDateString(t *testing.T) {
	assert.Equal(t, "2026-09-01", NewDate(2026, time.September, 1).String())
	assert.Equal(t, "", Date{}.String())
}

func TestParseDateTimeFormats(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Time
	}{
		{"space separated with offset", "2026-10-13 09:40:57 +00:00", time.Date(2026, 10, 13, 9, 40, 57, 0, time.UTC)},
		{"space separated with non-utc offset", "2026-10-13 11:40:57 +02:00", time.Date(2026, 10, 13, 9, 40, 57, 0, time.UTC)},
		{"iso with millis and offset", "2026-10-13T09:40:57.289+00:00", time.Date(2026, 10, 13, 9, 40, 57, 289_000_000, time.UTC)},
		{"iso with 7 fractional digits", "2026-10-05T09:06:52.2277519+00:00", time.Date(2026, 10, 5, 9, 6, 52, 227_751_900, time.UTC)},
		{"rfc3339 zulu", "2026-10-13T09:40:57Z", time.Date(2026, 10, 13, 9, 40, 57, 0, time.UTC)},
		{"iso without offset", "2026-10-13T09:40:57", time.Date(2026, 10, 13, 9, 40, 57, 0, time.UTC)},
		{"space separated without offset", "2021-11-10 00:00:00", time.Date(2021, 11, 10, 0, 0, 0, 0, time.UTC)},
		{"date only", "2026-10-13", time.Date(2026, 10, 13, 0, 0, 0, 0, time.UTC)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseDateTime(tc.input)
			require.NoError(t, err)
			assert.True(t, tc.want.Equal(got.Time()), "got %s, want %s", got.Time(), tc.want)
		})
	}
}

func TestParseDateTimeRejectsGarbage(t *testing.T) {
	_, err := ParseDateTime("13/10/2026")
	assert.Error(t, err)
}

func TestDateTimeUnmarshalJSON(t *testing.T) {
	var dt DateTime
	require.NoError(t, json.Unmarshal([]byte(`"2026-10-13 09:40:57 +00:00"`), &dt))
	assert.Equal(t, "2026-10-13 09:40:57", dt.String())

	require.NoError(t, json.Unmarshal([]byte(`null`), &dt))
	assert.True(t, dt.IsZero())

	var pointer *DateTime
	require.NoError(t, json.Unmarshal([]byte(`null`), &pointer))
	assert.Nil(t, pointer)
}

func TestDateTimeMarshalJSONUsesUTC(t *testing.T) {
	stockholm := time.FixedZone("CEST", 2*60*60)
	data, err := json.Marshal(DateTimeOf(time.Date(2026, 10, 13, 11, 40, 57, 0, stockholm)))
	require.NoError(t, err)
	assert.Equal(t, `"2026-10-13 09:40:57"`, string(data))

	data, err = json.Marshal(DateTime{})
	require.NoError(t, err)
	assert.Equal(t, `null`, string(data))
}

func TestDateTimeRoundTrip(t *testing.T) {
	original := DateTimeOf(time.Date(2026, 10, 13, 9, 40, 57, 0, time.UTC))
	data, err := json.Marshal(original)
	require.NoError(t, err)
	var decoded DateTime
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.True(t, original.Time().Equal(decoded.Time()))
}
