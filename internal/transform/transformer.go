package transform

import (
	"strings"

	"github.com/MarcinZ20/bankAPI/pkg/models"
)

// Handles all model transformations
type ModelTransformer struct{}

// Normalizes and sanitizes input data
func (t *ModelTransformer) CleanRequestModel(branch *models.Branch) {
	branch.SwiftCode = strings.ToUpper(branch.SwiftCode)
	branch.CountryISO2 = strings.ToUpper(branch.CountryISO2)
	branch.BankName = strings.ToUpper(branch.BankName)
	branch.CountryName = strings.ToUpper(branch.CountryName)
	branch.Address = strings.Trim(branch.Address, " ")
}

// Transforms Bank entity into Headquarter object
func (t *ModelTransformer) ToHeadquarter(bank models.Bank) models.Headquarter {
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
func (t *ModelTransformer) ToBranch(bank models.Bank) models.Branch {
	return models.Branch{
		Address:       bank.Address,
		BankName:      bank.Name,
		CountryISO2:   bank.CountryISO2Code,
		CountryName:   bank.CountryName,
		IsHeadquarter: bank.IsHeadquarter(),
		SwiftCode:     bank.SwiftCode,
	}
}

// Transforms raw bank data into database-ready format
func (t *ModelTransformer) TransformBankData(banks *[]models.Bank) *map[string]models.Headquarter {
	hqs := make(map[string]models.Headquarter)
	brs := make(map[string][]models.Branch)

	t.convertIntoMaps(banks, hqs, brs)
	t.mergeMaps(hqs, brs)

	return &hqs
}

// Converts a list of Bank models to headquarter and branch mappings
func (t *ModelTransformer) convertIntoMaps(banks *[]models.Bank, hqs map[string]models.Headquarter, brs map[string][]models.Branch) {
	for _, bank := range *banks {
		keyCode := bank.SwiftCode[0:8]

		if bank.IsHeadquarter() {
			hqs[keyCode] = t.ToHeadquarter(bank)
		} else {
			brs[keyCode] = append(brs[keyCode], t.ToBranch(bank))
		}
	}
}

// Merges branches into corresponding headquarters
func (t *ModelTransformer) mergeMaps(hqs map[string]models.Headquarter, brs map[string][]models.Branch) {
	for keyCode := range hqs {
		if branches, ok := brs[keyCode]; ok {
			hq := hqs[keyCode]
			hq.Branches = append(hq.Branches, branches...)
			hqs[keyCode] = hq
		}
	}
}
