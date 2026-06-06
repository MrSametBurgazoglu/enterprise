package models

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	FieldTypeBool = iota + 1
	FieldTypeEnum
	FieldTypeInt
	FieldTypeSmallInt
	FieldTypeBigInt
	FieldTypeFloat32
	FieldTypeFloat64
	FieldTypeString
	FieldTypeTime
	FieldTypeUUID
	FieldTypeUint
	FieldTypeByte
	FieldTypeJSON
	FieldTypeCustom
)

type FieldI interface {
	GetRequiredPackages() []string
	GetFieldType() int
	GetType() string
	GetName() string
	GetDBName() string
	GetCustomType() string
	IsNillable() bool
	IsSerial() bool
	IsCanIn() bool
	GetDefault() string
	IsNumeric() bool
	IsComparable() bool
	GetSQLDefault() (string, bool)
}

type Field struct {
	FieldType         int
	Name              string
	DBName            string
	Type              string
	BaseType          string
	Nillable          bool
	HaveDefault       bool
	IsPrepare         bool
	IsGreater         bool
	HaveCustomType    bool
	IsTime            bool
	IsUUID            bool
	IsBool            bool
	CanIn             bool
	Serial            bool
	CustomDBType      string
	RequiredPackages  []string
	defaultFunc       reflect.Value
	DefaultFuncStruct *FuncStruct
}

func (f *Field) GetName() string {
	return f.Name
}

func (f *Field) GetNameLower() string {
	return strings.ToLower(f.Name)
}

func (f *Field) GetNameTitle() string {
	return strings.ToTitle(f.Name)
}

func (f *Field) GetDBName() string {
	return f.DBName
}

func (f *Field) GetCustomType() string {
	return f.CustomDBType
}

func (f *Field) SetDBName(v string) {
	f.DBName = v
}

func (f *Field) GetType() string {
	return f.Type
}

func (f *Field) GetBaseType() string {
	return f.BaseType
}

func (f *Field) GetAddressType() string {
	if f.Type[0] == '*' {
		return f.Type
	}
	return "*" + f.Type
}

func (f *Field) SetNillable() *Field {
	f.Nillable = true
	f.BaseType = f.Type
	f.Type = "*" + f.Type
	return f
}

func (f *Field) IsNillable() bool {
	return f.Nillable
}

func (f *Field) IsSerial() bool {
	return f.Serial
}

func (f *Field) IsDefault() bool {
	return f.HaveDefault
}

func (f *Field) IsCanIn() bool {
	return f.CanIn
}

func (f *Field) NeedPrepare() bool {
	return f.IsPrepare
}

func (f *Field) CanBeGreater() bool {
	return f.IsGreater
}

func (f *Field) CanTime() bool {
	return f.IsTime
}

func (f *Field) CanUUID() bool {
	return f.IsUUID
}

func (f *Field) CanUseIf() bool {
	return f.IsBool && f.IsNillable()
}

func (f *Field) IsCustomType() bool {
	return f.HaveCustomType
}

func (f *Field) IsNumeric() bool {
	return f.FieldType == FieldTypeInt ||
		f.FieldType == FieldTypeSmallInt ||
		f.FieldType == FieldTypeBigInt ||
		f.FieldType == FieldTypeFloat32 ||
		f.FieldType == FieldTypeFloat64 ||
		f.FieldType == FieldTypeUint
}

func (f *Field) IsComparable() bool {
	return f.IsNumeric() ||
		f.FieldType == FieldTypeString ||
		f.FieldType == FieldTypeTime
}

func (f *Field) GetRequiredPackages() []string {
	if f.defaultFunc.IsValid() {
		f.GetDefault()
	}
	return f.RequiredPackages
}

func (f *Field) GetFieldType() int {
	return f.FieldType
}

func (f *Field) DefaultFunc(v any) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Func {
		panic("not a function")
	}
	f.defaultFunc = val
	f.DefaultFuncStruct.DefaultFunc(v)
	f.HaveDefault = true
}

func (f *Field) GetDefault() string {
	if f.defaultFunc.IsValid() {
		pc := f.defaultFunc.Pointer()
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			fullName := fn.Name()
			if !strings.Contains(fullName, ".func") {
				i1 := strings.LastIndex(fullName, ".")
				packageAddress := fullName[:i1]
				i2 := strings.LastIndex(fullName, "/")
				packageFunc := fullName[i2+1:]
				if packageAddress != "" && packageAddress != "main" {
					f.RequiredPackages = append(f.RequiredPackages, packageAddress)
				}
				return packageFunc + "()"
			}
		}

		// Fallback: execute function at generation time
		results := f.defaultFunc.Call(nil)
		val := results[0].Interface()
		switch f.FieldType {
		case FieldTypeInt, FieldTypeSmallInt, FieldTypeBigInt:
			return fmt.Sprintf("%d", val)
		case FieldTypeUint:
			return fmt.Sprintf("%d", val)
		case FieldTypeFloat32, FieldTypeFloat64:
			return fmt.Sprintf("%f", val)
		case FieldTypeString:
			return fmt.Sprintf("%q", val)
		case FieldTypeBool:
			return fmt.Sprintf("%t", val)
		case FieldTypeUUID:
			u := val.(uuid.UUID)
			f.RequiredPackages = append(f.RequiredPackages, "github.com/google/uuid")
			return fmt.Sprintf("uuid.MustParse(%q)", u.String())
		case FieldTypeTime:
			t := val.(time.Time)
			f.RequiredPackages = append(f.RequiredPackages, "time")
			return fmt.Sprintf("time.Unix(%d, %d)", t.Unix(), t.UnixNano()%1e9)
		default:
			return fmt.Sprintf("%v", val)
		}
	}
	return ""
}

func (f *Field) GetSQLDefault() (string, bool) {
	if !f.HaveDefault {
		return "", false
	}
	if f.defaultFunc.IsValid() {
		pc := f.defaultFunc.Pointer()
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			fullName := fn.Name()
			if strings.HasSuffix(fullName, "uuid.New") {
				return "gen_random_uuid()", true
			}
			if strings.HasSuffix(fullName, "time.Now") {
				return "now()", true
			}
		}
	}
	return "", false
}

func (f *Field) setField(name, typeName string, fieldType int) {
	f.FieldType = fieldType
	f.Name = name
	f.Type = typeName
	f.BaseType = typeName
	f.SetDBNameManually(name)
	f.DefaultFuncStruct = new(FuncStruct)
}

func (f *Field) SetDBNameManually(name string) {
	snake := ConvertToSnakeCase(name)
	f.SetDBName(snake)
}
