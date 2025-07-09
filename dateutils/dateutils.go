package dateutils

import (
	"log"
	"strings"
	"time"
)

// Date is a custom time.Time that handles the "YYYY-MM-DD" format.
type DateOnly struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler to parse "YYYY-MM-DD" formatted dates.
func (d *DateOnly) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
		d.Time = time.Time{}
		return
	}
	d.Time, err = time.Parse(time.DateOnly, s)
	return
}

// MarshalJSON implements json.Marshaler to format dates as "YYYY-MM-DD".
func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Format(time.DateOnly) + `"`), nil
}

// PreviousWeek calculates the start (Monday) and end (Sunday) dates of the previous week in the Stockholm timezone.
func PreviousWeek(t ...time.Time) (DateOnly, DateOnly) {
	loc, err := time.LoadLocation("Europe/Stockholm")
	if err != nil {
		log.Fatalf("Failed to load location Europe/Stockholm: %v", err)
	}

	var now time.Time
	if len(t) > 0 {
		now = t[0].In(loc)
	} else {
		now = time.Now().In(loc)
	}

	// Navigate to the most recent Monday.
	// If today is Sunday (0), go back 6 days to Monday.
	// If today is Monday (1), go back 0 days.
	// If today is Tuesday (2), go back 1 day.
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7 // Treat Sunday as the 7th day to simplify logic
	}
	offset := 1 - int(weekday)
	mostRecentMonday := now.AddDate(0, 0, offset)

	// We want the *previous* week, so subtract 7 days.
	previousMonday := mostRecentMonday.AddDate(0, 0, -7)
	previousSunday := previousMonday.AddDate(0, 0, 6)

	year, month, day := previousMonday.Date()
	previousMonday = time.Date(year, month, day, 0, 0, 0, 0, loc)

	year, month, day = previousSunday.Date()
	previousSunday = time.Date(year, month, day, 0, 0, 0, 0, loc)

	return DateOnly{previousMonday}, DateOnly{previousSunday}
}

// PreviousMonth calculates the first and last day of the previous month in the Stockholm timezone.
func PreviousMonth(t ...time.Time) (DateOnly, DateOnly) {
	loc, err := time.LoadLocation("Europe/Stockholm")
	if err != nil {
		log.Fatalf("Failed to load location Europe/Stockholm: %v", err)
	}
	var now time.Time
	if len(t) > 0 {
		now = t[0].In(loc)
	} else {
		now = time.Now().In(loc)
	}

	// First day of the current month in the specified location
	firstDayCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)

	// Last day of the previous month is one day before the first day of the current month
	lastDayPreviousMonth := firstDayCurrentMonth.AddDate(0, 0, -1)

	// First day of the previous month in the specified location
	firstDayPreviousMonth := time.Date(lastDayPreviousMonth.Year(), lastDayPreviousMonth.Month(), 1, 0, 0, 0, 0, loc)

	return DateOnly{firstDayPreviousMonth}, DateOnly{lastDayPreviousMonth}
}

func FirstDayOfMonth(year int, month time.Month) DateOnly {
	return DateOnly{time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)}
}

func LastDayOfMonth(year int, month time.Month) DateOnly {
	return DateOnly{time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)}
}
