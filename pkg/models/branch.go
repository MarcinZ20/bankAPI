package models

// Branch represents a branch data structure to be processed by the parser
type Branch struct {
	Address       string `bson:"address" json:"address"`
	BankName      string `bson:"bankName" json:"bankName"`
	CountryISO2   string `bson:"countryISO2" json:"countryISO2"`
	CountryName   string `bson:"countryName" json:"countryName"`
	IsHeadquarter bool   `bson:"isHeadquarter" json:"isHeadquarter"`
	SwiftCode     string `bson:"swiftCode" json:"swiftCode"`
}

func (b *Branch) GetAddress() string {
	return b.Address
}

func (b *Branch) GetBankName() string {
	return b.BankName
}

func (b *Branch) GetCountryISO2() string {
	return b.CountryISO2
}

func (b *Branch) GetCountryName() string {
	return b.CountryName
}

func (b *Branch) IsHq() bool {
	return b.IsHeadquarter
}

func (b *Branch) GetSwiftCode() string {
	return b.SwiftCode
}
