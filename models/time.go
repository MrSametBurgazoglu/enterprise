package models

import (
	"time"
)

type TimeDBField struct {
	*Field
	DefaultFuncStruct *FuncStruct
}

func (t *TimeDBField) DefaultFunc(v func() time.Time) *TimeDBField {
	t.DefaultFuncStruct.DefaultFunc(v)
	t.HaveDefault = true
	return t
}

func (t *TimeDBField) GetDefault() string {
	if t.DefaultFuncStruct.PackageFunc != "" {
		t.RequiredPackages = append(t.RequiredPackages, t.DefaultFuncStruct.PackageAddress)
		return t.DefaultFuncStruct.PackageFunc + "()"
	}
	return ""
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
