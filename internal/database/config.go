package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"slices"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config holds database configuration and connection details
type Config struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

var instance *Config

// Connect establishes database connection and initializes indexes
func Connect(ctx context.Context) (*Config, error) {
	if instance != nil {
		return instance, nil
	}

	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		return nil, fmt.Errorf("MONGO_URI environment variable is not set")
	}

	connCtx, connCancel := context.WithTimeout(ctx, 5*time.Second)
	defer connCancel()

	clientOptions := options.Client().SetTimeout(3 * time.Second).ApplyURI(mongoUri)
	client, err := mongo.Connect(connCtx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to mongoDB: %w", err)
	}

	defer func() {
		if err != nil {
			client.Disconnect(context.Background())
		}
	}()

	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("error while pinging the database: %w", err)
	}

	dbName := os.Getenv("MONGO_DATABASE")
	collName := os.Getenv("MONGO_COLLECTION")
	if dbName == "" || collName == "" {
		return nil, fmt.Errorf("MONGO_DATABASE and MONGO_COLLECTION environment variables must be set")
	}

	collectionFilter := bson.D{{}}

	collections, err := client.Database(dbName).ListCollectionNames(ctx, collectionFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections")
	}

	if !slices.Contains(collections, collName) {
		if err := createCollection(client.Database(dbName), collName); err != nil {
			return nil, fmt.Errorf("there was no existing collection named %v", collName)
		}
	}

	collection := client.Database(dbName).Collection(collName)

	if err := createIndexes(ctx, collection); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	instance = &Config{
		Client:     client,
		Collection: collection,
	}

	return instance, nil
}

// createIndexes ensures all required indexes exist
func createIndexes(ctx context.Context, collection *mongo.Collection) error {
	indexCtx, indexCancel := context.WithTimeout(ctx, 10*time.Second)
	defer indexCancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "swiftCode", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("swiftCode_unique"),
		},
		{
			Keys:    bson.D{{Key: "countryISO2", Value: 1}},
			Options: options.Index().SetUnique(false).SetName("countryISO2"),
		},
	}

	// Drop all existing indexes except _id_
	if _, err := collection.Indexes().DropAll(indexCtx); err != nil {
		return fmt.Errorf("failed to drop existing indexes: %w", err)
	}

	// Create new indexes
	if _, err := collection.Indexes().CreateMany(indexCtx, indexes); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	// Verify indexes were created correctly
	cursor, err := collection.Indexes().List(indexCtx)
	if err != nil {
		return fmt.Errorf("failed to list indexes: %w", err)
	}
	defer cursor.Close(indexCtx)

	var createdIndexes []bson.M
	if err = cursor.All(indexCtx, &createdIndexes); err != nil {
		return fmt.Errorf("failed to read created indexes: %w", err)
	}

	// Verify all required indexes exist
	requiredIndexes := map[string]bool{
		"_id_":             false,
		"swiftCode_unique": false,
		"countryISO2":      false,
	}

	for _, idx := range createdIndexes {
		name := idx["name"].(string)
		if _, ok := requiredIndexes[name]; ok {
			requiredIndexes[name] = true
		}
	}

	for name, found := range requiredIndexes {
		if !found {
			return fmt.Errorf("required index %s was not created", name)
		}
	}

	return nil
}

func createCollection(db *mongo.Database, name string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	if err := db.CreateCollection(ctx, name); err != nil {
		return fmt.Errorf("error while creating a collection: %w", err)
	}

	return nil
}

// GetInstance returns the current database configuration instance
func GetInstance() *Config {
	return instance
}

// Disconnect closes the database connection
func (c *Config) Disconnect(ctx context.Context) error {
	if c.Client != nil {
		return c.Client.Disconnect(ctx)
	}
	return nil
}
