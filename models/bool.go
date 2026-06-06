package models

type BoolDBField struct {
	*Field
	DefaultValue      bool
	DefaultFuncStruct *FuncStruct
}

func (i *BoolDBField) GetDefault() string {
	if i.defaultFunc.IsValid() {
		return i.Field.GetDefault()
	} else if i.DefaultValue {
		return "true"
	} else {
		return "false"
	}
}

func (i *BoolDBField) GetSQLDefault() (string, bool) {
	if !i.HaveDefault {
		return "", false
	}
	if i.defaultFunc.IsValid() {
		return i.Field.GetSQLDefault()
	}
	if i.DefaultValue {
		return "true", true
	}
	return "false", true
}

func (i *BoolDBField) Default(v bool) *BoolDBField {
	i.DefaultValue = v
	i.HaveDefault = true
	return i
}

func (i *BoolDBField) DefaultFunc(v func() bool) *BoolDBField {
	i.Field.DefaultFunc(v)
	return i
}

func BoolField(name string) *BoolDBField {
	f := &BoolDBField{}
	f.Field = new(Field)
	f.IsBool = true
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "bool", FieldTypeBool)
	return f
}

func (i *BoolDBField) SetDBName(v string) *BoolDBField {
	i.DBName = v
	return i
}
