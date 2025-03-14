package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MarcinZ20/bankAPI/internal/database"
	"github.com/MarcinZ20/bankAPI/internal/services"
	"github.com/MarcinZ20/bankAPI/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Collection, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	err = client.Ping(ctx, nil)
	require.NoError(t, err, "Failed to connect to MongoDB. Make sure the MongoDB container is running")

	db := client.Database("test_db")
	collection := db.Collection("test_banks")

	cleanup := func() {
		collection.Drop(ctx)
		client.Disconnect(ctx)
	}

	return collection, cleanup
}

func setupTestApp(collection *mongo.Collection) *fiber.App {
	app := fiber.New()

	dbConfig := &database.Config{
		Collection: collection,
	}

	sm := services.NewServiceManager(dbConfig)

	app.Get("/api/v1/swift-codes/:swiftCode", GetSwiftCodesBySwiftCode)
	app.Post("/api/v1/swift-codes", AddNewSwiftCode)
	app.Delete("/api/v1/swift-codes/:swiftCode", DeleteSwiftCode)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("serviceManager", sm)
		return c.Next()
	})

	return app
}

func TestGetSwiftCodesBySwiftCode(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	app := setupTestApp(collection)

	hq := &models.Headquarter{
		SwiftCode:     "DEUTDEFFXXX",
		BankName:      "Deutsche Bank",
		CountryISO2:   "DE",
		CountryName:   "Germany",
		Address:       "Taunusanlage 12",
		IsHeadquarter: true,
		Branches:      []models.Branch{},
	}

	ctx := context.Background()
	_, err := collection.InsertOne(ctx, hq)
	require.NoError(t, err)

	tests := []struct {
		name           string
		swiftCode      string
		expectedStatus int
		validateResp   func(*testing.T, []byte)
	}{
		{
			name:           "Get existing headquarter",
			swiftCode:      "DEUTDEFFXXX",
			expectedStatus: fiber.StatusOK,
			validateResp: func(t *testing.T, body []byte) {
				var resp map[string]any
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)

				data, ok := resp["data"].(map[string]any)
				require.True(t, ok, "Response data should be an object")

				assert.Equal(t, "DEUTDEFFXXX", data["swiftCode"])
				assert.Equal(t, "Deutsche Bank", data["bankName"])
				assert.Equal(t, "DE", data["countryISO2"])
			},
		},
		{
			name:           "Get non-existing headquarter",
			swiftCode:      "NONEXISTXXX",
			expectedStatus: fiber.StatusNotFound,
		},
		{
			name:           "Invalid SWIFT code format",
			swiftCode:      "INVALID",
			expectedStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/swift-codes/"+tt.swiftCode, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.validateResp != nil {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				tt.validateResp(t, body)
			}
		})
	}
}
