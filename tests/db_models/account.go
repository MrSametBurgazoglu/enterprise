package db_models

import (
	"reflect"

	"github.com/MrSametBurgazoglu/enterprise/models"
	"github.com/MrSametBurgazoglu/enterprise/tests/custom_data_type"
	"github.com/google/uuid"
)

func Account() *models.Table {
	idField := models.UUIDField("ID").DefaultFunc(uuid.New)

	tb := &models.Table{
		Fields: []models.FieldI{
			idField,
			models.StringField("Name"),
			models.StringField("Surname"),
			models.StringField("Status").Default("active"),
			models.UUIDField("DenemeID").SetNillable(),
			models.UintField("Serial").AddSerial(),
			models.EnumField("Role", []string{"admin", "user"}).GoType(reflect.TypeOf(custom_data_type.UserRole(""))).Default("user"),
			models.DecimalField("Balance", 10, 2).Default("0.00"),
			models.StringArrayField("Tags").SetNillable(),
		},
		Relations: []*models.Relation{
			models.ManyToOne(DenemeName, idField.DBName, "deneme_id"),
			models.ManyToMany(GroupName, "account_id", "group_id", "id", AccountGroupName),
		},
	}

	tb.SetTableName(AccountName)
	tb.SetIDField(idField)
	tb.AddIndex("name_surname_index", "Name", "Surname")
	tb.AddUniqueIndex("account_name_unique", "Name")

	return tb
}
