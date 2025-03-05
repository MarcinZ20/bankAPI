package main

import (
	"context"
	"fmt"

	"github.com/MarcinZ20/bankAPI/db"
	"github.com/MarcinZ20/bankAPI/internal/parser"
	"github.com/MarcinZ20/bankAPI/utils"
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
		validationResult := utils.ValidateBankEntity(bank)
		if !validationResult.IsValid {
			fmt.Println(validationResult.Errors)
		}
	}

	fmt.Println("Validation completed")
}
