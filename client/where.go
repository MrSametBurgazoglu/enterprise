package client

import "fmt"

const (
	EQ   = "\"%s\".\"%s\" = %s"
	NEQ  = "\"%s\".\"%s\" != %s"
	GT   = "\"%s\".\"%s\" > %s"
	GTE  = "\"%s\".\"%s\" >= %s"
	LT   = "\"%s\".\"%s\" < %s"
	LTE  = "\"%s\".\"%s\" <= %s"
	NIL  = "\"%s\".\"%s\" IS NIL"
	NNIL = "\"%s\".\"%s\" IS NOT NIL"
	ANY  = "\"%s\".\"%s\" = ANY(%s)"
	NANY = "\"%s\".\"%s\" != ANY(%s)"
)

type Where struct {
	Type     string
	Name     string
	HasValue bool
	Value    any
}

func (w *Where) GetSqlString(tableName string) string {
	if w.HasValue {
		return fmt.Sprintf(w.Type, tableName, w.Name, "@"+w.Name)
	} else {
		return fmt.Sprintf(w.Type, tableName, w.Name)
	}
}
