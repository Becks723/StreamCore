package base

import (
	"github.com/kitex-contrib/obs-opentelemetry/provider"
)

// NewOtelProvider OTLP instrumentation https://github.com/kitex-contrib/obs-opentelemetry
func NewOtelProvider(serviceName string, exporterAddr string) provider.OtelProvider {
	return provider.NewOpenTelemetryProvider(
		provider.WithServiceName(serviceName),
		provider.WithExportEndpoint(exporterAddr),
		provider.WithInsecure())
}
