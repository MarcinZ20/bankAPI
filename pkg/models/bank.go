package models

import "strings"

// This represents an object parsed from google spreadsheet
type Bank struct {
	CountryISO2Code string `json:"countryISO2Code"`
	SwiftCode       string `json:"swiftCode"`
	CodeType        string `json:"codeType"`
	Name            string `json:"name"`
	Address         string `json:"address"`
	TownName        string `json:"townName"`
	CountryName     string `json:"countryName"`
	Timezone        string `json:"timezone"`
}

func (b *Bank) IsHeadquarter() bool {
	return strings.HasSuffix(b.SwiftCode, "XXX")
}

func (b *Bank) GetName() string {
	return b.Name
}

func (b *Bank) GetCountryISO2Code() string {
	return b.CountryISO2Code
}

func (b *Bank) GetCountryName() string {
	return b.CountryName
}

func (b *Bank) GetSwiftCode() string {
	return b.SwiftCode
}
