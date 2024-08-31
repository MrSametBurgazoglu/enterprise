package tests

import (
	"context"
	"github.com/MrSametBurgazoglu/enterprise/client"
	"github.com/MrSametBurgazoglu/enterprise/mock"
	"github.com/MrSametBurgazoglu/enterprise/tests/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHook(t *testing.T) {
	expectedSQLQuery := "SELECT \"account\".\"id\", \"account\".\"name\", \"account\".\"surname\", \"account\".\"deneme_id\", \"account\".\"serial\" FROM account WHERE ((\"account\".\"id\" = @account__id));"
	id := uuid.New()
	ctx := context.TODO()
	var serial uint = 5

	mockDB := mock.NewMockClient()
	defer mockDB.Close()

	operationStarted, operationEnded := false, false
	var operationType client.OperationType
	var operationContext context.Context

	models.SetDatabaseAccountOperationHook(func(operationInfo *client.OperationInfo, model *models.Account, operationFunc func() error) error {
		operationStarted = true
		operationType = operationInfo.OperationType
		operationContext = model.GetContext()
		err := operationFunc()
		operationEnded = true
		return err
	})
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

	assert.Equal(t, true, operationStarted)
	assert.Equal(t, true, operationEnded)
	assert.Equal(t, client.OperationTypeGet, operationType)
	assert.Equal(t, ctx, operationContext)
}
