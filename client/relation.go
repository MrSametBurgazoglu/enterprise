package client

import (
	"fmt"
	"strings"
)

type RelationJoinType string

const (
	RelationJoinTypeDefault RelationJoinType = ""
	RelationJoinTypeLeft    RelationJoinType = "LEFT"
	RelationJoinTypeRight   RelationJoinType = "RIGHT"
	RelationJoinTypeFull    RelationJoinType = "FULL"
	RelationJoinTypeInner   RelationJoinType = "INNER"
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
	RelationModel    Model
	RelationResult   Result
	RelationTable    string
	RelationWhere    *RelationCondition
	RelationJoinType RelationJoinType
	Where            []*WhereList
	ManyToManyTable  string
}

func (r *Relation) SetJoinType(t RelationJoinType) *Relation {
	r.RelationJoinType = t
	return r
}

func (r *Relation) getJoinType() string {
	if r.RelationJoinType == RelationJoinTypeDefault {
		return string(RelationJoinTypeLeft)
	}
	return string(r.RelationJoinType)
}

func (r *Relation) getJoinString(tableName string) string {
	if r.ManyToManyTable != "" {
		return fmt.Sprintf(
			"%s JOIN \"%s\" ON \"%s\".\"%s\" = \"%s\".\"%s\" LEFT JOIN \"%s\" ON \"%s\".\"%s\" = \"%s\".\"%s\" ",
			r.getJoinType(),
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
	return fmt.Sprintf("%s JOIN \"%s\" ON %s", r.getJoinType(), r.RelationTable, r.RelationWhere.String(tableName, r.RelationTable))
}

func (r *Relation) isRelationHaveWhereClause() bool {
	return len(r.Where) != 0
}

func (r *Relation) parseWhere() *Res {
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
