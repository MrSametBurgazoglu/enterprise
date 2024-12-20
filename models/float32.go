package models

import "strconv"

type Float32DBField struct {
	*Field
	DefaultValue      float32
	DefaultFuncStruct *FuncStruct
}

func (i *Float32DBField) GetDefault() string {
	if i.DefaultFuncStruct.PackageFunc != "" {
		i.RequiredPackages = append(i.RequiredPackages, i.DefaultFuncStruct.PackageAddress)
		return i.DefaultFuncStruct.PackageFunc + "()"
	} else {
		return strconv.FormatFloat(float64(i.DefaultValue), 'g', -1, 64)
	}
}

func (i *Float32DBField) Default(v float32) *Float32DBField {
	i.DefaultValue = v
	i.HaveDefault = true
	return i
}

func (i *Float32DBField) DefaultFunc(v func() float32) *Float32DBField {
	i.DefaultFuncStruct.DefaultFunc(v)
	i.RequiredPackages = append(i.RequiredPackages, i.DefaultFuncStruct.PackageAddress)
	i.HaveDefault = true
	return i
}

func Float32Field(name string) *Float32DBField {
	f := &Float32DBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "float32", FieldTypeFloat32)
	f.IsGreater = true
	f.CanIn = true
	return f
}

func (i *Float32DBField) SetDBName(v string) *Float32DBField {
	i.DBName = v
	return i
}
