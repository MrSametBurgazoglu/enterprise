package tests

import (
	"context"
	"github.com/MrSametBurgazoglu/enterprise/logger"
	"github.com/MrSametBurgazoglu/enterprise/mock"
	"github.com/jackc/pgx/v5"
	"testing"
)

func TestLog(t *testing.T) {
	ctx := context.TODO()
	mockHandler := mock.NewHandler()
	mockHandler.On("DebugContext", ctx, "%s executed with args:%v", []any{"TEST SQL", pgx.NamedArgs{"hello": "world"}})

	l := logger.NewLogger(mockHandler)
	l.Log(ctx, "TEST SQL", pgx.NamedArgs{"hello": "world"})
}
