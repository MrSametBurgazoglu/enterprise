package models

import (
	"github.com/google/uuid"
)

type UUIDDBField struct {
	*Field
	DefaultFuncStruct *FuncStruct
}

func (u *UUIDDBField) DefaultFunc(v func() uuid.UUID) *UUIDDBField {
	u.DefaultFuncStruct.DefaultFunc(v)
	u.HaveDefault = true
	return u
}

func (u *UUIDDBField) GetDefault() string {
	if u.DefaultFuncStruct.PackageFunc != "" {
		u.RequiredPackages = append(u.RequiredPackages, u.DefaultFuncStruct.PackageAddress)
		return u.DefaultFuncStruct.PackageFunc + "()"
	}
	return ""
}

func (u *UUIDDBField) PrepareFunc() string {
	return "&uuid.Nil"
}

func UUIDField(name string) *UUIDDBField {
	f := &UUIDDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "uuid.UUID", FieldTypeUUID)
	f.IsUUID = true
	f.RequiredPackages = append(f.RequiredPackages, "github.com/google/uuid")
	return f
}

func (u *UUIDDBField) SetDBName(v string) *UUIDDBField {
	u.DBName = v
	return u
}

func (u *UUIDDBField) SetNillable() *UUIDDBField {
	u.Field.SetNillable()
	u.IsPrepare = true
	return u
}
