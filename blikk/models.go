package blikk

type NameID struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProjectRef struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number string `json:"number"`
}

type CostCenterRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type TagRef struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

type Address struct {
	StreetAddress     string `json:"streetAddress"`
	AdditionalAddress string `json:"additionalAddress"`
	PostalCode        string `json:"postalCode"`
	City              string `json:"city"`
	State             string `json:"state"`
	CountryID         int    `json:"countryId"`
	CountryName       string `json:"countryName"`
}

type Location struct {
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	StreetAddress string  `json:"streetAddress"`
	PostalCode    string  `json:"postalCode"`
	City          string  `json:"city"`
	CountryName   string  `json:"countryName"`
}

type ProjectStatusRef struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	IsCompletedStatus bool   `json:"isCompletedStatus"`
}

type ProjectCategory struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Tag struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Color       string    `json:"color"`
	CreatedDate DateTime  `json:"createdDate"`
	CreatedBy   *NameID   `json:"createdBy"`
	UpdatedDate *DateTime `json:"updatedDate"`
	UpdatedBy   *NameID   `json:"updatedBy"`
}

type Comment struct {
	ID            int       `json:"id"`
	Text          string    `json:"text"`
	CreatedByUser *NameID   `json:"createdByUser"`
	UpdatedByUser *NameID   `json:"updatedByUser"`
	CreatedDate   *DateTime `json:"createdDate"`
	UpdatedDate   *DateTime `json:"updatedDate"`
}

type VatCode struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	VAT       float64   `json:"vat"`
	IsVatFree bool      `json:"isVatFree"`
	CreatedBy *NameID   `json:"createdBy"`
	UpdatedBy *NameID   `json:"updatedBy"`
	Created   *DateTime `json:"created"`
	Updated   *DateTime `json:"updated"`
}

type commentInput struct {
	Text string `json:"text"`
}

func Ptr[T any](v T) *T { return &v }
