package client

import "fmt"

type Order struct {
	Desc  bool
	Field string
}

func (o Order) String() string {
	order := "ASC"
	if o.Desc {
		order = "DESC"
	}
	return fmt.Sprintf("ORDER BY %s %s", o.Field, order)
}
