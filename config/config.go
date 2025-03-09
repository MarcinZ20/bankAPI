package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MarcinZ20/bankAPI/api/middleware"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client
var ApiServer *fiber.App

func ConnectDb() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoUri := os.Getenv("MONGO_URI")
	clientOptions := options.Client().SetTimeout(3 * time.Second).ApplyURI(mongoUri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Printf("error while connecting to mongoDB: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Printf("error while pinging the database: %v", err)
	}

	MongoClient = client
	fmt.Println("Connected succesfuly to mongoDB!")
}

func ConfigAPI() {
	config := fiber.Config{
		AppName:       "bankAPI v1.0",
		CaseSensitive: true,
		StrictRouting: true,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		BodyLimit:     10 * 1024 * 1024, // 10MB
		ReadTimeout:   5 * time.Second,
	}

	ApiServer = fiber.New(config)
	ApiServer.Use(middleware.WithTimeout(5 * time.Second))
}
