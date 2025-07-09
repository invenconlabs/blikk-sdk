package dateutils

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDateOnly_MarshalJSON(t *testing.T) {
	d := DateOnly{Time: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)}
	bytes, err := json.Marshal(d)
	require.NoError(t, err)
	assert.Equal(t, `"2024-01-02"`, string(bytes))

	var zero DateOnly
	bytes, err = json.Marshal(zero)
	require.NoError(t, err)
	assert.Equal(t, `null`, string(bytes))
}

func TestDateOnly_UnmarshalJSON(t *testing.T) {
	var d DateOnly
	err := json.Unmarshal([]byte(`"2024-01-02"`), &d)
	require.NoError(t, err)
	assert.Equal(t, 2024, d.Year())
	assert.Equal(t, time.January, d.Month())
	assert.Equal(t, 2, d.Day())

	var d2 DateOnly
	err = json.Unmarshal([]byte(`null`), &d2)
	require.NoError(t, err)
	assert.True(t, d2.IsZero())
}

func TestPreviousWeek(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Stockholm")
	require.NoError(t, err)

	testCases := []struct {
		name     string
		today    time.Time
		wantFrom string
		wantTo   string
	}{
		{"Mid-week Wednesday", time.Date(2024, 3, 27, 12, 0, 0, 0, loc), "2024-03-18", "2024-03-24"},
		{"Sunday", time.Date(2024, 3, 24, 12, 0, 0, 0, loc), "2024-03-11", "2024-03-17"},
		{"Monday", time.Date(2024, 3, 25, 12, 0, 0, 0, loc), "2024-03-18", "2024-03-24"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			from, to := PreviousWeek(tc.today)
			assert.Equal(t, tc.wantFrom, from.Format(time.DateOnly))
			assert.Equal(t, tc.wantTo, to.Format(time.DateOnly))
		})
	}
}

func TestPreviousMonth(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Stockholm")
	require.NoError(t, err)

	testCases := []struct {
		name     string
		today    time.Time
		wantFrom string
		wantTo   string
	}{
		{"Mid-month", time.Date(2024, 4, 15, 0, 0, 0, 0, loc), "2024-03-01", "2024-03-31"},
		{"Beginning of month", time.Date(2024, 3, 1, 0, 0, 0, 0, loc), "2024-02-01", "2024-02-29"}, // Leap year
		{"January", time.Date(2024, 1, 10, 0, 0, 0, 0, loc), "2023-12-01", "2023-12-31"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			from, to := PreviousMonth(tc.today)
			assert.Equal(t, tc.wantFrom, from.Format(time.DateOnly))
			assert.Equal(t, tc.wantTo, to.Format(time.DateOnly))
		})
	}
}

func TestFirstDayOfMonth(t *testing.T) {
	d := FirstDayOfMonth(2024, time.February)
	assert.Equal(t, "2024-02-01", d.Format(time.DateOnly))
}

func TestLastDayOfMonth(t *testing.T) {
	d := LastDayOfMonth(2024, time.February) // Leap year
	assert.Equal(t, "2024-02-29", d.Format(time.DateOnly))

	d2 := LastDayOfMonth(2023, time.February) // Not a leap year
	assert.Equal(t, "2023-02-28", d2.Format(time.DateOnly))

	d3 := LastDayOfMonth(2024, time.December)
	assert.Equal(t, "2024-12-31", d3.Format(time.DateOnly))
}
