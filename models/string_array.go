package models

import (
	"strconv"
	"strings"
)

type StringArrayDBField struct {
	*Field
	DefaultValue      []string
	DefaultFuncStruct *FuncStruct
}

func StringArrayField(name string) *StringArrayDBField {
	f := &StringArrayDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "[]string", FieldTypeStringArray)
	return f
}

func (s *StringArrayDBField) SetDBName(v string) *StringArrayDBField { s.DBName = v; return s }

func (s *StringArrayDBField) Default(v []string) *StringArrayDBField {
	s.DefaultValue = v
	s.HaveDefault = true
	return s
}

func (s *StringArrayDBField) DefaultFunc(v func() []string) *StringArrayDBField {
	s.Field.DefaultFunc(v)
	return s
}

func (s *StringArrayDBField) GetDefault() string {
	if s.defaultFunc.IsValid() {
		return s.Field.GetDefault()
	}
	if s.DefaultValue == nil {
		return "nil"
	}
	quoted := make([]string, len(s.DefaultValue))
	for i, val := range s.DefaultValue {
		quoted[i] = strconv.Quote(val)
	}
	return "[]string{" + strings.Join(quoted, ", ") + "}"
}

func (s *StringArrayDBField) GetSQLDefault() (string, bool) {
	if !s.HaveDefault {
		return "", false
	}
	if s.defaultFunc.IsValid() {
		return s.Field.GetSQLDefault()
	}
	
	var builder strings.Builder
	builder.WriteString("'{")
	for i, val := range s.DefaultValue {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteByte('"')
		r := strings.NewReplacer(`\`, `\\`, `"`, `\"`)
		builder.WriteString(r.Replace(val))
		builder.WriteByte('"')
	}
	builder.WriteString("}'")
	return builder.String(), true
}
