package blikk

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type queryTestParams struct {
	ListOptions
	Text      string          `query:"filter.text"`
	Number    int             `query:"filter.number"`
	IDs       []int           `query:"filter.ids"`
	Statuses  []InvoiceStatus `query:"filter.statuses"`
	Flag      bool            `query:"filter.flag"`
	TriState  *bool           `query:"filter.triState"`
	Day       Date            `query:"filter.day"`
	Moment    DateTime        `query:"filter.moment"`
	Untagged  string
	Enumerish InvoiceType `query:"filter.enum"`
}

func TestEncodeQueryOmitsZeroValues(t *testing.T) {
	query, err := encodeQuery(queryTestParams{})
	require.NoError(t, err)
	assert.Empty(t, query)
}

func TestEncodeQueryEncodesEverySupportedType(t *testing.T) {
	params := queryTestParams{
		ListOptions: ListOptions{Page: 2, PageSize: 50, SortBy: "name", SortOrder: SortDescending},
		Text:        "roof & walls",
		Number:      42,
		IDs:         []int{1, 2, 3},
		Statuses:    []InvoiceStatus{InvoiceStatusPending, InvoiceStatusDraft},
		Flag:        true,
		TriState:    Ptr(true),
		Day:         NewDate(2026, time.September, 1),
		Moment:      DateTimeOf(time.Date(2021, 1, 1, 14, 0, 0, 0, time.UTC)),
		Untagged:    "ignored",
		Enumerish:   InvoiceTypeFixedPrice,
	}

	query, err := encodeQuery(params)
	require.NoError(t, err)

	assert.Equal(t, url.Values{
		"page":            {"2"},
		"pageSize":        {"50"},
		"sortBy":          {"name"},
		"sortOrder":       {"descending"},
		"filter.text":     {"roof & walls"},
		"filter.number":   {"42"},
		"filter.ids":      {"1,2,3"},
		"filter.statuses": {"Pending,Draft"},
		"filter.flag":     {"true"},
		"filter.triState": {"true"},
		"filter.day":      {"2026-09-01"},
		"filter.moment":   {"2021-01-01 14:00:00"},
		"filter.enum":     {"FixedPrice"},
	}, query)
}

func TestEncodeQuerySendsFalseForTriStatePointer(t *testing.T) {
	query, err := encodeQuery(queryTestParams{TriState: Ptr(false)})
	require.NoError(t, err)
	assert.Equal(t, url.Values{"filter.triState": {"false"}}, query)
}

func TestEncodeQueryOmitsPlainFalseBool(t *testing.T) {
	query, err := encodeQuery(queryTestParams{Flag: false})
	require.NoError(t, err)
	assert.NotContains(t, query, "filter.flag")
}

func TestEncodeQueryOmitsEmptySlice(t *testing.T) {
	query, err := encodeQuery(queryTestParams{IDs: []int{}})
	require.NoError(t, err)
	assert.NotContains(t, query, "filter.ids")
}

func TestEncodeQueryEncodesDateTimeInUTC(t *testing.T) {
	stockholm := time.FixedZone("CEST", 2*60*60)
	query, err := encodeQuery(queryTestParams{Moment: DateTimeOf(time.Date(2021, 1, 1, 16, 0, 0, 0, stockholm))})
	require.NoError(t, err)
	assert.Equal(t, "2021-01-01 14:00:00", query.Get("filter.moment"))
}

func TestEncodeQueryAcceptsPointerAndNil(t *testing.T) {
	query, err := encodeQuery(&queryTestParams{Number: 7})
	require.NoError(t, err)
	assert.Equal(t, "7", query.Get("filter.number"))

	query, err = encodeQuery(nil)
	require.NoError(t, err)
	assert.Empty(t, query)
}

func TestEncodeQueryRejectsUnsupportedTypes(t *testing.T) {
	type unsupported struct {
		Ratio float64 `query:"filter.ratio"`
	}
	_, err := encodeQuery(unsupported{Ratio: 1.5})
	assert.ErrorContains(t, err, "Ratio")

	_, err = encodeQuery("not a struct")
	assert.Error(t, err)
}

func TestEncodeQueryPercentEncodesDateTimeSpace(t *testing.T) {
	query, err := encodeQuery(queryTestParams{Moment: DateTimeOf(time.Date(2021, 1, 1, 14, 0, 0, 0, time.UTC))})
	require.NoError(t, err)
	assert.Equal(t, "filter.moment=2021-01-01+14%3A00%3A00", query.Encode())
}
