package models

type EnumDBField struct {
	*Field
	DefaultValue      string
	DefaultFuncStruct *FuncStruct
	TypeName          string
	Values            []string
}

func (s *EnumDBField) Default(v string) *EnumDBField {
	s.DefaultValue = v
	s.HaveDefault = true
	return s
}

func (s *EnumDBField) DefaultFunc(v func() string) *EnumDBField {
	s.DefaultFuncStruct.DefaultFunc(v)
	s.HaveDefault = true
	return s
}

func (s *EnumDBField) GetDefault() string {
	if s.DefaultFuncStruct.PackageFunc != "" {
		s.RequiredPackages = append(s.RequiredPackages, s.DefaultFuncStruct.PackageAddress)
		return s.DefaultFuncStruct.PackageFunc + "()"
	} else {
		return s.DefaultValue
	}
}

func EnumField(name string, values []string) *EnumDBField {
	f := &EnumDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.TypeName = name
	f.setField(name, f.TypeName, FieldTypeEnum)
	f.Values = values
	f.HaveCustomType = true
	return f
}

func (s *EnumDBField) SetDBName(v string) *EnumDBField {
	s.DBName = v
	return s
}
