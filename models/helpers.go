package models

import (
	"regexp"
	"strings"
)

var caseConverterRegex = regexp.MustCompile("([a-z])([A-Z])")

func ConvertToSnakeCase(input string) string {
	// Use regular expression to find the positions where the uppercase letters are
	snake := caseConverterRegex.ReplaceAllString(input, "${1}_${2}")

	// Convert the whole string to lowercase
	return strings.ToLower(snake)
}
