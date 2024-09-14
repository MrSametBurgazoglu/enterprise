package models

type JSONDBField struct {
	*Field
	DefaultFuncStruct *FuncStruct
}

func (u *JSONDBField) DefaultFunc(v func() map[string]any) *JSONDBField {
	u.DefaultFuncStruct.DefaultFunc(v)
	u.HaveDefault = true
	return u
}

func (u *JSONDBField) GetDefault() string {
	if u.DefaultFuncStruct.PackageFunc != "" {
		u.RequiredPackages = append(u.RequiredPackages, u.DefaultFuncStruct.PackageAddress)
		return u.DefaultFuncStruct.PackageFunc + "()"
	}
	return "map[string]any"
}

func (u *JSONDBField) PrepareFunc() string {
	return "make(map[string]any)"
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
