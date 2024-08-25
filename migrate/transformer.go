package migrate

import (
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"fmt"
	"github.com/MrSametBurgazoglu/enterprise/models"
)

func TransformSchemaToAtlasSchema(schemaName string, tables []*models.Table) *schema.Schema {
	dbSchema := &schema.Schema{Name: schemaName}
	for _, table := range tables {
		dbSchema.Tables = append(dbSchema.Tables, TransformTableToAtlasTable(table))
	}
	manyToManyTables := make(map[string]bool)
	for _, table := range tables {
		for _, relation := range table.Relations {
			if relation.ManyTableDBName != "" {
				_, ok := manyToManyTables[relation.ManyTableDBName]
				if ok {
					continue
				}
				manyToManyTables[relation.ManyTableDBName] = true
				relationTable := &schema.Table{
					Name: relation.ManyTableDBName,
					PrimaryKey: &schema.Index{
						Parts: []*schema.IndexPart{{C: &schema.Column{Name: "id"}}},
					},
				}
				relationTable.Columns = []*schema.Column{
					schema.NewIntColumn("id", postgres.TypeSmallSerial),
					{Name: fmt.Sprintf("%s_id", table.DBName), Type: &schema.ColumnType{Type: &postgres.UUIDType{T: postgres.TypeUUID}}},
					{Name: fmt.Sprintf("%s_id", relation.RelationTableDBName), Type: &schema.ColumnType{Type: &postgres.UUIDType{T: postgres.TypeUUID}}},
				}
				dbSchema.Tables = append(dbSchema.Tables, relationTable)
			}
		}
	}

	return dbSchema
}

func TransformTableToAtlasTable(table *models.Table) *schema.Table {
	dbTable := &schema.Table{
		Name: table.DBName,
		PrimaryKey: &schema.Index{
			Parts: []*schema.IndexPart{{C: &schema.Column{Name: table.IDDBField}}},
		},
	}
	for _, field := range table.Fields {
		dbTable.Columns = append(dbTable.Columns, TransformFieldToAtlasColumn(field))
	}
	return dbTable
}

func TransformFieldToAtlasColumn(field models.FieldI) *schema.Column {
	column := &schema.Column{Name: field.GetDBName(), Type: &schema.ColumnType{}}
	column.SetNull(field.IsNillable())
	var t schema.Type
	switch field.GetFieldType() {
	case models.FieldTypeBool:
		t = &schema.BoolType{T: postgres.TypeBoolean}
	case models.FieldTypeSmallInt: //int16
		t = &schema.IntegerType{T: postgres.TypeSmallInt}
		if field.IsSerial() {
			t = &schema.IntegerType{T: postgres.TypeSmallSerial}
		}
	case models.FieldTypeInt: //int32
		t = &schema.IntegerType{T: postgres.TypeInt}
		if field.IsSerial() {
			t = &schema.IntegerType{T: postgres.TypeSerial}
		}
	case models.FieldTypeBigInt: //int64
		t = &schema.IntegerType{T: postgres.TypeBigInt}
		if field.IsSerial() {
			t = &schema.IntegerType{T: postgres.TypeSerial}
		}
	case models.FieldTypeUint:
		t = &schema.IntegerType{T: postgres.TypeInt, Unsigned: true}
		if field.IsSerial() {
			t = &schema.IntegerType{T: postgres.TypeSmallSerial}
		}
	case models.FieldTypeFloat32:
		t = &schema.FloatType{T: postgres.TypeReal}
	case models.FieldTypeFloat64:
		t = &schema.FloatType{T: postgres.TypeDouble}
	case models.FieldTypeUUID:
		t = &postgres.UUIDType{T: postgres.TypeUUID}
	case models.FieldTypeString:
		t = &schema.StringType{T: postgres.TypeVarChar}
	case models.FieldTypeTime:
		t = &schema.TimeType{T: postgres.TypeTimestampWTZ}
	case models.FieldTypeEnum:
		t = &schema.StringType{T: postgres.TypeVarChar}
	default:
		return nil
	}

	column.Type.Type = t
	return column
}
