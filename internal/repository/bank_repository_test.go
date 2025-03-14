package repository

import (
	"context"
	"testing"
	"time"

	"github.com/MarcinZ20/bankAPI/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Sets up a test db
func setupTestDB(t *testing.T) (*mongo.Collection, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	db := client.Database("test_db")
	collection := db.Collection("test_banks")

	cleanup := func() {
		collection.Drop(ctx)
		client.Disconnect(ctx)
	}

	return collection, cleanup
}

func TestBankRepository_FindHeadquarter(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBankRepository(collection)
	ctx := context.Background()

	// random test data
	hq := &models.Headquarter{
		SwiftCode:     "DEUTDEFF",
		BankName:      "Deutsche Bank",
		CountryISO2:   "DE",
		IsHeadquarter: true,
		Branches:      []models.Branch{},
	}

	_, err := collection.InsertOne(ctx, hq)
	require.NoError(t, err)

	tests := []struct {
		name      string
		swiftCode string
		wantErr   bool
	}{
		{
			name:      "Existing headquarter",
			swiftCode: "DEUTDEFF",
			wantErr:   false,
		},
		{
			name:      "Non-existing headquarter",
			swiftCode: "NONEXIST",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindHeadquarter(ctx, tt.swiftCode)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.swiftCode, result.SwiftCode)
		})
	}
}

func TestBankRepository_FindBranch(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBankRepository(collection)
	ctx := context.Background()

	// random test data
	branch := models.Branch{
		SwiftCode:   "DEUTDEFF100",
		BankName:    "Deutsche Bank Berlin",
		CountryISO2: "DE",
	}

	hq := &models.Headquarter{
		SwiftCode:     "DEUTDEFF",
		BankName:      "Deutsche Bank",
		CountryISO2:   "DE",
		IsHeadquarter: true,
		Branches:      []models.Branch{branch},
	}

	_, err := collection.InsertOne(ctx, hq)
	require.NoError(t, err)

	tests := []struct {
		name        string
		branchSwift string
		parentSwift string
		wantErr     bool
	}{
		{
			name:        "Existing branch",
			branchSwift: "DEUTDEFF100",
			parentSwift: "DEUTDEFF",
			wantErr:     false,
		},
		{
			name:        "Non-existing branch",
			branchSwift: "NONEXIST",
			parentSwift: "DEUTDEFF",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindBranch(ctx, tt.branchSwift, tt.parentSwift)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.branchSwift, result.SwiftCode)
		})
	}
}

func TestBankRepository_FindBanksByCountry(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBankRepository(collection)
	ctx := context.Background()

	// just some random test data
	hqs := []any{
		&models.Headquarter{
			SwiftCode:     "DEUTDEFF",
			BankName:      "Deutsche Bank",
			CountryISO2:   "DE",
			IsHeadquarter: true,
		},
		&models.Headquarter{
			SwiftCode:     "COMMDEFF",
			BankName:      "Commerzbank",
			CountryISO2:   "DE",
			IsHeadquarter: true,
		},
	}

	_, err := collection.InsertMany(ctx, hqs)
	require.NoError(t, err)

	tests := []struct {
		name        string
		countryCode string
		wantCount   int
		wantErr     bool
	}{
		{
			name:        "Existing country",
			countryCode: "DE",
			wantCount:   2,
			wantErr:     false,
		},
		{
			name:        "Non-existing country",
			countryCode: "XX",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := repo.FindBanksByCountry(ctx, tt.countryCode)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, results, tt.wantCount)
		})
	}
}

func TestBankRepository_CreateHeadquarter(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBankRepository(collection)
	ctx := context.Background()

	hq := &models.Headquarter{
		SwiftCode:     "DEUTDEFF",
		BankName:      "Deutsche Bank",
		CountryISO2:   "DE",
		IsHeadquarter: true,
	}

	tests := []struct {
		name    string
		hq      *models.Headquarter
		wantErr bool
		setup   func()
	}{
		{
			name:    "New headquarter",
			hq:      hq,
			wantErr: false,
			setup:   func() {},
		},
		{
			name:    "Duplicate headquarter",
			hq:      hq,
			wantErr: true,
			setup: func() {
				_, err := collection.InsertOne(ctx, hq)
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection.Drop(ctx)
			tt.setup()

			err := repo.CreateHeadquarter(ctx, tt.hq)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			var result models.Headquarter
			err = collection.FindOne(ctx, bson.D{{Key: "swiftCode", Value: tt.hq.SwiftCode}}).Decode(&result)
			assert.NoError(t, err)
			assert.Equal(t, tt.hq.SwiftCode, result.SwiftCode)
		})
	}
}

