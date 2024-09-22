package mock

import (
	"context"
	"github.com/MrSametBurgazoglu/enterprise/logger"
	"github.com/stretchr/testify/mock"
)

var _ logger.Handler = (*Handler)(nil)

type Handler struct {
	mock.Mock
}

func (h *Handler) DebugContext(ctx context.Context, msg string, args ...any) {
	h.Called(ctx, msg, args)
}

func (h *Handler) ErrorContext(ctx context.Context, msg string, args ...any) {
	h.Called(ctx, msg, args)
}

func NewHandler() *Handler {
	return &Handler{}
}
