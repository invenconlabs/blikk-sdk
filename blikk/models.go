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

type Tag struct {
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

type ProjectStatus struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	IsCompletedStatus bool   `json:"isCompletedStatus"`
}

type ProjectCategory struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func Ptr[T any](v T) *T { return &v }
