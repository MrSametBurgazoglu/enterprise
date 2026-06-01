package models

import "strconv"

type UintDBField struct {
	*Field
	DefaultValue      uint
	DefaultFuncStruct *FuncStruct
}

func (i *UintDBField) GetDefault() string {
	if i.defaultFunc.IsValid() {
		return i.Field.GetDefault()
	} else {
		return strconv.Itoa(int(i.DefaultValue))
	}
}

func (i *UintDBField) Default(v uint) *UintDBField {
	i.DefaultValue = v
	i.HaveDefault = true
	return i
}

func (i *UintDBField) DefaultFunc(v func() uint) *UintDBField {
	i.Field.DefaultFunc(v)
	return i
}

func UintField(name string) *UintDBField {
	f := &UintDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "uint", FieldTypeUint)
	f.IsGreater = true
	f.CanIn = true
	return f
}

func (i *UintDBField) SetDBName(v string) *UintDBField {
	i.DBName = v
	return i
}

func (i *UintDBField) AddSerial() *UintDBField {
	i.Serial = true
	return i
}
