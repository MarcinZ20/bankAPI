package transformer

import (
	"strings"

	"github.com/MarcinZ20/bankAPI/handlers/parser"
)

// IsBranchOf checks if a bank is a branch of a headquarters
func IsBranchOf(bank parser.Bank, hq Headquarter) bool {
	return bank.SwiftCode[:8] == hq.SwiftCode[:8] && bank.CountryISO2Code == hq.CountryISO2Code
}

func TransformBankEntity(bank *parser.Bank) {
	bank.SwiftCode = strings.ToUpper(bank.SwiftCode)
	bank.CountryISO2Code = strings.ToUpper(bank.CountryISO2Code)
	bank.Name = strings.ToUpper(bank.Name)
	bank.CountryName = strings.ToUpper(bank.CountryName)
}

func TransformIntoHeadquarter(bank parser.Bank) Headquarter {
	return Headquarter{
		Address:         bank.Address,
		BankName:        bank.Name,
		CountryISO2Code: bank.CountryISO2Code,
		CountryName:     bank.CountryName,
		IsHeadquarter:   bank.IsHeadquarter(),
		SwiftCode:       bank.SwiftCode,
	}
}

func TransformIntoBranch(bank parser.Bank) Branch {
	return Branch{
		Address:         bank.Address,
		Name:            bank.Name,
		CountryISO2Code: bank.CountryISO2Code,
		CountryName:     bank.CountryName,
		IsHeadquarter:   bank.IsHeadquarter(),
		SwiftCode:       bank.SwiftCode,
	}
}

func Transform(banks []parser.Bank) map[string]Headquarter {
	hqs := make(map[string]Headquarter)
	brs := make(map[string][]Branch)

	ConvertIntoMaps(banks, hqs, brs)
	MergeMaps(hqs, brs)

	return hqs
}

func ConvertIntoMaps(banks []parser.Bank, hqs map[string]Headquarter, brs map[string][]Branch) {
	for _, bank := range banks {
		key_code := bank.SwiftCode[0:8]

		if bank.IsHeadquarter() {
			hqs[key_code] = TransformIntoHeadquarter(bank)
		} else {
			brs[key_code] = append(brs[key_code], TransformIntoBranch(bank))
		}
	}
}

// TODO: Refactor using pointers
func MergeMaps(hqs map[string]Headquarter, brs map[string][]Branch) {
	for key_code := range hqs {
		if branches, ok := brs[key_code]; ok {
			hq := hqs[key_code]
			hq.Branches = append(hq.Branches, branches...)
			hqs[key_code] = hq
		}
	}
}
