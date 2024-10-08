package tests

import (
	"context"
	"github.com/MrSametBurgazoglu/enterprise/client"
	"github.com/MrSametBurgazoglu/enterprise/mock"
	"github.com/MrSametBurgazoglu/enterprise/tests/custom_data_type"
	"github.com/MrSametBurgazoglu/enterprise/tests/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	expectedSQLQuery := "SELECT \"account\".\"id\", \"account\".\"name\", \"account\".\"surname\", \"account\".\"deneme_id\", \"account\".\"serial\" FROM account WHERE ((\"account\".\"id\" = @account__id));"
	id := uuid.New()
	ctx := context.TODO()
	var serial uint = 5

	mockDB := mock.NewMockClient()
	defer mockDB.Close()

	account := models.NewAccount(ctx, mockDB)
	account.Where(account.IsIDEqual(id))

	resultRow := pgxmock.NewRows([]string{"id", "name", "surname", "deneme_id", "serial"}).AddRow(id, "name", "surname", nil, serial)

	namedArgs := pgx.NamedArgs{"account__id": id}
	mockDB.ExpectQuery(expectedSQLQuery).
		WithArgs(namedArgs).WillReturnRows(resultRow)

	err := account.Get()

	assert.Equal(t, nil, err)
	assert.Equal(t, "name", account.GetName())
	assert.Equal(t, "surname", account.GetSurname())
	assert.Equal(t, &uuid.Nil, account.GetDenemeID())
	assert.Equal(t, serial, account.GetSerial())
}