func TestBankRepository_DeleteHeadquarter(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBankRepository(collection)
	ctx := context.Background()

	hq := &models.Headquarter{
		SwiftCode:     "DEUTDEFF",
		BankName:      "Deutsche Bank",
		CountryISO2:   "DE",
		IsHeadquarter: true,
	}

	tests := []struct {
		name      string
		swiftCode string
		wantErr   bool
		setup     func()
	}{
		{
			name:      "Existing headquarter",
			swiftCode: "DEUTDEFF",
			wantErr:   false,
			setup: func() {
				_, err := collection.InsertOne(ctx, hq)
				require.NoError(t, err)
			},
		},
		{
			name:      "Non-existing headquarter",
			swiftCode: "NONEXIST",
			wantErr:   true,
			setup:     func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection.Drop(ctx)
			tt.setup()

			err := repo.DeleteHeadquarter(ctx, tt.swiftCode)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			count, err := collection.CountDocuments(ctx, bson.D{{Key: "swiftCode", Value: tt.swiftCode}})
			assert.NoError(t, err)
			assert.Equal(t, int64(0), count)
		})
	}
}

func TestBankRepository_AddBranch(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBankRepository(collection)
	ctx := context.Background()

	hq := &models.Headquarter{
		SwiftCode:     "DEUTDEFF",
		BankName:      "Deutsche Bank",
		CountryISO2:   "DE",
		IsHeadquarter: true,
		Branches:      []models.Branch{},
	}

	branch := &models.Branch{
		SwiftCode:   "DEUTDEFF100",
		BankName:    "Deutsche Bank Berlin",
		CountryISO2: "DE",
	}

	tests := []struct {
		name            string
		parentSwiftCode string
		branch          *models.Branch
		wantErr         bool
		setup           func()
	}{
		{
			name:            "Add new branch",
			parentSwiftCode: "DEUTDEFF",
			branch:          branch,
			wantErr:         false,
			setup: func() {
				_, err := collection.InsertOne(ctx, hq)
				require.NoError(t, err)
			},
		},
		{
			name:            "Non-existing headquarter",
			parentSwiftCode: "NONEXIST",
			branch:          branch,
			wantErr:         true,
			setup:           func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection.Drop(ctx)
			tt.setup()

			err := repo.AddBranch(ctx, tt.parentSwiftCode, tt.branch)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			var result models.Headquarter
			err = collection.FindOne(ctx, bson.D{{Key: "swiftCode", Value: tt.parentSwiftCode}}).Decode(&result)
			assert.NoError(t, err)
			assert.Contains(t, result.Branches, *tt.branch)
		})
	}
}

func TestBankRepository_DeleteBranch(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBankRepository(collection)
	ctx := context.Background()

	branch := models.Branch{
		SwiftCode:   "DEUTDEFF100",
		BankName:    "Deutsche Bank Berlin",
		CountryISO2: "DE",
	}

	hq := &models.Headquarter{
		SwiftCode:     "DEUTDEFF",
		BankName:      "Deutsche Bank",
		CountryISO2:   "DE",
		IsHeadquarter: true,
		Branches:      []models.Branch{branch},
	}

	tests := []struct {
		name        string
		branchSwift string
		parentSwift string
		wantErr     bool
		setup       func()
	}{
		{
			name:        "Delete existing branch",
			branchSwift: "DEUTDEFF100",
			parentSwift: "DEUTDEFF",
			wantErr:     false,
			setup: func() {
				_, err := collection.InsertOne(ctx, hq)
				require.NoError(t, err)
			},
		},
		{
			name:        "Delete from non-existing headquarter",
			branchSwift: "DEUTDEFF100",
			parentSwift: "NONEXIST",
			wantErr:     true,
			setup:       func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection.Drop(ctx)
			tt.setup()

			err := repo.DeleteBranch(ctx, tt.branchSwift, tt.parentSwift)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			var result models.Headquarter
			err = collection.FindOne(ctx, bson.D{{Key: "swiftCode", Value: tt.parentSwift}}).Decode(&result)
			assert.NoError(t, err)
			for _, b := range result.Branches {
				assert.NotEqual(t, tt.branchSwift, b.SwiftCode)
			}
		})
	}
}
