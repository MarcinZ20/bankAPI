package validation

import (
	"fmt"
	"strings"

	"github.com/MarcinZ20/bankAPI/pkg/models"
	"github.com/MarcinZ20/bankAPI/pkg/utils"
)

type ValidationResult struct {
	IsValid bool
	Errors  []string
}

type Validator interface {
	Validate(value string) []error
}

type SwiftCodeValidator struct{}

// Validates SWIFT code integrity
func (v SwiftCodeValidator) Validate(value string) []error {
	length := len(value)
	errors := []error{}

	if length < 8 || length > 11 {
		errors = append(errors, fmt.Errorf("invalid swift code: length must be between 8 and 11 characters long, but is %v", length))
	}

	if !utils.IsValidSwiftCodeFormat(value) {
		errors = append(errors, fmt.Errorf("invalid swift code: %v does not match the expected format", value))
	}

	return errors
}

type CountryISO2Validator struct{}

// Validates ISO2 country code integrity
func (v CountryISO2Validator) Validate(value string) []error {
	length := len(value)
	errors := []error{}

	if length != 2 {
		errors = append(errors, fmt.Errorf("invalid countryISO2 code: length must be 2 characters long, but is %v", length))
	}

	if strings.ContainsAny(value, "0123456789") {
		errors = append(errors, fmt.Errorf("invalid countryISO2 code: cannot contain numbers"))
	}

	return errors
}

// Validates BankEntity object by checking SWIFT code and ISO2 code
func ValidateBankEntity(entity any) ValidationResult {
	result := ValidationResult{
		IsValid: true,
	}

	switch e := entity.(type) {
	case models.Bank:
		validateField("SwiftCode", e.SwiftCode, []Validator{SwiftCodeValidator{}}).appendErrors(&result)
		validateField("CountryISO2Code", e.CountryISO2Code, []Validator{CountryISO2Validator{}}).appendErrors(&result)
	default:
		result.IsValid = false
		result.Errors = append(result.Errors, "unsupported input type")
	}

	return result
}

// Validates single field by running propper checks
func validateField(fieldName string, value string, validators []Validator) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	for _, v := range validators {
		validationResult := v.Validate(value)
		if len(validationResult) != 0 {
			result.IsValid = false
			for _, err := range validationResult {
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %s", fieldName, err))
			}
		}
	}

	return result
}

// Appends thrown errors into a ValidationResult
func (v *ValidationResult) appendErrors(result *ValidationResult) {
	if !v.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, v.Errors...)
	}
}
