package _const

import (
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

var (
	Tracer *tracesdk.TracerProvider
)