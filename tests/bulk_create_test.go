package tests

import (
	"context"
	"github.com/MrSametBurgazoglu/enterprise/mock"
	"github.com/MrSametBurgazoglu/enterprise/tests/models"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBulkCreate(t *testing.T) {
	expectedSQLQuery := "INSERT INTO \"account\" (id,name,surname) VALUES (@0id,@0name,@0surname), (@1id,@1name,@1surname);"
	ctx := context.TODO()
	mockDB := mock.NewMockClient()
	defer mockDB.Close()

	account1 := models.NewAccount(ctx, mockDB)
	account1.SetDefaults()
	account1.SetName("name1")
	account1.SetSurname("surname1")

	account2 := models.NewAccount(ctx, mockDB)
	account2.SetDefaults()
	account2.SetName("name2")
	account2.SetSurname("surname2")

	namedArgs := pgx.NamedArgs{
		"0id":      account1.GetID(),
		"0name":    "name1",
		"0surname": "surname1",
		"1id":      account2.GetID(),
		"1name":    "name2",
		"1surname": "surname2",
	}

	mockDB.ExpectExec(expectedSQLQuery).
		WithArgs(namedArgs).WillReturnResult(pgxmock.NewResult("INSERT", 2))

	accountList := models.NewAccountList(ctx, mockDB)
	err := accountList.Create(account1, account2)

	assert.Equal(t, nil, err)
}
