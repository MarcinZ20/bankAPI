package importer

import (
	"context"
	"fmt"
	"time"

	"github.com/MarcinZ20/bankAPI/internal/database"
	"github.com/MarcinZ20/bankAPI/internal/parser"
	"github.com/MarcinZ20/bankAPI/internal/spreadsheet"
	"github.com/MarcinZ20/bankAPI/internal/transform"
	"github.com/MarcinZ20/bankAPI/internal/validation"
	"github.com/MarcinZ20/bankAPI/pkg/models"
)

// Handles the data import process from a Google Spreadsheet
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

	response, err := spreadsheet.FetchData(googleSpreadsheet)
	if err != nil {
		return fmt.Errorf("failed to fetch spreadsheet data: %w", err)
	}

	var rawData []models.Bank
	if err := parser.ParseBankData(response, &rawData); err != nil {
		return fmt.Errorf("failed to parse bank data: %w", err)
	}

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

	transformer := transform.ModelTransformer{}
	transformedData := transformer.TransformBankData(&rawData)

	// For clean setup, clean existing data
	if err := db.Collection.Drop(ctx); err != nil {
		return fmt.Errorf("failed to clear existing data: %w", err)
	}

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
