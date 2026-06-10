package models

import "strconv"

type BigIntDBField struct {
	*Field
	DefaultValue      int64
	DefaultFuncStruct *FuncStruct
}

func (i *BigIntDBField) GetDefault() string {
	if i.defaultFunc.IsValid() {
		return i.Field.GetDefault()
	} else {
		return strconv.FormatInt(i.DefaultValue, 10)
	}
}

func (i *BigIntDBField) GetSQLDefault() (string, bool) {
	if !i.HaveDefault {
		return "", false
	}
	if i.defaultFunc.IsValid() {
		return i.Field.GetSQLDefault()
	}
	return strconv.FormatInt(i.DefaultValue, 10), true
}

func (i *BigIntDBField) Default(v int64) *BigIntDBField {
	i.DefaultValue = v
	i.HaveDefault = true
	return i
}

func (i *BigIntDBField) DefaultFunc(v func() int64) *BigIntDBField {
	i.Field.DefaultFunc(v)
	return i
}

func BigIntField(name string) *BigIntDBField {
	f := &BigIntDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "int64", FieldTypeBigInt)
	f.IsGreater = true
	f.CanIn = true
	return f
}

func (i *BigIntDBField) SetDBName(v string) *BigIntDBField {
	i.DBName = v
	return i
}

func (i *BigIntDBField) AddSerial() *BigIntDBField {
	i.Serial = true
	return i
}
