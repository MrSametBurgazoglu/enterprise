package client

import "fmt"

type Paging struct {
	Skip  int
	Limit int
}

func (p Paging) String() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.Limit, p.Skip)
}
