package main

import (
	"context"
	"fmt"

	"github.com/MarcinZ20/bankAPI/db"
	"github.com/MarcinZ20/bankAPI/handlers/parser"
	"github.com/MarcinZ20/bankAPI/handlers/validator"
)

func main() {

	client, err := db.ConnectDb()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Disconnect(context.Background())

	bankData, err := parser.ParseBankData(parser.SpreedsheetData)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, bank := range bankData {
		validationResult := validator.ValidateBankEntity(bank)
		if !validationResult.IsValid {
			fmt.Println(validationResult.Errors)
		}
	}

	fmt.Println("Validation completed")
}
