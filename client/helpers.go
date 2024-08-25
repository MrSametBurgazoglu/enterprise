package client

import "strings"

func withAndClause(whereStrings []string) string {
	return "(" + strings.Join(whereStrings, " AND ") + ")"
}
