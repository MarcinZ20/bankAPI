package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MarcinZ20/bankAPI/pkg/models"
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

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	mongoUri := os.Getenv("MONGO_URI")
	clientOptions := options.Client().SetTimeout(3 * time.Second).ApplyURI(mongoUri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to mongoDB: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("error while pinging the database: %w", err)
	}

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION"))

	// Create indexes
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "swiftCode", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "countryISO2", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
	}

	if _, err := collection.Indexes().CreateMany(ctx, indexes); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	instance = &Config{
		Client:     client,
		Collection: collection,
	}

	return instance, nil
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

// GetHeadquarter retrieves a headquarter by SWIFT code
func (c *Config) GetHeadquarter(ctx context.Context, swiftCode string) (*models.Headquarter, error) {
	filter := bson.D{
		{Key: "swiftCode", Value: swiftCode},
		{Key: "isHeadquarter", Value: true},
	}

	var hq models.Headquarter
	err := c.Collection.FindOne(ctx, filter).Decode(&hq)
	if err != nil {
		return nil, err
	}

	return &hq, nil
}

// GetBranch retrieves a branch by SWIFT code
func (c *Config) GetBranch(ctx context.Context, swiftCode string) (*models.Branch, error) {
	parentHqSwiftCode := swiftCode[0:8] + "XXX"

	filter := bson.D{
		{Key: "swiftCode", Value: parentHqSwiftCode},
		{Key: "isHeadquarter", Value: true},
		{Key: "branches.swiftCode", Value: swiftCode},
	}

	opts := options.FindOne().SetProjection(bson.D{
		{Key: "branches.$", Value: 1},
	})

	var hq models.Headquarter
	err := c.Collection.FindOne(ctx, filter, opts).Decode(&hq)
	if err != nil {
		return nil, err
	}

	if len(hq.Branches) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &hq.Branches[0], nil
}

// GetBanksByCountry retrieves all banks in a given country
func (c *Config) GetBanksByCountry(ctx context.Context, countryCode string) ([]models.Headquarter, error) {
	filter := bson.D{
		{Key: "countryISO2", Value: countryCode},
		{Key: "isHeadquarter", Value: true},
	}

	cursor, err := c.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var banks []models.Headquarter
	if err := cursor.All(ctx, &banks); err != nil {
		return nil, err
	}

	if len(banks) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return banks, nil
}

// AddHeadquarter creates a new headquarter
func (c *Config) AddHeadquarter(ctx context.Context, hq *models.Headquarter) error {
	_, err := c.Collection.InsertOne(ctx, hq)
	return err
}

// AddBranch adds a branch to an existing headquarter
func (c *Config) AddBranch(ctx context.Context, parentSwiftCode string, branch *models.Branch) error {
	filter := bson.D{
		{Key: "swiftCode", Value: parentSwiftCode},
		{Key: "isHeadquarter", Value: true},
	}

	update := bson.D{{
		Key: "$push",
		Value: bson.D{{
			Key:   "branches",
			Value: branch,
		}},
	}}

	result, err := c.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// DeleteHeadquarter deletes a headquarter and all its branches
func (c *Config) DeleteHeadquarter(ctx context.Context, swiftCode string) error {
	filter := bson.D{
		{Key: "swiftCode", Value: swiftCode},
		{Key: "isHeadquarter", Value: true},
	}

	result, err := c.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// DeleteBranch removes a branch from its headquarter
func (c *Config) DeleteBranch(ctx context.Context, swiftCode, parentSwiftCode string) error {
	filter := bson.D{
		{Key: "swiftCode", Value: parentSwiftCode},
		{Key: "isHeadquarter", Value: true},
	}

	update := bson.D{{
		Key: "$pull",
		Value: bson.D{{
			Key: "branches",
			Value: bson.D{{
				Key:   "swiftCode",
				Value: swiftCode,
			}},
		}},
	}}

	result, err := c.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
