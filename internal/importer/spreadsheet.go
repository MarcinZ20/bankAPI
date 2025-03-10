package importer

import (
	"context"
	"fmt"
	"time"

	"github.com/MarcinZ20/bankAPI/internal/database"
	"github.com/MarcinZ20/bankAPI/internal/spreadsheet"
	"github.com/MarcinZ20/bankAPI/pkg/models"
	"github.com/MarcinZ20/bankAPI/internal/parser"
	"github.com/MarcinZ20/bankAPI/internal/transform"
	"github.com/MarcinZ20/bankAPI/internal/validation"
)

// ImportSpreadsheetData handles the data import process from a Google Spreadsheet
func ImportSpreadsheetData(ctx context.Context, spreadsheetID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	db := database.GetInstance()
	if db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	googleSpreadsheet := &models.GoogleSpreadsheet{
		SpreadsheetId: spreadsheetID,
	}

	// Fetch data from spreadsheet
	response, err := spreadsheet.FetchData(googleSpreadsheet)
	if err != nil {
		return fmt.Errorf("failed to fetch spreadsheet data: %w", err)
	}

	// Parse raw data
	var rawData []models.Bank
	if err := parser.ParseBankData(response, &rawData); err != nil {
		return fmt.Errorf("failed to parse bank data: %w", err)
	}

	// Validate data
	var validationErrors []error
	for i, bank := range rawData {
		result := validation.ValidateBankEntity(bank)
		if !result.IsValid {
			validationErrors = append(validationErrors,
				fmt.Errorf("validation failed for bank at index %d: %v", i, result.Errors))
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("validation errors occurred: %v", validationErrors)
	}

	// Transform data
	transformedData := transform.Transform(&rawData)

	// Clear existing data
	if err := db.Collection.Drop(ctx); err != nil {
		return fmt.Errorf("failed to clear existing data: %w", err)
	}

	// Insert new data
	var documents []any
	for _, bank := range *transformedData {
		documents = append(documents, bank)
	}

	_, err = db.Collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	return nil
}
