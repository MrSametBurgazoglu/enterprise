package models

import "strconv"

type Float64DBField struct {
	*Field
	DefaultValue      float64
	DefaultFuncStruct *FuncStruct
}

func (i *Float64DBField) GetDefault() string {
	if i.defaultFunc.IsValid() {
		return i.Field.GetDefault()
	} else {
		return strconv.FormatFloat(i.DefaultValue, 'g', -1, 64)
	}
}

func (i *Float64DBField) GetSQLDefault() (string, bool) {
	if !i.HaveDefault {
		return "", false
	}
	if i.defaultFunc.IsValid() {
		return i.Field.GetSQLDefault()
	}
	return strconv.FormatFloat(i.DefaultValue, 'g', -1, 64), true
}

func (i *Float64DBField) Default(v float64) *Float64DBField {
	i.DefaultValue = v
	i.HaveDefault = true
	return i
}

func (i *Float64DBField) DefaultFunc(v func() float64) *Float64DBField {
	i.Field.DefaultFunc(v)
	return i
}

func Float64Field(name string) *Float64DBField {
	f := &Float64DBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "float64", FieldTypeFloat64)
	f.IsGreater = true
	f.CanIn = true
	return f
}

func (i *Float64DBField) SetDBName(v string) *Float64DBField {
	i.DBName = v
	return i
}