func TestGetWithRelations(t *testing.T) {
	expectedSQLQuery := `
	SELECT
	"test"."id",
    "test"."name",
    "test"."created_at",
    "test"."info",
    "deneme"."id",
    "deneme"."test_id",
    "deneme"."count",
    "deneme"."is_active",
    "deneme"."deneme_type",
    "account"."id",
    "account"."name",
    "account"."surname",
    "account"."deneme_id",
    "account"."serial",
    "group"."id",
    "group"."name",
    "group"."surname",
    "group"."data"
	FROM test
	    LEFT JOIN "deneme" ON "test"."id" = "deneme"."test_id"
		RIGHT JOIN "account" ON "deneme"."id" = "account"."deneme_id"
		LEFT JOIN "group" ON "account"."group_id" = "group"."account_id"
	WHERE (("deneme"."count" = @deneme__count)
		AND ("account"."name" = @account__name) 
		AND ("group"."name" = @group__name));`

	testID := uuid.New()
	denemeID := uuid.New()
	accountID1 := uuid.New()
	accountID2 := uuid.New()
	groupID1 := uuid.New()
	groupID2 := uuid.New()
	groupID3 := uuid.New()
	createdAt := time.Now()

	ctx := context.TODO()
	var serial uint = 5

	mockDB := mock.NewMockClient()
	defer mockDB.Close()

	custom := custom_data_type.Custom{Hello: "hello", World: "world"}
	v, _ := custom.Value()

	values := [][]any{
		{testID, "test_name", createdAt, nil,
			denemeID, testID, 20, true, models.DenemeTypeDenemeType,
			accountID1, "account_name", "account1_surname", denemeID, serial,
			groupID1, "group_name", "group1_surname", map[string]any{"deneme": "value1"}},
		{testID, "test_name", createdAt, nil,
			denemeID, testID, 20, true, models.DenemeTypeDenemeType,
			accountID2, "account_name", "account2_surname", denemeID, serial,
			groupID2, "group_name", "group2_surname", map[string]any{"deneme": "value2"}},
		{testID, "test_name", createdAt, any(v),
			denemeID, testID, 20, true, models.DenemeTypeDenemeType,
			accountID2, "account_name", "account2_surname", denemeID, serial,
			groupID3, "group_name", "group3_surname", map[string]any{"deneme": "value3"}},
	}

	resultRow := pgxmock.NewRows([]string{
		"id",
		"name",
		"created_at",
		"info",
		"id",
		"test_id",
		"count",
		"is_active",
		"deneme_type",
		"id",
		"name",
		"surname",
		"deneme_id",
		"serial",
		"id",
		"name",
		"surname",
		"data"},
	).
		AddRows(values...)

	namedArgs := pgx.NamedArgs{"deneme__count": 20, "account__name": "account_name", "group__name": "group_name"}
	mockDB.ExpectQuery(expectedSQLQuery).
		WithArgs(namedArgs).WillReturnRows(resultRow)

	test := models.NewTest(ctx, mockDB)
	test.WithDenemeList(func(denemeList *models.DenemeList) {
		denemeList.Where(denemeList.IsCountEqual(20))
		denemeList.WithAccountList(func(accountList *models.AccountList) {
			accountList.Where(accountList.IsNameEqual("account_name"))
			accountList.WithGroupList(func(groupList *models.GroupList) {
				groupList.Where(groupList.IsNameEqual("group_name"))
			})
		}).SetJoinType(client.RelationJoinTypeRight)
	})

	err := test.Get()

	assert.Equal(t, nil, err)
	assert.Equal(t, testID, test.GetID())
	assert.Equal(t, "test_name", test.GetName())
	assert.Equal(t, createdAt, test.GetCreatedAt())

	assert.Equal(t, denemeID, test.DenemeList.Items[0].GetID())
	assert.Equal(t, &testID, test.DenemeList.Items[0].GetTestID())
	assert.Equal(t, 20, test.DenemeList.Items[0].GetCount())
	assert.Equal(t, true, test.DenemeList.Items[0].GetIsActive())
	assert.Equal(t, models.DenemeTypeDenemeType, test.DenemeList.Items[0].GetDenemeType())

	assert.Equal(t, accountID1, test.DenemeList.Items[0].AccountList.Items[0].GetID())
	assert.Equal(t, "account_name", test.DenemeList.Items[0].AccountList.Items[0].GetName())
	assert.Equal(t, "account1_surname", test.DenemeList.Items[0].AccountList.Items[0].GetSurname())
	assert.Equal(t, &denemeID, test.DenemeList.Items[0].AccountList.Items[0].GetDenemeID())
	assert.Equal(t, serial, test.DenemeList.Items[0].AccountList.Items[0].GetSerial())

	assert.Equal(t, accountID2, test.DenemeList.Items[0].AccountList.Items[1].GetID())
	assert.Equal(t, "account_name", test.DenemeList.Items[0].AccountList.Items[1].GetName())
	assert.Equal(t, "account2_surname", test.DenemeList.Items[0].AccountList.Items[1].GetSurname())
	assert.Equal(t, &denemeID, test.DenemeList.Items[0].AccountList.Items[1].GetDenemeID())
	assert.Equal(t, serial, test.DenemeList.Items[0].AccountList.Items[1].GetSerial())

	assert.Equal(t, groupID1, test.DenemeList.Items[0].AccountList.Items[0].GroupList.Items[0].GetID())
	assert.Equal(t, "group_name", test.DenemeList.Items[0].AccountList.Items[0].GroupList.Items[0].GetName())
	assert.Equal(t, "group1_surname", test.DenemeList.Items[0].AccountList.Items[0].GroupList.Items[0].GetSurname())
	assert.Equal(t, "value1", test.DenemeList.Items[0].AccountList.Items[0].GroupList.Items[0].GetData()["deneme"])

	assert.Equal(t, groupID2, test.DenemeList.Items[0].AccountList.Items[1].GroupList.Items[0].GetID())
	assert.Equal(t, "group_name", test.DenemeList.Items[0].AccountList.Items[1].GroupList.Items[0].GetName())
	assert.Equal(t, "group2_surname", test.DenemeList.Items[0].AccountList.Items[1].GroupList.Items[0].GetSurname())
	assert.Equal(t, "value2", test.DenemeList.Items[0].AccountList.Items[1].GroupList.Items[0].GetData()["deneme"])

	assert.Equal(t, groupID3, test.DenemeList.Items[0].AccountList.Items[1].GroupList.Items[1].GetID())
	assert.Equal(t, "group_name", test.DenemeList.Items[0].AccountList.Items[1].GroupList.Items[1].GetName())
	assert.Equal(t, "group3_surname", test.DenemeList.Items[0].AccountList.Items[1].GroupList.Items[1].GetSurname())
	assert.Equal(t, "value3", test.DenemeList.Items[0].AccountList.Items[1].GroupList.Items[1].GetData()["deneme"])
}
