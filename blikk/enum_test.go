package blikk

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnumsNormaliseCaseAndPassUnknownValuesThrough(t *testing.T) {
	var decoded struct {
		License       License               `json:"license"`
		Lowercase     InvoiceType           `json:"lowercase"`
		Camel         InvoiceType           `json:"camel"`
		ProjectType   TimeReportProjectType `json:"projectType"`
		InvoiceStatus InvoiceStatus         `json:"invoiceStatus"`
		Unknown       License               `json:"unknown"`
		Null          InvoiceStatus         `json:"null"`
	}
	input := `{"license":"total","lowercase":"ongoing","camel":"fixedPrice","projectType":"internalproject","invoiceStatus":"INVOICED","unknown":"Enterprise","null":null}`

	require.NoError(t, json.Unmarshal([]byte(input), &decoded))

	assert.Equal(t, LicenseTotal, decoded.License)
	assert.Equal(t, InvoiceTypeOngoing, decoded.Lowercase)
	assert.Equal(t, InvoiceTypeFixedPrice, decoded.Camel)
	assert.Equal(t, TimeReportProjectTypeInternalProject, decoded.ProjectType)
	assert.Equal(t, InvoiceStatusInvoiced, decoded.InvoiceStatus)
	assert.Equal(t, License("Enterprise"), decoded.Unknown)
	assert.Equal(t, InvoiceStatus(""), decoded.Null)
}

func TestEnumKeepsNumberAsDecimalString(t *testing.T) {
	var decoded struct {
		License       License       `json:"license"`
		InvoiceStatus InvoiceStatus `json:"invoiceStatus"`
	}

	require.NoError(t, json.Unmarshal([]byte(`{"license":42,"invoiceStatus":3}`), &decoded))

	assert.Equal(t, License("42"), decoded.License)
	assert.Equal(t, InvoiceStatus("3"), decoded.InvoiceStatus)
}

func TestEnumRejectsNonStringNonNumber(t *testing.T) {
	var license License
	assert.Error(t, json.Unmarshal([]byte(`true`), &license))
	assert.Error(t, json.Unmarshal([]byte(`{"name":"Total"}`), &license))
}
