package services

import (
	"github.com/MarcinZ20/bankAPI/internal/database"
)

// ServiceManager handles all services in the application
type ServiceManager struct {
	BankService *BankService
}

var instance *ServiceManager

// NewServiceManager creates a new service manager with all services initialized
func NewServiceManager(db *database.Config) *ServiceManager {
	if instance != nil {
		return instance
	}

	instance = &ServiceManager{
		BankService: NewBankService(db.Collection),
	}

	return instance
}

// GetInstance returns the current service manager instance
func GetInstance() *ServiceManager {
	return instance
}

// IsInitialized checks if all services are properly initialized
func (sm *ServiceManager) IsInitialized() bool {
	return sm != nil && sm.BankService != nil && sm.BankService.IsInitialized()
}
