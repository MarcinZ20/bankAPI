package models

// TODO: Check if getters are really necessary
// BankData represents a bank data structure to be processed by the parser
type Headquarter struct {
	Address         string   `bson:"address" json:"address"`
	BankName        string   `bson:"bankName" json:"bankName"`
	CountryISO2Code string   `bson:"countryISO2" json:"countryISO2"`
	CountryName     string   `bson:"countryName" json:"countryName"`
	IsHeadquarter   bool     `bson:"isHeadquarter" json:"isHeadquarter"`
	SwiftCode       string   `bson:"swift_code" json:"swift_code"`
	Branches        []Branch `bson:"branches" json:"branches"`
}

func (h *Headquarter) GetAddress() string {
	return h.Address
}

func (h *Headquarter) GetName() string {
	return h.BankName
}

func (h *Headquarter) GetCountryISO2Code() string {
	return h.CountryISO2Code
}

func (h *Headquarter) GetCountryName() string {
	return h.CountryName
}

func (h *Headquarter) IsHQ() bool {
	return h.IsHeadquarter
}

func (h *Headquarter) GetSwiftCode() string {
	return h.SwiftCode
}

func (h *Headquarter) GetBranches() []Branch {
	return h.Branches
}
