package transform

import (
	"testing"

	"github.com/MarcinZ20/bankAPI/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestModelTransformer_CleanRequestModel(t *testing.T) {
	transformer := &ModelTransformer{}

	tests := []struct {
		name     string
		input    *models.Branch
		expected *models.Branch
	}{
		{
			name: "Clean mixed case input",
			input: &models.Branch{
				SwiftCode:     "deutdeff",
				CountryISO2:   "de",
				BankName:      "Deutsche Bank",
				CountryName:   "Germany",
				Address:       "  Sample Street 123  ",
				IsHeadquarter: false,
			},
			expected: &models.Branch{
				SwiftCode:     "DEUTDEFF",
				CountryISO2:   "DE",
				BankName:      "DEUTSCHE BANK",
				CountryName:   "GERMANY",
				Address:       "Sample Street 123",
				IsHeadquarter: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer.CleanRequestModel(tt.input)
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}

func TestModelTransformer_ToHeadquarter(t *testing.T) {
	transformer := &ModelTransformer{}

	tests := []struct {
		name     string
		input    models.Bank
		expected models.Headquarter
	}{
		{
			name: "Transform bank to headquarter",
			input: models.Bank{
				SwiftCode:       "DEUTDEFFXXX",
				CountryISO2Code: "DE",
				Name:            "Deutsche Bank",
				Address:         "Frankfurt",
				CountryName:     "Germany",
			},
			expected: models.Headquarter{
				SwiftCode:     "DEUTDEFFXXX",
				CountryISO2:   "DE",
				BankName:      "Deutsche Bank",
				Address:       "Frankfurt",
				CountryName:   "Germany",
				IsHeadquarter: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformer.ToHeadquarter(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestModelTransformer_ToBranch(t *testing.T) {
	transformer := &ModelTransformer{}

	tests := []struct {
		name     string
		input    models.Bank
		expected models.Branch
	}{
		{
			name: "Transform bank to branch",
			input: models.Bank{
				SwiftCode:       "DEUTDEFF100",
				CountryISO2Code: "DE",
				Name:            "Deutsche Bank Branch",
				Address:         "Berlin",
				CountryName:     "Germany",
			},
			expected: models.Branch{
				SwiftCode:     "DEUTDEFF100",
				CountryISO2:   "DE",
				BankName:      "Deutsche Bank Branch",
				Address:       "Berlin",
				CountryName:   "Germany",
				IsHeadquarter: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformer.ToBranch(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestModelTransformer_TransformBankData(t *testing.T) {
	transformer := &ModelTransformer{}

	tests := []struct {
		name             string
		input            []models.Bank
		expectedHQCount  int
		expectedBranches map[string]int
	}{
		{
			name: "Transform banks with HQ and branches",
			input: []models.Bank{
				{
					SwiftCode:       "DEUTDEFFXXX",
					CountryISO2Code: "DE",
					Name:            "Deutsche Bank HQ",
					Address:         "Frankfurt",
					CountryName:     "Germany",
				},
				{
					SwiftCode:       "DEUTDEFF111",
					CountryISO2Code: "DE",
					Name:            "Deutsche Bank Branch 1",
					Address:         "Berlin",
					CountryName:     "Germany",
				},
				{
					SwiftCode:       "DEUTDEFF222",
					CountryISO2Code: "DE",
					Name:            "Deutsche Bank Branch 2",
					Address:         "Munich",
					CountryName:     "Germany",
				},
			},
			expectedHQCount: 1,
			expectedBranches: map[string]int{
				"DEUTDEFF": 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformer.TransformBankData(&tt.input)

			// number of hq
			assert.Equal(t, tt.expectedHQCount, len(*result))

			// branches for each hq
			for hqCode, expectedBranchCount := range tt.expectedBranches {
				hq, exists := (*result)[hqCode]
				assert.True(t, exists, "Headquarter should exist")
				assert.Equal(t, expectedBranchCount, len(hq.Branches),
					"Wrong number of branches for the given hq %s", hqCode)
			}
		})
	}
}
