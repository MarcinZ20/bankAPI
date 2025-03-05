package transformer

import (
	"strings"

	"github.com/MarcinZ20/bankAPI/internal/parser"
)

func TransformBankEntity(bank parser.Bank) parser.Bank {
	bank.SwiftCode = strings.ToUpper(bank.SwiftCode)
	bank.CountryISO2Code = strings.ToUpper(bank.CountryISO2Code)
	bank.Name = strings.ToUpper(bank.Name)
	bank.CountryName = strings.ToUpper(bank.CountryName)

	return bank
}
