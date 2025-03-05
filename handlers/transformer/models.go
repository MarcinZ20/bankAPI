package transformer

// Bank defines the interface for bank-related data structures
// It provides methods to access common bank information like address, name, country details, etc.
type Bank interface {
	GetAddress() string
	GetName() string
	GetCountryISO2Code() string
	GetCountryName() string
	IsHeadquarter() bool
	GetSwiftCode() string
}

// Branch represents a branch data structure to be processed by the parser
type Branch struct {
	Address         string `bson:"address" json:"address"`
	Name            string `bson:"bankName" json:"bankName"`
	CountryISO2Code string `bson:"countryISO2" json:"countryISO2"`
	CountryName     string `bson:"countryName" json:"countryName"`
	IsHeadquarter   bool   `bson:"isHeadquarter" json:"isHeadquarter"`
	SwiftCode       string `bson:"swift_code" json:"swift_code"`
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
