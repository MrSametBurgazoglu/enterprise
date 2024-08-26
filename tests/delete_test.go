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

func TestDelete(t *testing.T) {
	expectedSQLQuery := "DELETE FROM \"account\" WHERE id = @idvalue;"
	ctx := context.TODO()
	mockDB := mock.NewMockClient()
	defer mockDB.Close()

	account := models.NewAccount(ctx, mockDB)

	namedArgs := pgx.NamedArgs{"idvalue": account.GetID()}
	mockDB.ExpectExec(expectedSQLQuery).
		WithArgs(namedArgs).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := account.Delete()

	assert.Equal(t, nil, err)
}
