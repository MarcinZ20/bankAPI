package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MarcinZ20/bankAPI/config"
	"github.com/MarcinZ20/bankAPI/internal/utils"
	"github.com/MarcinZ20/bankAPI/pkg/handlers/db"
	"github.com/MarcinZ20/bankAPI/pkg/models"
	"github.com/MarcinZ20/bankAPI/pkg/transform"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetSwiftCodesBySwiftCode(c *fiber.Ctx) error {
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()

	collection, err := db.GetCollection(config.MongoClient)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	swiftCode := c.Params("swiftCode")

	// 1. Case swift code is not valid
	if !utils.IsValidSwiftCodeFormat(swiftCode) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid SWIFT code format: %v", swiftCode),
		})
	}

	// 2. Case headquarter is queried
	if strings.HasSuffix(swiftCode, "XXX") {
		var hq models.Headquarter

		filter := bson.D{
			{Key: "swiftCode", Value: swiftCode},
			{Key: "isHeadquarter", Value: true},
		}

		err := collection.FindOne(ctx, filter).Decode(&hq)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{
					"error": fmt.Sprintf("No headquarter found with this SWIFT code: %v", swiftCode),
				})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}
		fmt.Printf("Headquarter: %v\n", hq)
		return c.Status(http.StatusOK).JSON(hq)
	}

	// 3. Case branch is queried
	var branch models.Branch
	parentHqSwiftCode := swiftCode[0:8] + "XXX"

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "swiftCode", Value: parentHqSwiftCode},
			{Key: "isHeadquarter", Value: true},
		}}},
		{{Key: "$unwind", Value: "$branches"}},
		{{Key: "$match", Value: bson.D{
			{Key: "branches.swiftCode", Value: swiftCode},
		}}},
		{{Key: "$replaceRoot", Value: bson.D{
			{Key: "newRoot", Value: "$branches"},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Database error: %v", err),
		})
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("No branch found with SWIFT code: %v", swiftCode),
		})
	}

	if err := cursor.Decode(&branch); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Error decoding branch: %v", err),
		})
	}

	return c.Status(http.StatusOK).JSON(branch)
}

func GetSwiftCodesByCountryCode(c *fiber.Ctx) error {
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()

	collection, err := db.GetCollection(config.MongoClient)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	countryCode := c.Params("countryISO2")

	if !utils.IsValidCountryCode(countryCode) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid country code format: %v", countryCode),
		})
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "countryISO2", Value: countryCode},
			{Key: "isHeadquarter", Value: true},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Database error: %v", err),
		})
	}
	defer cursor.Close(ctx)

	var results []models.Headquarter
	if err := cursor.All(ctx, &results); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Error decoding results: %v", err),
		})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("No banks found for country code: %v", countryCode),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": len(results),
		"banks": results,
	})
}

func AddNewSwiftCode(c *fiber.Ctx) error {
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()

	collection, err := db.GetCollection(config.MongoClient)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	var record models.Branch
	if err := c.BodyParser(&record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid request body: %v", err),
		})
	}

	if !utils.IsValidSwiftCodeFormat(record.SwiftCode) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid SWIFT code format: %v", record.SwiftCode),
		})
	}

	if !utils.IsValidCountryCode(record.CountryISO2Code) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid country code format: %v", record.CountryISO2Code),
		})
	}

	transform.TransformRequestModel(&record)

	// Case 1: Adding a headquarter
	if strings.HasSuffix(record.SwiftCode, "XXX") {
		if !record.IsHeadquarter {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "SWIFT code ends with XXX but isHeadquarter is false",
			})
		}

		hq := models.Headquarter{
			SwiftCode:       record.SwiftCode,
			BankName:        record.Name,
			Address:         record.Address,
			CountryName:     record.CountryName,
			CountryISO2Code: record.CountryISO2Code,
			IsHeadquarter:   true,
			Branches:        []models.Branch{},
		}

		// Check if headquarter already exists
		exists := bson.D{
			{Key: "swiftCode", Value: record.SwiftCode},
			{Key: "isHeadquarter", Value: true},
		}

		var existing models.Headquarter
		if err := collection.FindOne(ctx, exists).Decode(&existing); err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": fmt.Sprintf("Headquarter with SWIFT code %s already exists", record.SwiftCode),
			})
		}

		res, err := collection.InsertOne(ctx, hq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Database error: %v", err),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Headquarter created successfully",
			"id":      res.InsertedID,
		})
	}

	// Case 2: Adding a branch
	parentHqSwiftCode := record.SwiftCode[0:8] + "XXX"

	// Check if parent headquarter exists
	filter := bson.D{
		{Key: "swiftCode", Value: parentHqSwiftCode},
		{Key: "isHeadquarter", Value: true},
	}

	var hq models.Headquarter
	if err := collection.FindOne(ctx, filter).Decode(&hq); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("Parent headquarter with SWIFT code %s not found", parentHqSwiftCode),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Database error: %v", err),
		})
	}

	// Check if branch already exists
	for _, branch := range hq.Branches {
		if branch.SwiftCode == record.SwiftCode {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": fmt.Sprintf("Branch with SWIFT code %s already exists", record.SwiftCode),
			})
		}
	}

	// Add branch to headquarter
	update := bson.D{{
		Key: "$push",
		Value: bson.D{{
			Key:   "branches",
			Value: record,
		}},
	}}

	x, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Database error: %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": fmt.Sprintf("Branch added successfully: %v", x),
	})
}

func DeleteSwiftCode(c *fiber.Ctx) error {
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()

	collection, err := db.GetCollection(config.MongoClient)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}

	swiftCode := c.Params("swiftCode")

	// 1. Case swift code is not valid
	if !utils.IsValidSwiftCodeFormat(swiftCode) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid SWIFT code format: %v", swiftCode),
		})
	}

	// 2. Case headquarter
	if strings.HasSuffix(swiftCode, "XXX") {
		filter := bson.D{
			{Key: "swiftCode", Value: swiftCode},
			{Key: "isHeadquarter", Value: true},
		}

		deleted, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		if deleted.DeletedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("No headquarter found with SWIFT code: %v", swiftCode),
			})
		}

		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"message": "Headquarter was deleted successfully",
		})
	}

	// 3. Case branch
	parentHqSwiftCode := swiftCode[0:8] + "XXX"

	filter := bson.D{
		{Key: "swiftCode", Value: parentHqSwiftCode},
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

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Database error: %v", err),
		})
	}

	if result.ModifiedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("No branch found with SWIFT code: %v", swiftCode),
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Branch was deleted successfully",
	})
}
