package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MarcinZ20/bankAPI/handlers/transformer"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDb() (*mongo.Client, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}

	mongoUri := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoUri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}

func SaveEntities(collection *mongo.Collection, data map[string]transformer.Headquarter) {
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
