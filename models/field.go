package models

import (
	"strings"
)

const (
	FieldTypeBool = iota + 1
	FieldTypeEnum
	FieldTypeInt
	FieldTypeSmallInt
	FieldTypeBigInt
	FieldTypeFloat32
	FieldTypeFloat64
	FieldTypeString
	FieldTypeTime
	FieldTypeUUID
	FieldTypeUint
)

type FieldI interface {
	GetRequiredPackages() []string
	GetFieldType() int
	GetType() string
	GetName() string
	GetDBName() string
	IsNillable() bool
	IsSerial() bool
}

type Field struct {
	FieldType        int
	Name             string
	DBName           string
	Type             string
	BaseType         string
	Nillable         bool
	HaveDefault      bool
	IsPrepare        bool
	IsGreater        bool
	HaveCustomType   bool
	IsTime           bool
	IsUUID           bool
	Serial           bool
	RequiredPackages []string
}

func (f *Field) GetName() string {
	return f.Name
}

func (f *Field) GetNameLower() string {
	return strings.ToLower(f.Name)
}

func (f *Field) GetNameTitle() string {
	return strings.ToTitle(f.Name)
}

func (f *Field) GetDBName() string {
	return f.DBName
}

func (f *Field) SetDBName(v string) {
	f.DBName = v
}

func (f *Field) GetType() string {
	return f.Type
}

func (f *Field) GetBaseType() string {
	return f.BaseType
}

func (f *Field) GetAddressType() string {
	if f.Type[0] == '*' {
		return f.Type
	}
	return "*" + f.Type
}

func (f *Field) SetNillable() *Field {
	f.Nillable = true
	f.BaseType = f.Type
	f.Type = "*" + f.Type
	return f
}

func (f *Field) IsNillable() bool {
	return f.Nillable
}

func (f *Field) IsSerial() bool {
	return f.Serial
}

func (f *Field) IsDefault() bool {
	return f.HaveDefault
}

func (f *Field) NeedPrepare() bool {
	return f.IsPrepare
}

func (f *Field) CanBeGreater() bool {
	return f.IsGreater
}

func (f *Field) CanTime() bool {
	return f.IsTime
}

func (f *Field) CanUUID() bool {
	return f.IsUUID
}

func (f *Field) IsCustomType() bool {
	return f.HaveCustomType
}

func (f *Field) GetRequiredPackages() []string {
	return f.RequiredPackages
}

func (f *Field) GetFieldType() int {
	return f.FieldType
}

func (f *Field) setField(name, typeName string, fieldType int) {
	f.FieldType = fieldType
	f.Name = name
	f.Type = typeName
	f.BaseType = typeName
	f.SetDBNameManually(name)
}

func (f *Field) SetDBNameManually(name string) {
	snake := ConvertToSnakeCase(name)
	f.SetDBName(snake)
}
