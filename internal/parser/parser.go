package parser

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/MarcinZ20/bankAPI/pkg/models"
)

// Parser handles all data parsing operations
type Parser struct {
	skipHeaderRow bool
	minCols       int
}

// Creates a new parser instance
func NewParser() *Parser {
	return &Parser{
		skipHeaderRow: true,
	}
}

// ParseBankData parses bank data from CSV string into Bank objects
func (p *Parser) ParseBankData(response string, data *[]models.Bank) error {
	if data == nil {
		return fmt.Errorf("nil slice: data parameter cannot be nil")
	}

	if response == "" {
		return fmt.Errorf("empty input: response cannot be empty")
	}

	reader := csv.NewReader(strings.NewReader(response))
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error while parsing data from .csv: %v", err)
	}

	if len(rows) == 0 {
		return fmt.Errorf("invalid CSV: no data found")
	}

	// Validate header row
	if len(rows[0]) < p.minCols {
		return fmt.Errorf("invalid CSV format: expected at least 7 columns in header")
	}

	expectedColumns := len(rows[0])

	for i, row := range rows {
		if i == 0 && p.skipHeaderRow {
			continue
		}

		if len(row) != expectedColumns {
			return fmt.Errorf("invalid CSV format at line %d: expected %d columns, got %d", i+1, expectedColumns, len(row))
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
