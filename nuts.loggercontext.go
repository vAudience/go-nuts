package gonuts

import (
	"context"

	"go.uber.org/zap"
)

const (
	RequestIdFieldKey     = "requestId"
	RequestOrgIdFieldKey  = "requestOrgId"
	RequestUserIdFieldKey = "requestUserId"
	requestIdPrefix       = "rid"
	requestIdLength       = 16
)

type ctxKey int

const (
	requestIdCtxKey ctxKey = iota
	requestOrgIdCtxKey
	requestUserIdCtxKey
)

func GenerateRequestId() string {
	return NID(requestIdPrefix, requestIdLength)
}

func NewContextWithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, requestIdCtxKey, requestId)
}

func NewContextWithRequestOrgId(ctx context.Context, requestOrgId string) context.Context {
	return context.WithValue(ctx, requestOrgIdCtxKey, requestOrgId)
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

func RequestOrgIdFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	requestOrgId, ok := ctx.Value(requestOrgIdCtxKey).(string)
	if !ok {
		return ""
	}
	return requestOrgId
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

func LoggerFieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	requestId := RequestIdFromContext(ctx)
	if requestId != "" {
		fields = append(fields, zap.String(RequestIdFieldKey, requestId))
	}

	requestOrgId := RequestOrgIdFromContext(ctx)
	if requestOrgId != "" {
		fields = append(fields, zap.String(RequestOrgIdFieldKey, requestOrgId))
	}

	requestUserId := RequestUserIdFromContext(ctx)
	if requestUserId != "" {
		fields = append(fields, zap.String(RequestUserIdFieldKey, requestUserId))
	}
	return fields
}

func NewLoggerFromContext(ctx context.Context, logger *zap.SugaredLogger) *zap.SugaredLogger {
	if logger == nil {
		logger = L
	}

	requestId := RequestIdFromContext(ctx)
	if requestId != "" {
		logger = logger.With(RequestIdFieldKey, requestId)
	}

	requestOrgId := RequestOrgIdFromContext(ctx)
	if requestOrgId != "" {
		logger = logger.With(RequestOrgIdFieldKey, requestOrgId)
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
