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

func TestBulkDelete(t *testing.T) {
	expectedSQLQuery := "DELETE FROM \"account\" WHERE id IN (@idvalue);"
	ctx := context.TODO()
	mockDB := mock.NewMockClient()
	defer mockDB.Close()

	account1 := models.NewAccount(ctx, mockDB)
	account2 := models.NewAccount(ctx, mockDB)

	namedArgs := pgx.NamedArgs{
		"idvalue": []any{account1.GetID(), account2.GetID()},
	}

	mockDB.ExpectExec(expectedSQLQuery).
		WithArgs(namedArgs).WillReturnResult(pgxmock.NewResult("DELETE", 2))

	accountList := models.NewAccountList(ctx, mockDB)
	err := accountList.Delete(account1, account2)

	assert.Equal(t, nil, err)
}
