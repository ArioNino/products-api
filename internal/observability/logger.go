package observability

import (
	"context"
	"fmt"
	"log/slog"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/resource"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func InitLogger(serviceName string) (*slog.Logger, func(context.Context) error, error) {
	exporter, err := stdoutlog.New(
		stdoutlog.WithPrettyPrint(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal membuat stdout log exporter: %w", err)
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal membuat resource: %w", err)
	}

	processor := sdklog.NewSimpleProcessor(exporter)

	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(processor),
	)

	logger := otelslog.NewLogger(serviceName, otelslog.WithLoggerProvider(provider))

	shutdown := func(ctx context.Context) error {
		return provider.Shutdown(ctx)
	}

	return logger, shutdown, nil
}