package utils

import (
	"regexp"
	"strings"
)

// checks if a string is uppercase
func IsUppercase(s string) bool {
	return s == strings.ToUpper(s)
}

// checks if a string is lowercase
func IsLowercase(s string) bool {
	return s == strings.ToLower(s)
}

// checks if a swift code is valid
func IsValidSwiftCodeFormat(swiftCode string) bool {
	swiftRegex := `^[A-Z0-9]{8}(XXX|[A-Z0-9]{3})?$`
	matched, _ := regexp.MatchString(swiftRegex, swiftCode)
	return matched
}

// checks if a string is not empty
func IsNotEmpty(s string) bool {
	return s != ""
}

func IsValidCountryCode(code string) bool {
	countryISO2Regex := `^[A-Z]{2}$`
	match, _ := regexp.MatchString(countryISO2Regex, code)
	return match
}
