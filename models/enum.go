package models

import (
	"fmt"
)

type EnumValue struct {
	value string
}

func (e *EnumValue) Title() string {
	return ToCamelCase(e.value)
}


func (e *EnumValue) Value() string {
	return e.value
}

type EnumDBField struct {
	*Field
	DefaultValue      string
	DefaultFuncStruct *FuncStruct
	TypeName          string
	Values            []*EnumValue
}

func (s *EnumDBField) Default(v string) *EnumDBField {
	s.DefaultValue = v
	s.HaveDefault = true
	return s
}

func (s *EnumDBField) DefaultFunc(v func() string) *EnumDBField {
	s.Field.DefaultFunc(v)
	return s
}

func (s *EnumDBField) GetDefault() string {
	if s.defaultFunc.IsValid() {
		return s.Field.GetDefault()
	} else {
		return fmt.Sprintf("\"%s\"", s.DefaultValue)
	}
}

func EnumField(name string, values []string) *EnumDBField {
	f := &EnumDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.TypeName = name
	f.setField(name, f.TypeName, FieldTypeEnum)
	valuesArray := make([]*EnumValue, len(values))
	for i, v := range values{
		valuesArray[i] = &EnumValue{value: v}
	}
	f.Values = valuesArray
	f.HaveCustomType = true
	f.CanIn = true
	return f
}

func (s *EnumDBField) SetDBName(v string) *EnumDBField {
	s.DBName = v
	return s
}
