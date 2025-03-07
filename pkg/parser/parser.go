package parser

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/MarcinZ20/bankAPI/pkg/models"
)

func ParseBankData(response string, data *[]models.Bank) error {
	reader := csv.NewReader(strings.NewReader(response))
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error while parsing data from .csv: %v", err)
	}

	for i, row := range rows {
		if i == 0 { // Skip header row
			continue
		}
		bank := models.Bank{
			CountryISO2Code: row[0],
			SwiftCode:       row[1],
			Name:            row[3],
			Address:         row[4],
			CountryName:     row[6],
		}
		*data = append(*data, bank)
	}

	return nil
}
