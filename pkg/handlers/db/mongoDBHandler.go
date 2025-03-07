package db

import (
	"context"
	"fmt"
	"os"

	"github.com/MarcinZ20/bankAPI/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveEntities(collection *mongo.Collection, data map[string]models.Headquarter) {
	var docs []any

	for _, entity := range data {
		docs = append(docs, entity)
	}

	if len(docs) > 0 {
		_, err := collection.InsertMany(context.Background(), docs)
		if err != nil {
			fmt.Printf("Failed to insert records: %v", err)
		} else {
			fmt.Printf("Inserted %d records into the database\n", len(docs))
		}
	}
}

func GetCollection(client *mongo.Client) (*mongo.Collection, error) {
	if client == nil {
		return nil, fmt.Errorf("client is not working")
	}

	col := client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION"))

	return col, nil
}
