package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/MarcinZ20/bankAPI/internal/repository"
	"github.com/MarcinZ20/bankAPI/pkg/models"
	"github.com/MarcinZ20/bankAPI/pkg/utils"
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

func (s *BankService) IsInitialized() bool {
	return s != nil && s.repo != nil
}

// GetHeadquarter retrieves a headquarter by SWIFT code
func (s *BankService) GetHeadquarter(ctx context.Context, swiftCode string) (*models.Headquarter, error) {
	if !utils.IsValidSwiftCodeFormat(swiftCode) {
		return nil, fmt.Errorf("invalid SWIFT code format")
	}
	if !strings.HasSuffix(swiftCode, "XXX") {
		return nil, fmt.Errorf("SWIFT code must end with XXX for headquarters")
	}
	return s.repo.FindHeadquarter(ctx, swiftCode)
}

// GetBranch retrieves a branch by SWIFT code
func (s *BankService) GetBranch(ctx context.Context, swiftCode string) (*models.Branch, error) {
	if !utils.IsValidSwiftCodeFormat(swiftCode) {
		return nil, fmt.Errorf("invalid SWIFT code format")
	}
	if strings.HasSuffix(swiftCode, "XXX") {
		return nil, fmt.Errorf("branch SWIFT code cannot end with XXX")
	}
	parentHqSwiftCode := swiftCode[0:8] + "XXX"
	return s.repo.FindBranch(ctx, swiftCode, parentHqSwiftCode)
}

// GetBanksByCountryCode retrieves all banks in a given country
func (s *BankService) GetBanksByCountryCode(ctx context.Context, countryCode string) ([]models.Headquarter, error) {
	if !utils.IsValidCountryCode(countryCode) {
		return nil, fmt.Errorf("invalid country code format")
	}
	return s.repo.FindBanksByCountry(ctx, countryCode)
}

// AddHeadquarter creates a new headquarter
func (s *BankService) AddHeadquarter(ctx context.Context, hq *models.Headquarter) error {
	if err := s.validateHeadquarter(hq); err != nil {
		return err
	}
	return s.repo.CreateHeadquarter(ctx, hq)
}

// AddBranch adds a new branch to a headquarter
func (s *BankService) AddBranch(ctx context.Context, parentSwiftCode string, branch *models.Branch) error {
	if err := s.validateBranch(branch); err != nil {
		return err
	}
	if !strings.HasSuffix(parentSwiftCode, "XXX") {
		return fmt.Errorf("parent SWIFT code must end with XXX")
	}
	return s.repo.AddBranch(ctx, parentSwiftCode, branch)
}

// DeleteHeadquarter deletes a headquarter and all its branches
func (s *BankService) DeleteHeadquarter(ctx context.Context, swiftCode string) error {
	if !utils.IsValidSwiftCodeFormat(swiftCode) {
		return fmt.Errorf("invalid SWIFT code format")
	}
	if !strings.HasSuffix(swiftCode, "XXX") {
		return fmt.Errorf("SWIFT code must end with XXX for headquarters")
	}
	return s.repo.DeleteHeadquarter(ctx, swiftCode)
}

// DeleteBranch removes a branch from its headquarter
func (s *BankService) DeleteBranch(ctx context.Context, swiftCode, parentSwiftCode string) error {
	if !utils.IsValidSwiftCodeFormat(swiftCode) || !utils.IsValidSwiftCodeFormat(parentSwiftCode) {
		return fmt.Errorf("invalid SWIFT code format")
	}
	if strings.HasSuffix(swiftCode, "XXX") {
		return fmt.Errorf("branch SWIFT code cannot end with XXX")
	}
	if !strings.HasSuffix(parentSwiftCode, "XXX") {
		return fmt.Errorf("parent SWIFT code must end with XXX")
	}
	return s.repo.DeleteBranch(ctx, swiftCode, parentSwiftCode)
}

// validateHeadquarter validates headquarter data
func (s *BankService) validateHeadquarter(hq *models.Headquarter) error {
	if hq == nil {
		return fmt.Errorf("headquarter cannot be nil")
	}
	if !utils.IsValidSwiftCodeFormat(hq.SwiftCode) {
		return fmt.Errorf("invalid SWIFT code format")
	}
	if !strings.HasSuffix(hq.SwiftCode, "XXX") {
		return fmt.Errorf("headquarter SWIFT code must end with XXX")
	}
	if hq.BankName == "" {
		return fmt.Errorf("bank name is required")
	}
	if !utils.IsValidCountryCode(hq.CountryISO2) {
		return fmt.Errorf("invalid country code format")
	}
	if !hq.IsHeadquarter {
		return fmt.Errorf("isHeadquarter must be true")
	}
	return nil
}

// validateBranch validates branch data
func (s *BankService) validateBranch(branch *models.Branch) error {
	if branch == nil {
		return fmt.Errorf("branch cannot be nil")
	}
	if !utils.IsValidSwiftCodeFormat(branch.SwiftCode) {
		return fmt.Errorf("invalid SWIFT code format")
	}
	if strings.HasSuffix(branch.SwiftCode, "XXX") {
		return fmt.Errorf("branch SWIFT code cannot end with XXX")
	}
	if branch.BankName == "" {
		return fmt.Errorf("bank name is required")
	}
	if !utils.IsValidCountryCode(branch.CountryISO2) {
		return fmt.Errorf("invalid country code format")
	}
	if branch.IsHeadquarter {
		return fmt.Errorf("isHeadquarter must be false")
	}
	return nil
}
