package models

type StringDBField struct {
	*Field
	DefaultValue      string
	DefaultFuncStruct *FuncStruct
}

func (s *StringDBField) Default(v string) *StringDBField {
	s.DefaultValue = v
	s.HaveDefault = true
	return s
}

func (s *StringDBField) DefaultFunc(v func() string) *StringDBField {
	s.DefaultFuncStruct.DefaultFunc(v)
	s.HaveDefault = true
	return s
}

func (s *StringDBField) GetDefault() string {
	if s.DefaultFuncStruct.PackageFunc != "" {
		s.RequiredPackages = append(s.RequiredPackages, s.DefaultFuncStruct.PackageAddress)
		return s.DefaultFuncStruct.PackageFunc + "()"
	} else {
		return s.DefaultValue
	}
}

func StringField(name string) *StringDBField {
	f := &StringDBField{}
	f.Field = new(Field)
	f.DefaultFuncStruct = new(FuncStruct)
	f.setField(name, "string", FieldTypeString)
	f.CanIn = true
	return f
}

func (s *StringDBField) SetDBName(v string) *StringDBField {
	s.DBName = v
	return s
}
