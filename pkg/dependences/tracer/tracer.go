package tracer

import (
	"NotSmokeBot/config"
	"NotSmokeBot/pkg/constant/const"
	"fmt"
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

	_const.Tracer = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.OpenTelemetry.ServiceName),
		)),
	)
	otel.SetTracerProvider(_const.Tracer)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		),
	)
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
