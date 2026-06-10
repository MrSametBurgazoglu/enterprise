package models

import (
	"strconv"
)

type DecimalDBField struct {
	*Field
	DefaultValue      string
	DefaultFuncStruct *FuncStruct
}

// DecimalField defines a numeric(precision, scale) column. The value is
// carried as a string (to avoid precision loss).
func DecimalField(name string, precision, scale int) *DecimalDBField {
	f := &DecimalDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "string", FieldTypeDecimal)
	f.Precision = precision
	f.Scale = scale
	f.CanIn = true
	return f
}

func (d *DecimalDBField) SetDBName(v string) *DecimalDBField { d.DBName = v; return d }

func (d *DecimalDBField) Default(v string) *DecimalDBField {
	d.DefaultValue = v
	d.HaveDefault = true
	return d
}

func (d *DecimalDBField) DefaultFunc(v func() string) *DecimalDBField {
	d.Field.DefaultFunc(v)
	return d
}

func (d *DecimalDBField) GetDefault() string {
	if d.defaultFunc.IsValid() {
		return d.Field.GetDefault()
	}
	return strconv.Quote(d.DefaultValue)
}

func (d *DecimalDBField) GetSQLDefault() (string, bool) {
	if !d.HaveDefault {
		return "", false
	}
	if d.defaultFunc.IsValid() {
		return d.Field.GetSQLDefault()
	}
	return d.DefaultValue, true
}
