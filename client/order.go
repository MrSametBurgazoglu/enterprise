package client

import (
	"fmt"
	"strings"
)

type Order struct {
	Desc  bool
	Field string
}

func (o Order) String() string {
	order := "ASC"
	if o.Desc {
		order = "DESC"
	}
	field := ValidateIdentifier(o.Field)
	if !strings.Contains(field, "\"") && !strings.ContainsAny(field, "() ,") {
		if strings.Contains(field, ".") {
			parts := strings.Split(field, ".")
			for i, p := range parts {
				parts[i] = fmt.Sprintf(`"%s"`, p)
			}
			field = strings.Join(parts, ".")
		} else {
			field = fmt.Sprintf(`"%s"`, field)
		}
	}
	return fmt.Sprintf("%s %s", field, order)
}
