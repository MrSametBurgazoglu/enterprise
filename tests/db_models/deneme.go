package db_models

import (
	"github.com/MrSametBurgazoglu/enterprise/models"
	"github.com/google/uuid"
)

func Deneme() *models.Table {
	idField := models.UUIDField("ID").DefaultFunc(uuid.New)
	denemeTypeEnumValues := []string{"Test", "Deneme"}
	testRelationField := models.UUIDField("TestID").SetNillable()

	tb := &models.Table{
		Fields: []models.FieldI{
			idField,
			testRelationField,
			models.IntField("Count"),
			models.BoolField("IsActive").Default(true),
			models.EnumField("DenemeType", denemeTypeEnumValues),
		},
		Relations: []*models.Relation{
			models.ManyToOne("Test", "id", testRelationField.DBName),
			models.OneToMany("Account", idField.DBName, "deneme_id"),
		},
	}

	tb.SetTableName(DenemeName)
	tb.SetIDField(idField)

	return tb
}
