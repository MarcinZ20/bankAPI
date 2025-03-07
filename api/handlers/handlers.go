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

func GetSwiftCodesByCountryCode(ctx *fiber.Ctx) error {
	return nil
}

func AddNewSwiftCode(ctx *fiber.Ctx) error {
	return nil
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
