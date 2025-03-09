package handlers

import (
	"fmt"
	"strings"

	"github.com/MarcinZ20/bankAPI/api/middleware"
	"github.com/MarcinZ20/bankAPI/api/responses"
	"github.com/MarcinZ20/bankAPI/config"
	"github.com/MarcinZ20/bankAPI/internal/utils"
	"github.com/MarcinZ20/bankAPI/pkg/handlers/db"
	"github.com/MarcinZ20/bankAPI/pkg/models"
	"github.com/MarcinZ20/bankAPI/pkg/services"
	"github.com/MarcinZ20/bankAPI/pkg/transform"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetSwiftCodesBySwiftCode(c *fiber.Ctx) error {
	ctx, ok := middleware.GetRequestContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get request context")
	}

	collection, err := db.GetCollection(config.MongoClient)
	if err != nil {
		return responses.DatabaseError(err)
	}

	bankService := services.NewBankService(collection)
	swiftCode := c.Params("swiftCode")

	if !utils.IsValidSwiftCodeFormat(swiftCode) {
		return responses.ValidationError(fmt.Sprintf("Invalid SWIFT code format: %v", swiftCode))
	}

	if strings.HasSuffix(swiftCode, "XXX") {
		hq, err := bankService.GetHeadquarter(ctx, swiftCode)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return responses.NotFoundError("headquarter", swiftCode)
			}
			return responses.DatabaseError(err)
		}

		response := new(responses.HeadquarterResponse)
		if err := response.FromModel(hq); err != nil {
			return responses.FormattingResponseError("Error while formatting response")
		}

		return responses.NewSuccessResponse(c, response)
	}

	branch, err := bankService.GetBranch(ctx, swiftCode)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return responses.NotFoundError("branch", swiftCode)
		}
		return responses.DatabaseError(err)
	}

	response := new(responses.LongBankResponse)
	if err := response.FromModel(branch); err != nil {
		return responses.FormattingResponseError("Error while formatting response")
	}

	return responses.NewSuccessResponse(c, response)
}

func GetSwiftCodesByCountryCode(c *fiber.Ctx) error {
	ctx, ok := middleware.GetRequestContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get request context")
	}

	collection, err := db.GetCollection(config.MongoClient)
	if err != nil {
		return responses.DatabaseError(err)
	}

	bankService := services.NewBankService(collection)
	countryCode := c.Params("countryISO2")

	if !utils.IsValidCountryCode(countryCode) {
		return responses.ValidationError(fmt.Sprintf("Invalid country code format: %v", countryCode))
	}

	foundData, err := bankService.GetBanksByCountryCode(ctx, countryCode)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return responses.NotFoundError("records", countryCode)
		}
		return responses.DatabaseError(err)
	}

	response := responses.GetSwiftCodesByCountryCodeResponse{
		CountryISO2: countryCode,
		CountryName: foundData[0].CountryName,
		SwiftCodes:  make([]responses.ShortBankResponse, 0),
	}

	for _, hq := range foundData {
		shortResponse := responses.ShortBankResponse{
			Address:       hq.Address,
			BankName:      hq.BankName,
			CountryISO2:   hq.CountryISO2,
			IsHeadquarter: true,
			SwiftCode:     hq.SwiftCode,
		}
		response.SwiftCodes = append(response.SwiftCodes, shortResponse)

		for _, branch := range hq.Branches {
			branchResponse := responses.ShortBankResponse{
				Address:       branch.Address,
				BankName:      branch.BankName,
				CountryISO2:   branch.CountryISO2,
				IsHeadquarter: false,
				SwiftCode:     branch.SwiftCode,
			}
			response.SwiftCodes = append(response.SwiftCodes, branchResponse)
		}
	}

	return responses.NewSuccessResponse(c, response)
}

func AddNewSwiftCode(c *fiber.Ctx) error {
	ctx, ok := middleware.GetRequestContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get request context")
	}

	collection, err := db.GetCollection(config.MongoClient)
	if err != nil {
		return responses.DatabaseError(err)
	}

	record := new(models.Branch)
	if err := c.BodyParser(record); err != nil {
		return responses.ValidationError(fmt.Sprintf("Invalid request body: %v", err))
	}

	if !utils.IsValidSwiftCodeFormat(record.SwiftCode) {
		return responses.ValidationError(fmt.Sprintf("Invalid SWIFT code format: %v", record))
	}

	if !utils.IsValidCountryCode(record.CountryISO2) {
		return responses.ValidationError(fmt.Sprintf("Invalid country code format: %v", record.CountryISO2))
	}

	transform.TransformRequestModel(record)

	bankService := services.NewBankService(collection)

	if strings.HasSuffix(record.SwiftCode, "XXX") {
		if !record.IsHeadquarter {
			return responses.ValidationError("SWIFT code ends with XXX, but <isHeadquarter> is false")
		}

		hq := models.Headquarter{
			SwiftCode:     record.SwiftCode,
			BankName:      record.BankName,
			Address:       record.Address,
			CountryName:   record.CountryName,
			CountryISO2:   record.CountryISO2,
			IsHeadquarter: true,
			Branches:      []models.Branch{},
		}

		if err := bankService.AddHeadquarter(ctx, &hq); err != nil {
			if err.Error() == "headquarter already exists" {
				return responses.AlreadyExistsError(fmt.Sprintf("Headquarter with SWIFT code %s already exists", record.SwiftCode))
			}
			return responses.DatabaseError(err)
		}

		return responses.NewSuccessResponse(c, fiber.Map{
			"message": "Headquarter created successfully",
		})
	}

	parentHqSwiftCode := record.SwiftCode[0:8] + "XXX"
	if err := bankService.AddBranch(ctx, parentHqSwiftCode, record); err != nil {
		if err == mongo.ErrNoDocuments {
			return responses.NotFoundError("parent headquarter", parentHqSwiftCode)
		}
		if err.Error() == "branch already exists" {
			return responses.AlreadyExistsError(fmt.Sprintf("Branch with SWIFT code %s already exists", record.SwiftCode))
		}
		return responses.DatabaseError(err)
	}

	return responses.NewSuccessResponse(c, fiber.Map{
		"message": "Branch added successfully",
	})
}

func DeleteSwiftCode(c *fiber.Ctx) error {
	ctx, ok := middleware.GetRequestContext(c)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get request context")
	}

	collection, err := db.GetCollection(config.MongoClient)
	if err != nil {
		return responses.DatabaseError(err)
	}

	swiftCode := c.Params("swiftCode")
	if !utils.IsValidSwiftCodeFormat(swiftCode) {
		return responses.ValidationError(fmt.Sprintf("Invalid SWIFT code format: %v", swiftCode))
	}

	bankService := services.NewBankService(collection)

	if strings.HasSuffix(swiftCode, "XXX") {
		if err := bankService.DeleteHeadquarter(ctx, swiftCode); err != nil {
			if err == mongo.ErrNoDocuments {
				return responses.NotFoundError("headquarter", swiftCode)
			}
			return responses.DatabaseError(err)
		}

		return responses.NewSuccessResponse(c, fiber.Map{
			"message": "Headquarter was deleted successfully",
		})
	}

	parentHqSwiftCode := swiftCode[0:8] + "XXX"
	if err := bankService.DeleteBranch(ctx, swiftCode, parentHqSwiftCode); err != nil {
		if err == mongo.ErrNoDocuments {
			return responses.NotFoundError("branch", swiftCode)
		}
		return responses.DatabaseError(err)
	}

	return responses.NewSuccessResponse(c, fiber.Map{
		"message": "Branch was deleted successfully",
	})
}
