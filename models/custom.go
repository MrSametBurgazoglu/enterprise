package models

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
)

// CustomDBField represents a field with a custom type in the database.
// Note: defaults for CustomField must be function-based (via DefaultFunc)
// as there is no single-value Default method.
type CustomDBField struct {
	*Field
}

type CustomDBFieldI interface {
	Scan(src any) error
	Value() (driver.Value, error)
}

func (i *CustomDBField) GetSQLDefault() (string, bool) {
	if !i.HaveDefault {
		return "", false
	}
	return i.Field.GetSQLDefault()
}

func CustomField(name, postgresType string, v CustomDBFieldI) *CustomDBField {
	f := &CustomDBField{}
	f.Field = new(Field)
	t := reflect.TypeOf(v)

	// Get the name of the struct
	fmt.Println("Struct Name:", t.Name())

	// Get the package path of the struct
	fmt.Println("Package Path:")
	pkg := t.PkgPath()
	lastPackage := pkg[strings.LastIndex(pkg, "/")+1:]
	f.setField(name, fmt.Sprintf("%s.%s", lastPackage, t.Name()), FieldTypeCustom)
	f.CustomDBType = postgresType
	f.RequiredPackages = append(f.RequiredPackages, t.PkgPath())
	return f
}

func (i *CustomDBField) SetDBName(v string) *CustomDBField {
	i.DBName = v
	return i
}
