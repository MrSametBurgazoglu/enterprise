package models

import (
	"fmt"
)

type ByteDBField struct {
	*Field
	DefaultValue      []byte
	DefaultFuncStruct *FuncStruct
}

func (i *ByteDBField) GetDefault() string {
	if i.defaultFunc.IsValid() {
		return i.Field.GetDefault()
	}
	return ""
}

func (i *ByteDBField) GetSQLDefault() (string, bool) {
	if !i.HaveDefault {
		return "", false
	}
	if i.defaultFunc.IsValid() {
		return i.Field.GetSQLDefault()
	}
	return fmt.Sprintf("'\\x%x'", i.DefaultValue), true
}

func (i *ByteDBField) Default(v []byte) *ByteDBField {
	i.DefaultValue = v
	i.HaveDefault = true
	return i
}

func (i *ByteDBField) DefaultFunc(v func() []byte) *ByteDBField {
	i.Field.DefaultFunc(v)
	return i
}

func ByteField(name string) *ByteDBField {
	f := &ByteDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "[]byte", FieldTypeByte)
	return f
}

func (i *ByteDBField) SetDBName(v string) *ByteDBField {
	i.DBName = v
	return i
}
