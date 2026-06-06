package models

import (
	"strconv"
	"strings"
)

type StringDBField struct {
	*Field
	DefaultValue      string
	DefaultFuncStruct *FuncStruct
}

func (s *StringDBField) Default(v string) *StringDBField {
	s.DefaultValue = v
	s.HaveDefault = true
	return s
}

func (s *StringDBField) DefaultFunc(v func() string) *StringDBField {
	s.Field.DefaultFunc(v)
	return s
}

func (s *StringDBField) GetDefault() string {
	if s.defaultFunc.IsValid() {
		return s.Field.GetDefault()
	}
	return strconv.Quote(s.DefaultValue)
}

func (s *StringDBField) GetSQLDefault() (string, bool) {
	if !s.HaveDefault {
		return "", false
	}
	if s.defaultFunc.IsValid() {
		return s.Field.GetSQLDefault()
	}
	return "'" + strings.ReplaceAll(s.DefaultValue, "'", "''") + "'", true
}

func StringField(name string) *StringDBField {
	f := &StringDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "string", FieldTypeString)
	f.CanIn = true
	return f
}

func (s *StringDBField) SetDBName(v string) *StringDBField {
	s.DBName = v
	return s
}
