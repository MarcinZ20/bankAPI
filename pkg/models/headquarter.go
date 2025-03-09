package models

// TODO: Check if getters are really necessary
// BankData represents a bank data structure to be processed by the parser
type Headquarter struct {
	Address       string   `bson:"address" json:"address"`
	BankName      string   `bson:"bankName" json:"bankName"`
	CountryISO2   string   `bson:"countryISO2" json:"countryISO2"`
	CountryName   string   `bson:"countryName" json:"countryName"`
	IsHeadquarter bool     `bson:"isHeadquarter" json:"isHeadquarter"`
	SwiftCode     string   `bson:"swiftCode" json:"swiftCode"`
	Branches      []Branch `bson:"branches" json:"branches"`
}

func (h *Headquarter) GetAddress() string {
	return h.Address
}

func (h *Headquarter) GetBankName() string {
	return h.BankName
}

func (h *Headquarter) GetCountryISO2() string {
	return h.CountryISO2
}

func (h *Headquarter) GetCountryName() string {
	return h.CountryName
}

func (h *Headquarter) IsHq() bool {
	return h.IsHeadquarter
}

func (h *Headquarter) GetSwiftCode() string {
	return h.SwiftCode
}

func (h *Headquarter) GetBranches() []Branch {
	return h.Branches
}
