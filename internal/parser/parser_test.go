package parser

import (
	"testing"

	"github.com/MarcinZ20/bankAPI/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestParser_ParseBankData(t *testing.T) {
	parser := NewParser()
	parser.minCols = 7
	parser.skipHeaderRow = true

	tests := []struct {
		name        string
		csvData     string
		wantErr     bool
		expectedLen int
		validate    func(t *testing.T, banks []models.Bank)
	}{
		{
			name: "Valid CSV data",
			csvData: `CountryISO2,SwiftCode,Type,Name,Address,City,Country
DE,DEUTDEFF,HQ,Deutsche Bank,Taunusanlage 12,Frankfurt,Germany
FR,BNPAFRPP,HQ,BNP Paribas,16 Boulevard des Italiens,Paris,France`,
			wantErr:     false,
			expectedLen: 2,
			validate: func(t *testing.T, banks []models.Bank) {

				assert.Equal(t, "DE", banks[0].CountryISO2Code)
				assert.Equal(t, "DEUTDEFF", banks[0].SwiftCode)
				assert.Equal(t, "Deutsche Bank", banks[0].Name)
				assert.Equal(t, "Taunusanlage 12", banks[0].Address)
				assert.Equal(t, "Germany", banks[0].CountryName)

				assert.Equal(t, "FR", banks[1].CountryISO2Code)
				assert.Equal(t, "BNPAFRPP", banks[1].SwiftCode)
				assert.Equal(t, "BNP Paribas", banks[1].Name)
				assert.Equal(t, "16 Boulevard des Italiens", banks[1].Address)
				assert.Equal(t, "France", banks[1].CountryName)
			},
		},
		{
			name: "Empty CSV data",
			csvData: `CountryISO2,SwiftCode,Type,Name,Address,City,Country
`,
			wantErr:     false,
			expectedLen: 0,
		},
		{
			name:        "Invalid CSV format",
			csvData:     "invalid,csv,data",
			wantErr:     true,
			expectedLen: 0,
		},
		{
			name:        "Empty string",
			csvData:     "",
			wantErr:     true,
			expectedLen: 0,
		},
		{
			name: "Row with insufficient columns",
			csvData: `CountryISO2,SwiftCode,Type,Name,Address,City,Country
DE,DEUTDEFF,HQ,Deutsche Bank,Taunusanlage 12`,
			wantErr:     true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var banks []models.Bank
			err := parser.ParseBankData(tt.csvData, &banks)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, banks, tt.expectedLen)

			if tt.validate != nil {
				tt.validate(t, banks)
			}
		})
	}
}

func TestParser_ParseBankData_NilSlice(t *testing.T) {
	parser := NewParser()
	err := parser.ParseBankData("some data", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil slice")
}
