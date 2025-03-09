package services

import (
	"context"
	"fmt"

	"github.com/MarcinZ20/bankAPI/pkg/models"
	"github.com/MarcinZ20/bankAPI/pkg/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

// BankService handles business logic for bank operations
type BankService struct {
	repo *repository.BankRepository
}

// NewBankService creates a new bank service
func NewBankService(collection *mongo.Collection) *BankService {
	return &BankService{
		repo: repository.NewBankRepository(collection),
	}
}

// GetHeadquarter retrieves a headquarter by SWIFT code
func (s *BankService) GetHeadquarter(ctx context.Context, swiftCode string) (*models.Headquarter, error) {
	return s.repo.FindHeadquarter(ctx, swiftCode)
}

// GetBranch retrieves a branch by SWIFT code
func (s *BankService) GetBranch(ctx context.Context, swiftCode string) (*models.Branch, error) {
	parentHqSwiftCode := swiftCode[0:8] + "XXX"
	return s.repo.FindBranch(ctx, swiftCode, parentHqSwiftCode)
}

// GetBanksByCountryCode retrieves all banks in a given country
func (s *BankService) GetBanksByCountryCode(ctx context.Context, countryCode string) ([]models.Headquarter, error) {
	return s.repo.FindBanksByCountry(ctx, countryCode)
}

// AddHeadquarter creates a new headquarter
func (s *BankService) AddHeadquarter(ctx context.Context, hq *models.Headquarter) error {
	// Validate headquarter data
	if hq.SwiftCode == "" || hq.BankName == "" || hq.CountryISO2 == "" {
		return fmt.Errorf("missing required fields")
	}

	return s.repo.CreateHeadquarter(ctx, hq)
}

// AddBranch adds a new branch to a headquarter
func (s *BankService) AddBranch(ctx context.Context, parentSwiftCode string, branch *models.Branch) error {
	// Validate branch data
	if branch.SwiftCode == "" || branch.BankName == "" || branch.CountryISO2 == "" {
		return fmt.Errorf("missing required fields")
	}

	return s.repo.AddBranch(ctx, parentSwiftCode, branch)
}

// DeleteHeadquarter deletes a headquarter and all its branches
func (s *BankService) DeleteHeadquarter(ctx context.Context, swiftCode string) error {
	return s.repo.DeleteHeadquarter(ctx, swiftCode)
}

// DeleteBranch removes a branch from its headquarter
func (s *BankService) DeleteBranch(ctx context.Context, swiftCode string, parentSwiftCode string) error {
	return s.repo.DeleteBranch(ctx, swiftCode, parentSwiftCode)
}
