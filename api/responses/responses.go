package responses

import (
	"fmt"

	"github.com/MarcinZ20/bankAPI/pkg/models"
	"github.com/gofiber/fiber/v2"
)

// Defines how to convert domain models to response DTOs
type Response interface {
	FromModel(models.BankEntity) error
}

// Creates a consistent success response with status 200 OK
func NewSuccessResponse(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(data)
}

type HeadquarterResponse struct {
	Address       string              `json:"address"`
	BankName      string              `json:"bankName"`
	CountryISO2   string              `json:"countryISO2"`
	CountryName   string              `json:"countryName"`
	IsHeadquarter bool                `json:"isHeadquarter"`
	SwiftCode     string              `json:"swiftCode"`
	Branches      []ShortBankResponse `json:"branches,omitempty"`
}

type ShortBankResponse struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type LongBankResponse struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type GetSwiftCodesByCountryCodeResponse struct {
	CountryISO2 string              `json:"countryISO2"`
	CountryName string              `json:"countryName"`
	SwiftCodes  []ShortBankResponse `json:"swiftCodes"`
}

func (r *HeadquarterResponse) FromModel(model models.BankEntity) error {
	hq, ok := model.(*models.Headquarter)
	if !ok {
		return fmt.Errorf("model is not a Headquarter")
	}

	r.Address = hq.Address
	r.BankName = hq.BankName
	r.CountryISO2 = hq.CountryISO2
	r.CountryName = hq.CountryName
	r.IsHeadquarter = hq.IsHeadquarter
	r.SwiftCode = hq.SwiftCode

	if len(hq.Branches) > 0 {
		r.Branches = make([]ShortBankResponse, len(hq.Branches))
		for i, branch := range hq.Branches {
			r.Branches[i] = ShortBankResponse{
				Address:       branch.Address,
				BankName:      branch.BankName,
				CountryISO2:   branch.CountryISO2,
				IsHeadquarter: branch.IsHeadquarter,
				SwiftCode:     branch.SwiftCode,
			}
		}
	}
	return nil
}

func (r *LongBankResponse) FromModel(model models.BankEntity) error {
	r.Address = model.GetAddress()
	r.BankName = model.GetBankName()
	r.CountryISO2 = model.GetCountryISO2()
	r.CountryName = model.GetCountryName()
	r.IsHeadquarter = model.IsHq()
	r.SwiftCode = model.GetSwiftCode()
	return nil
}

func (r *ShortBankResponse) FromModel(model models.BankEntity) error {
	r.Address = model.GetAddress()
	r.BankName = model.GetBankName()
	r.CountryISO2 = model.GetCountryISO2()
	r.IsHeadquarter = model.IsHq()
	r.SwiftCode = model.GetSwiftCode()
	return nil
}
