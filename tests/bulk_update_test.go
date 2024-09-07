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

func TestBulkUpdate(t *testing.T) {
	expectedSQLQuery := "UPDATE \"account\" SET \"name\" = @name, \"surname\" = @surname WHERE id IN (@idvalue)"
	ctx := context.TODO()
	mockDB := mock.NewMockClient()
	defer mockDB.Close()

	account1 := models.NewAccount(ctx, mockDB)
	account1.SetName("name")
	account1.SetSurname("surname")

	account2 := models.NewAccount(ctx, mockDB)
	account2.SetName("name")
	account2.SetSurname("surname")

	namedArgs := pgx.NamedArgs{
		"name":    "name",
		"surname": "surname",
		"idvalue": []any{account1.GetID(), account2.GetID()},
	}

	mockDB.ExpectExec(expectedSQLQuery).
		WithArgs(namedArgs).WillReturnResult(pgxmock.NewResult("UPDATE", 2))

	accountList := models.NewAccountList(ctx, mockDB)
	err := accountList.Update(account1, account2)

	assert.Equal(t, nil, err)
}
