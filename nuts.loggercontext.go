package gonuts

import (
	"context"

	"go.uber.org/zap"
)

const (
	RequestIdFieldKey     = "requestId"
	RequestUserIdFieldKey = "requestUserId"
	requestIdPrefix       = "rid"
	requestIdLength       = 16
)

type ctxKey int

const (
	requestIdCtxKey ctxKey = iota
	requestUserIdCtxKey
)

func GenerateRequestId() string {
	return NID(requestIdPrefix, requestIdLength)
}

func NewContextWithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, requestIdCtxKey, requestId)
}

func NewContextWithRequestUserId(ctx context.Context, requestUserId string) context.Context {
	return context.WithValue(ctx, requestUserIdCtxKey, requestUserId)
}

func RequestIdFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	requestId, ok := ctx.Value(requestIdCtxKey).(string)
	if !ok {
		return ""
	}
	return requestId
}

func RequestUserIdFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	requestUserId, ok := ctx.Value(requestUserIdCtxKey).(string)
	if !ok {
		return ""
	}
	return requestUserId
}

func NewLoggerFromContext(ctx context.Context, logger *zap.SugaredLogger) *zap.SugaredLogger {
	if logger == nil {
		logger = L
	}

	requestId := RequestIdFromContext(ctx)
	if requestId != "" {
		logger = logger.With(RequestIdFieldKey, requestId)
	}

	requestUserId := RequestUserIdFromContext(ctx)
	if requestUserId != "" {
		logger = logger.With(RequestUserIdFieldKey, requestUserId)
	}

	return logger
}

func LoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	return NewLoggerFromContext(ctx, L)
}
