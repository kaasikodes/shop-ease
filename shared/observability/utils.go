package observability

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func TraceInfoFromContext(ctx context.Context) (traceId, spanId string) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return "", ""
	}
	return spanCtx.TraceID().String(), spanCtx.SpanID().String()
}
