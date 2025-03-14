package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwiftCodeValidator(t *testing.T) {
	validator := SwiftCodeValidator{}

	tests := []struct {
		name     string
		input    string
		wantErr  bool
		errCount int
	}{
		{
			name:     "SWIFT code: valid 8-character",
			input:    "DEUTDEFF",
			wantErr:  false,
			errCount: 0,
		},
		{
			name:     "SWIFT code: valid 11-character",
			input:    "DEUTDEFF100",
			wantErr:  false,
			errCount: 0,
		},
		{
			name:     "SWIFT code: too short",
			input:    "DEUT",
			wantErr:  true,
			errCount: 2, // Length and format are wrong
		},
		{
			name:     "SWIFT code: too long",
			input:    "DEUTDEFF1000",
			wantErr:  true,
			errCount: 2, // Length and format are wrong
		},
		{
			name:     "SWIFT code: invalid format",
			input:    "1234ABCD",
			wantErr:  true,
			errCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.Validate(tt.input)
			if tt.wantErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors\nRecieved: %d", tt.errCount, len(errors))
			} else {
				assert.Empty(t, errors, "Expected no errors\nRecieved: %v", errors)
			}
		})
	}
}

func TestCountryISO2Validator(t *testing.T) {
	validator := CountryISO2Validator{}

	tests := []struct {
		name     string
		input    string
		wantErr  bool
		errCount int
	}{
		{
			name:     "ISO2 code: valid",
			input:    "DE",
			wantErr:  false,
			errCount: 0,
		},
		{
			name:     "ISO2 code: too short",
			input:    "D",
			wantErr:  true,
			errCount: 1,
		},
		{
			name:     "ISO2 code: too long",
			input:    "DEU",
			wantErr:  true,
			errCount: 1,
		},
		{
			name:     "ISO2 code: contains numbers",
			input:    "D2",
			wantErr:  true,
			errCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.Validate(tt.input)
			if tt.wantErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors\nRecieved: %d", tt.errCount, len(errors))
			} else {
				assert.Empty(t, errors, "Expected no errors\nRecieved: %v", errors)
			}
		})
	}
}
