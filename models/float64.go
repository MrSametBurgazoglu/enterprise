package models

import "strconv"

type Float64DBField struct {
	*Field
	DefaultValue      int
	DefaultFuncStruct *FuncStruct
}

func (i *Float64DBField) GetDefault() string {
	if i.DefaultFuncStruct.PackageFunc != "" {
		i.RequiredPackages = append(i.RequiredPackages, i.DefaultFuncStruct.PackageAddress)
		return i.DefaultFuncStruct.PackageFunc + "()"
	} else {
		return strconv.Itoa(i.DefaultValue)
	}
}

func (i *Float64DBField) Default(v int) *Float64DBField {
	i.DefaultValue = v
	i.HaveDefault = true
	return i
}

func (i *Float64DBField) DefaultFunc(v func() int) *Float64DBField {
	i.DefaultFuncStruct.DefaultFunc(v)
	i.RequiredPackages = append(i.RequiredPackages, i.DefaultFuncStruct.PackageAddress)
	i.HaveDefault = true
	return i
}

func Float64Field(name string) *Float64DBField {
	f := &Float64DBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "float64", FieldTypeFloat64)
	f.IsGreater = true
	return f
}

func (i *Float64DBField) SetDBName(v string) *Float64DBField {
	i.DBName = v
	return i
}
