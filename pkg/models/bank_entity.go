package models

type BankEntity interface {
	GetAddress() string
	GetBankName() string
	GetCountryISO2() string
	GetCountryName() string
	GetSwiftCode() string
	IsHq() bool
}
