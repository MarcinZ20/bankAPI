package models

// TODO: Check if getters are really necessary
// Branch represents a branch data structure to be processed by the parser
type Branch struct {
	Address         string `bson:"address" json:"address"`
	Name            string `bson:"bankName" json:"bankName"`
	CountryISO2Code string `bson:"countryISO2" json:"countryISO2"`
	CountryName     string `bson:"countryName" json:"countryName"`
	IsHeadquarter   bool   `bson:"isHeadquarter" json:"isHeadquarter"`
	SwiftCode       string `bson:"swiftCode" json:"swiftCode"`
}

func (b *Branch) GetAddress() string {
	return b.Address
}

func (b *Branch) GetName() string {
	return b.Name
}

func (b *Branch) GetCountryISO2Code() string {
	return b.CountryISO2Code
}

func (b *Branch) GetCountryName() string {
	return b.CountryName
}

func (b *Branch) IsHQ() bool {
	return b.IsHeadquarter
}

func (b *Branch) GetSwiftCode() string {
	return b.SwiftCode
}
