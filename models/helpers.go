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

func ToCamelCase(s string) string {
	delimiters := []string{"_", "-", " "}
	for _, d := range delimiters {
		s = strings.ReplaceAll(s, d, " ")
	}

	words := strings.Fields(s)
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, "")
}
