package main

import (
	"context"
	"fmt"
	"os"

	"github.com/MarcinZ20/bankAPI/config"
	"github.com/MarcinZ20/bankAPI/pkg/handlers/spreadsheet"
	"github.com/MarcinZ20/bankAPI/pkg/models"
	"github.com/MarcinZ20/bankAPI/pkg/parser"
	"github.com/MarcinZ20/bankAPI/pkg/validation"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Print(fmt.Errorf("error while loading .env: %v", err))
	}

	client, err := config.ConnectDb()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Disconnect(context.Background())

	// collection := client.Database(os.Getenv("MONGO_DB")).Collection(os.Getenv("MONGO_COLLECTION"))
	mySpreadsheet := models.GoogleSpreadsheet{
		SpreadsheetId: os.Getenv("SPREADSHEET_ID"),
	}

	response, err := spreadsheet.FetchSpreadsheetData(&mySpreadsheet)
	if err != nil {
		fmt.Println("Error")
	}

	var data []models.Bank
	err = parser.ParseBankData(response, &data)
	if err != nil {
		fmt.Println("Error")
	}

	for _, bank := range data {
		validationResult := validation.ValidateBankEntity(bank)
		if !validationResult.IsValid {
			fmt.Println(validationResult.Errors)
		}
	}
	fmt.Println("Validation completed")

	// data := transformer.Transform(bankData)

	// db.SaveEntities(collection, data)
}
