package models

import "strconv"

type IntDBField struct {
	*Field
	DefaultValue      int
	DefaultFuncStruct *FuncStruct
}

func (i *IntDBField) GetDefault() string {
	if i.DefaultFuncStruct.PackageFunc != "" {
		i.RequiredPackages = append(i.RequiredPackages, i.DefaultFuncStruct.PackageAddress)
		return i.DefaultFuncStruct.PackageFunc + "()"
	} else {
		return strconv.Itoa(i.DefaultValue)
	}
}

func (i *IntDBField) Default(v int) *IntDBField {
	i.DefaultValue = v
	i.HaveDefault = true
	return i
}

func (i *IntDBField) DefaultFunc(v func() int) *IntDBField {
	i.DefaultFuncStruct.DefaultFunc(v)
	i.RequiredPackages = append(i.RequiredPackages, i.DefaultFuncStruct.PackageAddress)
	i.HaveDefault = true
	return i
}

func IntField(name string) *IntDBField {
	f := &IntDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "int", FieldTypeInt)
	f.IsGreater = true
	return f
}

func (i *IntDBField) SetDBName(v string) *IntDBField {
	i.DBName = v
	return i
}

func (i *IntDBField) AddSerial() *IntDBField {
	i.Serial = true
	return i
}
