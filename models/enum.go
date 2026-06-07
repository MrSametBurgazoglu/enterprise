package models

import (
	"fmt"
	"reflect"
	"strings"
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

func (s *EnumDBField) GetSQLDefault() (string, bool) {
	if !s.HaveDefault {
		return "", false
	}
	if s.defaultFunc.IsValid() {
		return s.Field.GetSQLDefault()
	}
	return "'" + strings.ReplaceAll(s.DefaultValue, "'", "''") + "'", true
}

func EnumField(name string, values []string) *EnumDBField {
	f := &EnumDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.TypeName = name
	f.setField(name, f.TypeName, FieldTypeEnum)
	valuesArray := make([]*EnumValue, len(values))
	for i, v := range values {
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

func (s *EnumDBField) GoType(t reflect.Type) *EnumDBField {
	pkg := t.PkgPath()
	if pkg != "" {
		lastPackage := pkg[strings.LastIndex(pkg, "/")+1:]
		s.Field.Type = fmt.Sprintf("%s.%s", lastPackage, t.Name())
		s.Field.BaseType = s.Field.Type
		s.Field.RequiredPackages = append(s.Field.RequiredPackages, pkg)
	} else {
		s.Field.Type = t.Name()
		s.Field.BaseType = s.Field.Type
	}
	s.HaveCustomType = false // prevent generating Enum definition locally in schema struct template
	return s
}
