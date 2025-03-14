package app

import (
	"time"

	"github.com/MarcinZ20/bankAPI/api/middleware"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

// Holds application configuration
type Config struct {
	Server *fiber.App
}

var instance *Config

// Sets up the application configuration
func Initialize() *Config {
	if instance != nil {
		return instance
	}

	fiberConfig := fiber.Config{
		AppName:       "bankAPI v1.0",
		CaseSensitive: true,
		StrictRouting: true,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		ReadTimeout:   5 * time.Second,
	}

	server := fiber.New(fiberConfig)
	server.Use(middleware.WithTimeout(5 * time.Second))

	instance = &Config{
		Server: server,
	}

	return instance
}

// Returns the current application configuration instance
func GetInstance() *Config {
	return instance
}
