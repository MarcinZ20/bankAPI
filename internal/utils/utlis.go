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

func IsValidSwiftCodeFormat(swiftCode string) bool {
	swiftRegex := `^[A-Z0-9]{8}(XXX|[A-Z0-9]{3})?$`
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
