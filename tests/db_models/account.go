package db_models

import (
	"github.com/MrSametBurgazoglu/enterprise/models"
	"github.com/google/uuid"
)

func Account() *models.Table {
	idField := models.UUIDField("ID").DefaultFunc(uuid.New)

	tb := &models.Table{
		Fields: []models.FieldI{
			idField,
			models.StringField("Name"),
			models.StringField("Surname"),
			models.UUIDField("DenemeID").SetNillable(),
			models.UintField("Serial").AddSerial(),
		},
		Relations: []*models.Relation{
			models.ManyToOne(DenemeName, idField.DBName, "deneme_id"),
			models.ManyToMany(GroupName, "account_id", "group_id", "id", AccountGroupName),
		},
	}

	tb.SetTableName(AccountName)
	tb.SetIDField(idField)

	return tb
}
