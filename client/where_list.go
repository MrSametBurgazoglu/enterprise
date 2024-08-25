package client

type WhereList struct {
	Items []*Where
}

func (w *WhereList) Parse(tableName string) *Res {
	res := new(Res)
	var whereStrings []string
	for _, item := range w.Items {
		whereStrings = append(whereStrings, item.GetSqlString(tableName))
		if !item.HasValue {
			continue
		}
		res.Values = append(res.Values, item.Value)
		res.Names = append(res.Names, item.Name)
	}
	sql := withAndClause(whereStrings)
	res.SqlString = sql
	return res
}
