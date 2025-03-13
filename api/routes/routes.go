package routes

import (
	"github.com/MarcinZ20/bankAPI/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func BankRoutes(app *fiber.App) {
	app.Get("/v1/swift-codes/:swiftCode", handlers.GetSwiftCodesBySwiftCode)
	app.Get("/v1/swift-codes/country/:countryISO2", handlers.GetSwiftCodesByCountryCode)
	app.Post("/v1/swift-codes", handlers.AddNewSwiftCode)
	app.Delete("/v1/swift-codes/:swiftCode", handlers.DeleteSwiftCode)
}
