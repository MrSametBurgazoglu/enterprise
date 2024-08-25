package models

import (
	"fmt"
	"strings"
)

type Generation struct {
	Tables []*Table
}

const (
	RelationTypeManyToOne = iota
	RelationTypeOneToMany
	RelationTypeManyToMany
)

type Relation struct {
	RelationTable       string
	RelationTableDBName string
	RelationField       string
	OnField             string
	ManyTableDBName     string
	RelationTableField  string
	RelationType        int
}

func (r *Relation) IsRelationList() bool {
	if r.RelationType == RelationTypeManyToOne {
		return false
	} else {
		return true
	}
}

func (r *Relation) GetRelationField() string {
	if r.RelationType == RelationTypeManyToOne {
		return fmt.Sprintf("%s", r.RelationTable)
	} else {
		return fmt.Sprintf("%sList", r.RelationTable)
	}
}

func (r *Relation) GetRelationTableLower() string {
	return strings.ToLower(r.RelationTable)
}

func (r *Relation) IsManyToMany() bool {
	if r.RelationType == RelationTypeManyToMany {
		return true
	}
	return false
}

func OneToMany(table, toField, fromField string) *Relation {
	r := &Relation{RelationTable: table, RelationField: fromField, OnField: toField, RelationType: RelationTypeOneToMany}
	r.RelationTableDBName = ConvertToSnakeCase(r.RelationTable)
	return r
}

func ManyToOne(table, toField, fromField string) *Relation {
	r := &Relation{RelationTable: table, RelationField: toField, OnField: fromField, RelationType: RelationTypeManyToOne}
	r.RelationTableDBName = ConvertToSnakeCase(r.RelationTable)
	return r
}

func ManyToMany(table, toField, fromField, relationField, manyTableName string) *Relation {
	r := &Relation{RelationTable: table, RelationField: toField, OnField: fromField, RelationTableField: relationField, ManyTableDBName: manyTableName, RelationType: RelationTypeManyToMany}
	r.RelationTableDBName = ConvertToSnakeCase(r.RelationTable)
	return r
}

type Table struct {
	PackageName      string
	TableName        string
	DBName           string
	IDField          string
	IDDBField        string
	IDFieldType      string
	Fields           []FieldI
	Relations        []*Relation
	RequiredPackages []string
}

func (t *Table) SetTableName(name string) {
	t.TableName = name
	t.DBName = ConvertToSnakeCase(t.TableName)
}

func (t *Table) SetIDField(field FieldI) {
	t.IDField = field.GetName()
	t.IDDBField = field.GetDBName()
	t.IDFieldType = field.GetType()
	t.DBName = ConvertToSnakeCase(t.TableName)
}

func (t *Table) IDFieldLower() string {
	return strings.ToLower(t.IDField)
}
