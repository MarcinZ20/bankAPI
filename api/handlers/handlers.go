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
	"go.mongodb.org/mongo-driver/mongo/options"
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
					"error": "No headquarter found with this SWIFT code",
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
	var hqResult bson.M

	parentHqSwiftCode := swiftCode[0:8] + "XXX"

	filter := bson.D{
		{Key: "swiftCode", Value: parentHqSwiftCode},
		{Key: "isHeadquarter", Value: true},
	}

	err = collection.FindOne(ctx, filter).Decode(&hqResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "No headquarter found with this SWIFT code",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	branchFilter := bson.M{
		"branches": bson.M{
			"$elemMatch": bson.M{"swift_code": swiftCode},
		},
	}

	projection := bson.M{
		"branches.$": 1,
		"_id":        0,
	}

	options := options.FindOne().SetProjection(projection)

	err = collection.FindOne(ctx, branchFilter, options).Decode(&branch)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "No headquarter found with this SWIFT code",
			})
		} else {
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}
	}

	return c.Status(http.StatusOK).JSON(branch)
}

func GetSwiftCodesByCountryCode(ctx *fiber.Ctx) error {
	return nil
}

func AddNewSwiftCode(ctx *fiber.Ctx) error {
	return nil
}

func DeleteSwiftCode(ctx *fiber.Ctx) error {
	return nil
}
