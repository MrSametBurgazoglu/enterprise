package logger

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Logger interface {
	Log(ctx context.Context, sqlString string, args pgx.NamedArgs)
	LogError(ctx context.Context, sqlString string, args pgx.NamedArgs, err error)
}

type Handler interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type logger struct {
	handler Handler
}

const DebugSqlFormat = "%s executed with args:%v"
const ErrSqlFormat = "%s executed with args:%v; err:%s"

func (l logger) Log(ctx context.Context, sqlString string, args pgx.NamedArgs) {
	l.handler.DebugContext(ctx, DebugSqlFormat, sqlString, args)
}

func (l logger) LogError(ctx context.Context, sqlString string, args pgx.NamedArgs, err error) {
	l.handler.ErrorContext(ctx, ErrSqlFormat, sqlString, args, err.Error())
}

func NewLogger(handler Handler) Logger {
	l := logger{}
	l.handler = handler
	return &l
}
