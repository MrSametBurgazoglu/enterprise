package db_models

import (
	"github.com/MrSametBurgazoglu/enterprise/models"
	"github.com/MrSametBurgazoglu/enterprise/tests/custom_data_type"
	"github.com/google/uuid"
)

func Test() *models.Table {
	idField := models.UUIDField("ID").DefaultFunc(uuid.New)

	tb := &models.Table{
		Fields: []models.FieldI{
			idField,
			models.StringField("Name"),
			models.TimeField("CreatedAt"),
			models.CustomField("Info", "text", custom_data_type.Custom{}).SetNillable(),
		},
		Relations: []*models.Relation{
			models.OneToMany(DenemeName, idField.DBName, "test_id"),
		},
	}

	tb.SetTableName(TestName)
	tb.SetIDField(idField)

	return tb
}
