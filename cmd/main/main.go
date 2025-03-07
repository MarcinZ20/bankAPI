package main

import (
	"context"
	"fmt"
	"os"

	"github.com/MarcinZ20/bankAPI/api/routes"
	"github.com/MarcinZ20/bankAPI/config"
	"github.com/MarcinZ20/bankAPI/pkg/handlers/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Print(fmt.Errorf("error while loading .env: %v", err))
	}

	config.ConnectDb()
	defer config.MongoClient.Disconnect(context.Background())

	collection, err := db.GetCollection(config.MongoClient)
	if err != nil {
		fmt.Println(err)
		return
	}

	indexModel := mongo.IndexModel{
		Keys: bson.D{{
			Key: "swiftCode", Value: 1,
		}},
	}

	name, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		fmt.Printf("Couldnt create collection %v: %v", name, err)
	}

	config.ConfigAPI()
	routes.BankRoutes(config.ApiServer)

	// mySpreadsheet := models.GoogleSpreadsheet{
	// 	SpreadsheetId: os.Getenv("SPREADSHEET_ID"),
	// }

	// response, err := spreadsheet.FetchSpreadsheetData(&mySpreadsheet)
	// if err != nil {
	// 	fmt.Println("Error")
	// }

	// var rawData []models.Bank
	// err = parser.ParseBankData(response, &rawData)
	// if err != nil {
	// 	fmt.Println("Error")
	// }

	// for _, bank := range rawData {
	// 	validationResult := validation.ValidateBankEntity(bank)
	// 	if !validationResult.IsValid {
	// 		fmt.Println(validationResult.Errors)
	// 	}
	// }

	// transformedData := transform.Transform(&rawData)

	// db.SaveEntities(collection, transformedData)

	if err := config.ApiServer.Listen(os.Getenv("API_SERVER_PORT")); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
