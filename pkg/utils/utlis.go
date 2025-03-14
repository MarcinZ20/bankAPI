package utils

import (
	"regexp"
	"strings"
)

func IsUppercase(s string) bool {
	return s == strings.ToUpper(s)
}

func IsLowercase(s string) bool {
	return s == strings.ToLower(s)
}

// SWIFT code validation was done via regex according to the rules stated here:
// https://www.geeksforgeeks.org/how-to-validate-swift-bic-code-using-regex/
func IsValidSwiftCodeFormat(swiftCode string) bool {
	swiftRegex := `^[A-Z]{6}[A-Z0-9]{2}([A-Z0-9]{3})?$`
	matched, _ := regexp.MatchString(swiftRegex, swiftCode)
	return matched
}

func IsNotEmpty(s string) bool {
	return s != ""
}

func IsValidCountryCode(code string) bool {
	countryISO2Regex := `^[A-Z]{2}$`
	match, _ := regexp.MatchString(countryISO2Regex, code)
	return match
}
