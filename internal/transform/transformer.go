package transform

import (
	"strings"

	"github.com/MarcinZ20/bankAPI/pkg/models"
)

// Checks if a bank is a branch of a headquarters
func IsBranchOf(bank models.Bank, hq models.Headquarter) bool {
	return bank.SwiftCode[:8] == hq.SwiftCode[:8] && bank.CountryISO2Code == hq.CountryISO2
}

// Cleans request model properties
func TransformRequestModel(branch *models.Branch) {
	branch.SwiftCode = strings.ToUpper(branch.SwiftCode)
	branch.CountryISO2 = strings.ToUpper(branch.CountryISO2)
	branch.BankName = strings.ToUpper(branch.BankName)
	branch.CountryName = strings.ToUpper(branch.CountryName)
	branch.Address = strings.Trim(branch.Address, " ")
}

// Transforms Bank entity into Headquarter object
func transformIntoHeadquarter(bank models.Bank) models.Headquarter {
	return models.Headquarter{
		Address:       bank.Address,
		BankName:      bank.Name,
		CountryISO2:   bank.CountryISO2Code,
		CountryName:   bank.CountryName,
		IsHeadquarter: bank.IsHeadquarter(),
		SwiftCode:     bank.SwiftCode,
	}
}

// Transforms Bank entity into Branch object
func transformIntoBranch(bank models.Bank) models.Branch {
	return models.Branch{
		Address:       bank.Address,
		BankName:      bank.Name,
		CountryISO2:   bank.CountryISO2Code,
		CountryName:   bank.CountryName,
		IsHeadquarter: bank.IsHeadquarter(),
		SwiftCode:     bank.SwiftCode,
	}
}

// Transforms existing raw data into format ready for database population
func Transform(banks *[]models.Bank) *map[string]models.Headquarter {
	hqs := make(map[string]models.Headquarter)
	brs := make(map[string][]models.Branch)

	convertIntoMaps(banks, hqs, brs)
	mergeMaps(hqs, brs)

	return &hqs
}

// Converts a list of Bank models to headquarter and branch mappings
func convertIntoMaps(banks *[]models.Bank, hqs map[string]models.Headquarter, brs map[string][]models.Branch) {
	for _, bank := range *banks {
		key_code := bank.SwiftCode[0:8]

		if bank.IsHeadquarter() {
			hqs[key_code] = transformIntoHeadquarter(bank)
		} else {
			brs[key_code] = append(brs[key_code], transformIntoBranch(bank))
		}
	}
}

// Merges branches into corresponding headquarters
func mergeMaps(hqs map[string]models.Headquarter, brs map[string][]models.Branch) {
	for key_code := range hqs {
		if branches, ok := brs[key_code]; ok {
			hq := hqs[key_code]
			hq.Branches = append(hq.Branches, branches...)
			hqs[key_code] = hq
		}
	}
}
