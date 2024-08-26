package db_models

import (
	"github.com/MrSametBurgazoglu/enterprise/models"
	"github.com/google/uuid"
)

func Group() *models.Table {
	idField := models.UUIDField("ID").DefaultFunc(uuid.New)

	tb := &models.Table{
		Fields: []models.FieldI{
			idField,
			models.StringField("Name"),
			models.StringField("Surname"),
		},
		Relations: []*models.Relation{
			models.ManyToMany(AccountName, "group_id", "account_id", "id", AccountGroupName),
		},
	}

	tb.SetTableName(GroupName)
	tb.SetIDField(idField)

	return tb
}
