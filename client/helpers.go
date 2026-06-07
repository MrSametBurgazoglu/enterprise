package client

import (
	"fmt"
	"regexp"
	"strings"
)

var identRegexp = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

var castWhitelist = map[string]bool{
	"text":              true,
	"varchar":           true,
	"character varying": true,
	"char":              true,
	"numeric":           true,
	"decimal":           true,
	"integer":           true,
	"int":               true,
	"bigint":            true,
	"smallint":          true,
	"boolean":           true,
	"bool":              true,
	"uuid":              true,
	"timestamp":         true,
	"timestamptz":       true,
	"date":              true,
	"time":              true,
	"json":              true,
	"jsonb":             true,
	"double precision":  true,
	"real":              true,
}

var castRegexp = regexp.MustCompile(`^[a-z0-9_]+(?:\s+[a-z0-9_]+)*(?:\s*\(\s*\d+\s*(?:,\s*\d+\s*)*\))?$`)

func IsValidIdentifier(s string) bool {
	if strings.HasPrefix(s, "RAW:") {
		return true
	}
	s = strings.Trim(s, `"`)
	if s == "" {
		return false
	}
	if strings.Contains(s, ".") {
		parts := strings.Split(s, ".")
		for _, part := range parts {
			if part != "*" && !identRegexp.MatchString(strings.ToLower(part)) {
				return false
			}
		}
		return true
	}
	return s == "*" || identRegexp.MatchString(strings.ToLower(s))
}

func IsValidCast(s string) bool {
	if s == "" {
		return true
	}
	if strings.HasPrefix(s, "RAW:") {
		return true
	}
	lower := strings.ToLower(s)
	if castWhitelist[lower] {
		return true
	}
	return castRegexp.MatchString(lower)
}

func ValidateIdentifier(s string) string {
	if strings.HasPrefix(s, "RAW:") {
		return s[4:]
	}
	if !IsValidIdentifier(s) {
		panic(fmt.Sprintf("enterprise ORM: invalid identifier %q", s))
	}
	return s
}

func ValidateCast(s string) string {
	if strings.HasPrefix(s, "RAW:") {
		return s[4:]
	}
	if !IsValidCast(s) {
		panic(fmt.Sprintf("enterprise ORM: invalid cast type %q", s))
	}
	return s
}

func Raw(s string) string {
	return "RAW:" + s
}

func Expr(s string) string {
	return "RAW:" + s
}

func withAndClause(whereStrings []string) string {
	if len(whereStrings) == 0 {
		return ""
	}
	return "(" + strings.Join(whereStrings, " AND ") + ")"
}
