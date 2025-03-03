package awsmetricsprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
)

var typeStr = component.MustNewType("awsmetrics")

func NewFactory() processor.Factory {
	return processor.NewFactory(
		typeStr,
		createDefaultConfig,
		processor.WithMetrics(createMetricsProcessor, component.StabilityLevelAlpha))
}

func createMetricsProcessor(ctx context.Context, set processor.Settings, cfg component.Config, next consumer.Metrics) (processor.Metrics, error) {
	if config, ok := cfg.(*Config); ok {
		logger := set.Logger
		processor := NewProcessor(logger, next, config)
		return &processor, nil
	} else {
		return nil, fmt.Errorf("configuration parsing error")
	}
}
