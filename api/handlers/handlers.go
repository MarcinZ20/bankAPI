package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MarcinZ20/bankAPI/config"
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

	filter := bson.D{
		{Key: "swiftCode", Value: swiftCode},
		{Key: "isHeadquarter", Value: true},
	}

	if strings.HasSuffix(swiftCode, "XXX") {
		var hq models.Headquarter

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

	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
		"error": "Invalid SWIFT code format",
	})
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
