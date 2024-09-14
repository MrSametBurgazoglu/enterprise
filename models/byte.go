package models

type ByteDBField struct {
	*Field
	DefaultValue      []byte
	DefaultFuncStruct *FuncStruct
}

func (i *ByteDBField) GetDefault() string {
	if i.DefaultFuncStruct.PackageFunc != "" {
		i.RequiredPackages = append(i.RequiredPackages, i.DefaultFuncStruct.PackageAddress)
		return i.DefaultFuncStruct.PackageFunc + "()"
	}
	return ""
}

func (i *ByteDBField) Default(v []byte) *ByteDBField {
	i.DefaultValue = v
	i.HaveDefault = true
	return i
}

func (i *ByteDBField) DefaultFunc(v func() bool) *ByteDBField {
	i.DefaultFuncStruct.DefaultFunc(v)
	i.RequiredPackages = append(i.RequiredPackages, i.DefaultFuncStruct.PackageAddress)
	i.HaveDefault = true
	return i
}

func ByteField(name string) *ByteDBField {
	f := &ByteDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "[]byte", FieldTypeByte)
	return f
}

func (i *ByteDBField) SetDBName(v string) *ByteDBField {
	i.DBName = v
	return i
}
