package client

import "fmt"

type Paging struct {
	Skip  int
	Limit int
}

func (p Paging) String() string {
	if p.Limit <= 0 {
		if p.Skip > 0 {
			return fmt.Sprintf("OFFSET %d", p.Skip)
		}
		return ""
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.Limit, p.Skip)
}
