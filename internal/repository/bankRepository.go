package repository

import (
	"context"
	"fmt"

	"github.com/MarcinZ20/bankAPI/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BankRepository handles all database operations
type BankRepository struct {
	collection *mongo.Collection
}

// NewBankRepository creates a new bank repository
func NewBankRepository(collection *mongo.Collection) *BankRepository {
	return &BankRepository{
		collection: collection,
	}
}

// FindHeadquarter finds a headquarter by SWIFT code
func (r *BankRepository) FindHeadquarter(ctx context.Context, swiftCode string) (*models.Headquarter, error) {
	filter := bson.D{
		{Key: "swiftCode", Value: swiftCode},
		{Key: "isHeadquarter", Value: true},
	}

	var hq models.Headquarter
	err := r.collection.FindOne(ctx, filter).Decode(&hq)
	if err != nil {
		return nil, fmt.Errorf("failed to find headquarter: %w", err)
	}

	return &hq, nil
}

// FindBranch finds a branch by SWIFT code
func (r *BankRepository) FindBranch(ctx context.Context, swiftCode, parentSwiftCode string) (*models.Branch, error) {
	filter := bson.D{
		{Key: "swiftCode", Value: parentSwiftCode},
		{Key: "isHeadquarter", Value: true},
		{Key: "branches.swiftCode", Value: swiftCode},
	}

	opts := options.FindOne().SetProjection(bson.D{
		{Key: "branches.$", Value: 1},
	})

	var hq models.Headquarter
	err := r.collection.FindOne(ctx, filter, opts).Decode(&hq)
	if err != nil {
		return nil, fmt.Errorf("failed to find branch: %w", err)
	}

	if len(hq.Branches) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &hq.Branches[0], nil
}

// FindBanksByCountry finds all banks in a given country
func (r *BankRepository) FindBanksByCountry(ctx context.Context, countryCode string) ([]models.Headquarter, error) {
	filter := bson.D{
		{Key: "countryISO2", Value: countryCode},
		{Key: "isHeadquarter", Value: true},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find banks: %w", err)
	}
	defer cursor.Close(ctx)

	var foundData []models.Headquarter
	if err := cursor.All(ctx, &foundData); err != nil {
		return nil, fmt.Errorf("failed to decode banks: %w", err)
	}

	if len(foundData) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return foundData, nil
}

// CreateHeadquarter creates a new headquarter
func (r *BankRepository) CreateHeadquarter(ctx context.Context, hq *models.Headquarter) error {
	exists := bson.D{
		{Key: "swiftCode", Value: hq.SwiftCode},
		{Key: "isHeadquarter", Value: true},
	}

	var existing models.Headquarter
	if err := r.collection.FindOne(ctx, exists).Decode(&existing); err == nil {
		return fmt.Errorf("headquarter already exists")
	}

	_, err := r.collection.InsertOne(ctx, hq)
	if err != nil {
		return fmt.Errorf("failed to create headquarter: %w", err)
	}

	return nil
}

// AddBranch adds a new branch to a headquarter
func (r *BankRepository) AddBranch(ctx context.Context, parentSwiftCode string, branch *models.Branch) error {
	filter := bson.D{
		{Key: "swiftCode", Value: parentSwiftCode},
		{Key: "isHeadquarter", Value: true},
	}

	var hq models.Headquarter
	if err := r.collection.FindOne(ctx, filter).Decode(&hq); err != nil {
		return fmt.Errorf("failed to find parent headquarter: %w", err)
	}

	for _, b := range hq.Branches {
		if b.SwiftCode == branch.SwiftCode {
			return fmt.Errorf("branch already exists")
		}
	}

	update := bson.D{{
		Key: "$push",
		Value: bson.D{{
			Key:   "branches",
			Value: branch,
		}},
	}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add branch: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("failed to modify headquarter")
	}

	return nil
}

// DeleteHeadquarter deletes a headquarter and all its branches
func (r *BankRepository) DeleteHeadquarter(ctx context.Context, swiftCode string) error {
	filter := bson.D{
		{Key: "swiftCode", Value: swiftCode},
		{Key: "isHeadquarter", Value: true},
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete headquarter: %w", err)
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// DeleteBranch removes a branch from its headquarter
func (r *BankRepository) DeleteBranch(ctx context.Context, swiftCode, parentSwiftCode string) error {
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

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
