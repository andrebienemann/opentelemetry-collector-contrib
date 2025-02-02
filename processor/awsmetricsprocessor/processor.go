package awsmetricsprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type MetricsProcessor struct {
	host   component.Host
	cancel context.CancelFunc
	logger *zap.Logger
	next   consumer.Metrics
	config *Config
}

func (processor *MetricsProcessor) Start(ctx context.Context, host component.Host) error {
	processor.host = host
	ctx = context.Background()
	ctx, processor.cancel = context.WithCancel(ctx)
	processor.logger.Debug("Starting processing metrics")

	return nil
}

func (processor *MetricsProcessor) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	processor.logger.Debug("Received Metrics")
	return processor.next.ConsumeMetrics(ctx, md)
}

func (processor *MetricsProcessor) Shutdown(ctx context.Context) error {
	if processor.cancel != nil {
		processor.cancel()
	}

	return nil
}

func (processor *MetricsProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}
