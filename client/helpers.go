package client

import "strings"

func withAndClause(whereStrings []string) string {
	if len(whereStrings) == 0 {
		return ""
	}
	return "(" + strings.Join(whereStrings, " AND ") + ")"
}
