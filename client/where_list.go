package client

import (
	"github.com/jackc/pgx/v5"
)

type WhereList struct {
	Items []PredicateI
}

func (w *WhereList) Parse(tableName string, args pgx.NamedArgs) *Res {
	res := new(Res)
	var whereStrings []string
	for _, item := range w.Items {
		s := item.Parse(tableName, args)
		if s != "" {
			whereStrings = append(whereStrings, s)
		}
	}
	sql := withAndClause(whereStrings)
	res.SqlString = sql
	return res
}
