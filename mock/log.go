package mock

import (
	"context"
	"github.com/MrSametBurgazoglu/enterprise/logger"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

var _ logger.Logger = (*Logger)(nil)

type Logger struct {
	mock.Mock
}

func (l *Logger) Log(ctx context.Context, sqlString string, args pgx.NamedArgs) {
	_ = l.Called(ctx, sqlString, args)
}

func (l *Logger) LogError(ctx context.Context, sqlString string, args pgx.NamedArgs, err error) {
	_ = l.Called(ctx, sqlString, args, err)
}

func NewLogger() *Logger {
	return &Logger{}
}
