package services

import (
	"github.com/MarcinZ20/bankAPI/internal/database"
)

// Handles all services in the application
type ServiceManager struct {
	BankService *BankService
}

var instance *ServiceManager

// Creates a new service manager with all services initialized
func NewServiceManager(db *database.Config) *ServiceManager {
	if instance != nil {
		return instance
	}

	instance = &ServiceManager{
		BankService: NewBankService(db.Collection),
	}

	return instance
}

// Returns the current service manager instance
func GetInstance() *ServiceManager {
	return instance
}

// Checks if all services are properly initialized
func (sm *ServiceManager) IsInitialized() bool {
	return sm != nil && sm.BankService != nil && sm.BankService.IsInitialized()
}
