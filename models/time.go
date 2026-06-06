package models

import (
	"time"
)

type TimeDBField struct {
	*Field
	DefaultFuncStruct *FuncStruct
}

func (t *TimeDBField) DefaultFunc(v func() time.Time) *TimeDBField {
	t.Field.DefaultFunc(v)
	return t
}

func (t *TimeDBField) GetDefault() string {
	if t.defaultFunc.IsValid() {
		return t.Field.GetDefault()
	}
	return ""
}

func (t *TimeDBField) GetSQLDefault() (string, bool) {
	if !t.HaveDefault {
		return "", false
	}
	return t.Field.GetSQLDefault()
}

func (t *TimeDBField) PrepareFunc() string {
	return "new(time.Time)"
}

func TimeField(name string) *TimeDBField {
	f := &TimeDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "time.Time", FieldTypeTime)
	f.IsTime = true
	f.IsGreater = true
	f.RequiredPackages = append(f.RequiredPackages, "time")
	return f
}

func (t *TimeDBField) SetDBName(v string) *TimeDBField {
	t.DBName = v
	return t
}

func (t *TimeDBField) SetNillable() *TimeDBField {
	t.Field.SetNillable()
	t.IsPrepare = true
	return t
}
