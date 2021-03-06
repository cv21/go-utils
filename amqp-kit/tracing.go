package amqp_kit

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/opentracing/opentracing-go"
	otext "github.com/opentracing/opentracing-go/ext"
)

const tagError = "error"

type amqpSpanCtx string

var amqpCtx amqpSpanCtx

// Get context value spanContext and start Span with given operationName.
// Set an error as tag if raised.
func TraceEndpoint(tracer opentracing.Tracer, operationName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			var sp opentracing.Span
			if spCtx, ok := ctx.Value(amqpCtx).(*opentracing.SpanContext); ok {
				sp = tracer.StartSpan(operationName, opentracing.FollowsFrom(*spCtx))
				defer sp.Finish()

				otext.SpanKindRPCServer.Set(sp)
				ctx = opentracing.ContextWithSpan(ctx, sp)
			}

			i, err := next(ctx, request)
			if err != nil && sp != nil {
				sp.SetTag(tagError, err.Error())
			}

			return i, err
		}
	}
}
