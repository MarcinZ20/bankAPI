package main

import (
	"context"
	"fmt"
	"os"

	"github.com/MarcinZ20/bankAPI/db"
	"github.com/MarcinZ20/bankAPI/handlers/parser"
	"github.com/MarcinZ20/bankAPI/handlers/transformer"
	"github.com/MarcinZ20/bankAPI/handlers/validator"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	client, err := db.ConnectDb()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(os.Getenv("MONGO_DB")).Collection(os.Getenv("MONGO_COLLECTION"))

	bankData, err := parser.ParseBankData(parser.SpreedsheetData)
	if err != nil {
		fmt.Println(err) // TODO: make this handle errors the propper way
		return
	}

	for _, bank := range bankData {
		validationResult := validator.ValidateBankEntity(bank)
		if !validationResult.IsValid {
			fmt.Println(validationResult.Errors)
		}
	}
	fmt.Println("Validation completed")

	data := transformer.Transform(bankData)

	for _, d := range data {
		fmt.Printf("Bank: %v \tBranches: %v\n", d.BankName, len(d.Branches))
	}

	db.SaveEntities(collection, data)
}
