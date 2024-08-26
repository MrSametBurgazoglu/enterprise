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

func TestUpdate(t *testing.T) {
	expectedSQLQuery := "UPDATE \"account\" SET \"id\" = @id,  \"name\" = @name, \"surname\" = @surname WHERE id = @idvalue"
	ctx := context.TODO()
	mockDB := mock.NewMockClient()
	defer mockDB.Close()

	account := models.NewAccount(ctx, mockDB)
	account.SetName("name")
	account.SetSurname("surname")

	namedArgs := pgx.NamedArgs{"id": account.GetID(), "name": account.GetName(), "surname": account.GetSurname(), "idvalue": account.GetID()}
	mockDB.ExpectExec(expectedSQLQuery).
		WithArgs(namedArgs).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := account.Update()

	assert.Equal(t, nil, err)
}
