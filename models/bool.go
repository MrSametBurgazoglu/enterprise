package models

type BoolDBField struct {
	*Field
	DefaultValue      bool
	DefaultFuncStruct *FuncStruct
}

func (i *BoolDBField) GetDefault() string {
	if i.DefaultFuncStruct.PackageFunc != "" {
		i.RequiredPackages = append(i.RequiredPackages, i.DefaultFuncStruct.PackageAddress)
		return i.DefaultFuncStruct.PackageFunc + "()"
	} else if i.DefaultValue {
		return "true"
	} else {
		return "false"
	}
}

func (i *BoolDBField) Default(v bool) *BoolDBField {
	i.DefaultValue = v
	i.HaveDefault = true
	return i
}

func (i *BoolDBField) DefaultFunc(v func() bool) *BoolDBField {
	i.DefaultFuncStruct.DefaultFunc(v)
	i.RequiredPackages = append(i.RequiredPackages, i.DefaultFuncStruct.PackageAddress)
	i.HaveDefault = true
	return i
}

func BoolField(name string) *BoolDBField {
	f := &BoolDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "bool", FieldTypeBool)
	return f
}

func (i *BoolDBField) SetDBName(v string) *BoolDBField {
	i.DBName = v
	return i
}
