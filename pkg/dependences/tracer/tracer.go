package tracer

import (
	"NotSmokeBot/config"
	"NotSmokeBot/pkg/constant"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"log"
)

func NewTracer(cfg config.Config) {
	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(cfg.OpenTelemetry.URL),
		),
	)
	if err != nil {
		log.Fatalf("Cannot create Jaeger exporter: %s", err.Error())
	}

	constant.Tracer = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.OpenTelemetry.ServiceName),
		)),
	)
	otel.SetTracerProvider(constant.Tracer)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		),
	)
}

const TraceIDHeader = "X-Trace-OrgAppID"

func StartFiberTrace(c *fiber.Ctx, spanName string) (context.Context, trace.Span) {
	ctx, span := otel.Tracer("").Start(c.Context(), spanName)
	traceID := span.SpanContext().TraceID().String()
	c.Response().Header.Set(TraceIDHeader, traceID)
	c.Locals(TraceIDHeader, traceID)
	return ctx, span
}

func SpanSetErrWrap(span trace.Span, err, errorValue error, errorPlace string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	errorInfo := fmt.Sprintf("%s|%s|%v|", errorPlace, errorValue, args)
	err = errors.Wrapf(err, errorInfo)
	span.SetStatus(codes.Error, err.Error())
	return err
}
