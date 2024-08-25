package client

import (
	"fmt"
	"strings"
)

type RelationModel interface {
	GetDBName() string
	IsExist() bool
}

type RelationCondition struct {
	RelationValue      string
	TableValue         string
	RelationTableValue string
}

func (r *RelationCondition) String(tableName, relationTableName string) string {
	return fmt.Sprintf("\"%s\".\"%s\" = \"%s\".\"%s\"", tableName, r.TableValue, relationTableName, r.RelationValue)
}

type Relation struct {
	RelationModel   Model
	RelationResult  Result
	RelationTable   string
	RelationWhere   *RelationCondition
	Where           []*WhereList
	ManyToManyTable string
}

func (r *Relation) GetJoinString(tableName string) string {
	if r.ManyToManyTable != "" {
		return fmt.Sprintf(
			"LEFT JOIN \"%s\" ON \"%s\".\"%s\" = \"%s\".\"%s\" LEFT JOIN \"%s\" ON \"%s\".\"%s\" = \"%s\".\"%s\" ",
			r.ManyToManyTable,
			tableName,
			r.RelationWhere.RelationTableValue, //todo get two tables primary key name
			r.ManyToManyTable,
			r.RelationWhere.RelationValue,
			r.RelationTable,
			r.RelationTable,
			r.RelationWhere.RelationTableValue,
			r.ManyToManyTable,
			r.RelationWhere.TableValue,
		)
	}
	return fmt.Sprintf("LEFT JOIN \"%s\" ON %s", r.RelationTable, r.RelationWhere.String(tableName, r.RelationTable))
}

func (r *Relation) IsRelationHaveWhereClause() bool {
	return len(r.Where) != 0
}

func (r *Relation) ParseWhere() *Res {
	res := new(Res)
	var whereStrings []string
	for _, list := range r.Where {
		resp := list.Parse(r.RelationTable)
		whereStrings = append(whereStrings, resp.SqlString)
		res.Names = append(res.Names, resp.Names...)
		res.Values = append(res.Values, resp.Values...)
	}
	res.SqlString = strings.Join(whereStrings, " OR ")
	return res
}

type RelationList struct {
	Relations   []*Relation
	RelationMap map[string]*Relation
}
