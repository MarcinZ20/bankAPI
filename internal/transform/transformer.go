package transform

import (
	"strings"

	"github.com/MarcinZ20/bankAPI/pkg/models"
)

// IsBranchOf checks if a bank is a branch of a headquarters
func IsBranchOf(bank models.Bank, hq models.Headquarter) bool {
	return bank.SwiftCode[:8] == hq.SwiftCode[:8] && bank.CountryISO2Code == hq.CountryISO2
}

func TransformBankEntity(bank *models.Bank) {
	bank.SwiftCode = strings.ToUpper(bank.SwiftCode)
	bank.CountryISO2Code = strings.ToUpper(bank.CountryISO2Code)
	bank.Name = strings.ToUpper(bank.Name)
	bank.CountryName = strings.ToUpper(bank.CountryName)
}

func TransformRequestModel(branch *models.Branch) {
	branch.SwiftCode = strings.ToUpper(branch.SwiftCode)
	branch.CountryISO2 = strings.ToUpper(branch.CountryISO2)
	branch.BankName = strings.ToUpper(branch.BankName)
	branch.CountryName = strings.ToUpper(branch.CountryName)
}

func TransformIntoHeadquarter(bank models.Bank) models.Headquarter {
	return models.Headquarter{
		Address:       bank.Address,
		BankName:      bank.Name,
		CountryISO2:   bank.CountryISO2Code,
		CountryName:   bank.CountryName,
		IsHeadquarter: bank.IsHeadquarter(),
		SwiftCode:     bank.SwiftCode,
	}
}

func TransformIntoBranch(bank models.Bank) models.Branch {
	return models.Branch{
		Address:       bank.Address,
		BankName:      bank.Name,
		CountryISO2:   bank.CountryISO2Code,
		CountryName:   bank.CountryName,
		IsHeadquarter: bank.IsHeadquarter(),
		SwiftCode:     bank.SwiftCode,
	}
}

func Transform(banks *[]models.Bank) *map[string]models.Headquarter {
	hqs := make(map[string]models.Headquarter)
	brs := make(map[string][]models.Branch)

	ConvertIntoMaps(banks, hqs, brs)
	MergeMaps(hqs, brs)

	return &hqs
}

func ConvertIntoMaps(banks *[]models.Bank, hqs map[string]models.Headquarter, brs map[string][]models.Branch) {
	for _, bank := range *banks {
		key_code := bank.SwiftCode[0:8]

		if bank.IsHeadquarter() {
			hqs[key_code] = TransformIntoHeadquarter(bank)
		} else {
			brs[key_code] = append(brs[key_code], TransformIntoBranch(bank))
		}
	}
}

// TODO: Refactor using pointers
func MergeMaps(hqs map[string]models.Headquarter, brs map[string][]models.Branch) {
	for key_code := range hqs {
		if branches, ok := brs[key_code]; ok {
			hq := hqs[key_code]
			hq.Branches = append(hq.Branches, branches...)
			hqs[key_code] = hq
		}
	}
}
