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

func TestCreate(t *testing.T) {
	expectedSQLQuery := "INSERT INTO \"account\" (id,name,surname) VALUES (@id,@name,@surname) RETURNING serial;"
	ctx := context.TODO()
	mockDB := mock.NewMockClient()
	defer mockDB.Close()

	account := models.NewAccount(ctx, mockDB)
	account.SetDefaults()
	account.SetName("name")
	account.SetSurname("surname")
	var serial uint = 1

	resultRow := pgxmock.NewRows([]string{"serial"}).AddRow(serial)

	namedArgs := pgx.NamedArgs{"id": account.GetID(), "name": account.GetName(), "surname": account.GetSurname()}

	mockDB.ExpectQuery(expectedSQLQuery).
		WithArgs(namedArgs).WillReturnRows(resultRow)

	err := account.Create()

	assert.Equal(t, nil, err)
	assert.Equal(t, serial, account.GetSerial())
}
