package migrate

import (
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"fmt"
	"github.com/MrSametBurgazoglu/enterprise/models"
)

func TransformSchemaToAtlasSchema(schemaName string, tables []*models.Table) *schema.Schema {
	dbSchema := &schema.Schema{Name: schemaName}
	tableMap := make(map[string]*schema.Table)
	for _, table := range tables {
		t := TransformTableToAtlasTable(table)
		dbSchema.Tables = append(dbSchema.Tables, t)
		tableMap[t.Name] = t
	}

	for i, table := range tables {
		for _, relation := range table.Relations {
			symbol := fmt.Sprintf("%s_%s", table.DBName, relation.OnField)

			if relation.RelationType == 1 { //many to one
				continue
			}
			fk := schema.NewForeignKey(symbol).
				SetTable(tableMap[table.DBName]).
				SetRefTable(tableMap[relation.RelationTableDBName]).
				SetOnUpdate(schema.NoAction).
				SetOnDelete(schema.NoAction).
				AddColumns(&schema.Column{Name: relation.OnField}).
				AddRefColumns(&schema.Column{Name: relation.RelationField})
			dbSchema.Tables[i].AddForeignKeys(fk)
		}
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

				primaryKey := schema.NewIntColumn("id", postgres.TypeSmallSerial)

				relationTable := &schema.Table{
					Name: relation.ManyTableDBName,
					PrimaryKey: &schema.Index{
						Parts: []*schema.IndexPart{{C: primaryKey}},
					},
				}

				pk := schema.NewPrimaryKey(primaryKey)
				relationTable.SetPrimaryKey(pk)

				c1 := &schema.Column{Name: fmt.Sprintf("%s_id", table.DBName), Type: &schema.ColumnType{Type: &postgres.UUIDType{T: postgres.TypeUUID}}}
				c2 := &schema.Column{Name: fmt.Sprintf("%s_id", relation.RelationTableDBName), Type: &schema.ColumnType{Type: &postgres.UUIDType{T: postgres.TypeUUID}}}

				relationTable.Columns = []*schema.Column{
					primaryKey,
					c1,
					c2,
				}

				symbol1 := fmt.Sprintf("%s_%s", relation.ManyTableDBName, table.DBName)
				symbol2 := fmt.Sprintf("%s_%s", relation.ManyTableDBName, relation.RelationTableDBName)

				fk1 := schema.NewForeignKey(symbol1).
					SetTable(relationTable).
					SetRefTable(tableMap[table.DBName]).
					SetOnUpdate(schema.NoAction).
					SetOnDelete(schema.NoAction).
					AddColumns(c1).
					AddRefColumns(&schema.Column{Name: relation.RelationField})
				relationTable.AddForeignKeys(fk1)

				fk2 := schema.NewForeignKey(symbol2).
					SetTable(relationTable).
					SetRefTable(tableMap[relation.RelationTableDBName]).
					SetOnUpdate(schema.NoAction).
					SetOnDelete(schema.NoAction).
					AddColumns(c2).
					AddRefColumns(&schema.Column{Name: relation.OnField})
				relationTable.AddForeignKeys(fk2)

				dbSchema.Tables = append(dbSchema.Tables, relationTable)
			}
		}
	}

	return dbSchema
}

func TransformTableToAtlasTable(table *models.Table) *schema.Table {
	primaryKey := TransformFieldToAtlasColumn(table.IDColumn)

	dbTable := &schema.Table{Name: table.DBName}

	pk := schema.NewPrimaryKey(primaryKey)
	dbTable.SetPrimaryKey(pk)

	for _, field := range table.Fields {
		dbTable.Columns = append(dbTable.Columns, TransformFieldToAtlasColumn(field))
	}

	for _, index := range table.Indexes {
		var columns []*schema.Column
		for _, columnName := range index.Columns {
			columns = append(columns, schema.NewColumn(columnName))
		}
		schemaIndex := schema.NewIndex(index.Name)
		schemaIndex.AddColumns(columns...)
		dbTable.Indexes = append(dbTable.Indexes, schemaIndex)
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
	case models.FieldTypeByte:
		t = &schema.BinaryType{T: postgres.TypeBytea}
	case models.FieldTypeJSON:
		t = &schema.JSONType{T: postgres.TypeJSONB}
	case models.FieldTypeCustom:
		t = &schema.UnsupportedType{T: field.GetCustomType()}
	default:
		return nil
	}

	column.Type.Type = t
	return column
}
