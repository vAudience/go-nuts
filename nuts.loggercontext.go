package gonuts

import (
	"context"

	"go.uber.org/zap"
)

const (
	RequestIdFieldKey = "requestId"
	requestIdPrefix   = "rid"
	requestIdLength   = 16
)

type ctxKey int

const (
	requestIdCtxKey ctxKey = iota
)

func GenerateRequestId() string {
	return NID(requestIdPrefix, requestIdLength)
}

func NewContextWithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, requestIdCtxKey, requestId)
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

func LoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	requestId := RequestIdFromContext(ctx)
	if requestId == "" {
		return L
	}
	return L.With(RequestIdFieldKey, requestId)
}
