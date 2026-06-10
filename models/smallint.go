package models

import "strconv"

type SmallIntDBField struct {
	*Field
	DefaultValue      int16
	DefaultFuncStruct *FuncStruct
}

func (i *SmallIntDBField) GetDefault() string {
	if i.defaultFunc.IsValid() {
		return i.Field.GetDefault()
	} else {
		return strconv.FormatInt(int64(i.DefaultValue), 10)
	}
}

func (i *SmallIntDBField) GetSQLDefault() (string, bool) {
	if !i.HaveDefault {
		return "", false
	}
	if i.defaultFunc.IsValid() {
		return i.Field.GetSQLDefault()
	}
	return strconv.FormatInt(int64(i.DefaultValue), 10), true
}

func (i *SmallIntDBField) Default(v int16) *SmallIntDBField {
	i.DefaultValue = v
	i.HaveDefault = true
	return i
}

func (i *SmallIntDBField) DefaultFunc(v func() int16) *SmallIntDBField {
	i.Field.DefaultFunc(v)
	return i
}

func SmallIntField(name string) *SmallIntDBField {
	f := &SmallIntDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "int16", FieldTypeSmallInt)
	f.IsGreater = true
	f.CanIn = true
	return f
}

func (i *SmallIntDBField) SetDBName(v string) *SmallIntDBField {
	i.DBName = v
	return i
}

func (i *SmallIntDBField) AddSerial() *SmallIntDBField {
	i.Serial = true
	return i
}
