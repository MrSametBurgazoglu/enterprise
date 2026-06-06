package models

type JSONDBField struct {
	*Field
	DefaultFuncStruct *FuncStruct
}

func (u *JSONDBField) DefaultFunc(v func() map[string]any) *JSONDBField {
	u.Field.DefaultFunc(v)
	return u
}

func (u *JSONDBField) GetDefault() string {
	if u.defaultFunc.IsValid() {
		return u.Field.GetDefault()
	}
	return "map[string]any"
}

func (u *JSONDBField) GetSQLDefault() (string, bool) {
	if !u.HaveDefault {
		return "", false
	}
	return u.Field.GetSQLDefault()
}

func (u *JSONDBField) PrepareFunc() string {
	return "&map[string]any{}"
}

func JSONField(name string) *JSONDBField {
	f := &JSONDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "map[string]any", FieldTypeJSON)
	return f
}

func (u *JSONDBField) SetDBName(v string) *JSONDBField {
	u.DBName = v
	return u
}

func (u *JSONDBField) SetNillable() *JSONDBField {
	u.Field.SetNillable()
	u.IsPrepare = true
	return u
}
