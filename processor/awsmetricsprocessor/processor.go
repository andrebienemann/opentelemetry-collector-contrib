package awsmetricsprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/awsmetricsprocessor/internal/data"
)

type Processor struct {
	host    component.Host
	cancel  context.CancelFunc
	logger  *zap.Logger
	next    consumer.Metrics
	config  *Config
	cache   data.Cache
	mapping data.Mapping
}

func NewProcessor(logger *zap.Logger, next consumer.Metrics, config *Config) Processor {
	return Processor{logger: logger, next: next, config: config}
}

func (processor *Processor) Start(ctx context.Context, host component.Host) error {
	processor.host = host

	ctx = context.Background()
	ctx, processor.cancel = context.WithCancel(ctx)

	processor.cache = data.NewCache()
	processor.mapping = data.NewMapping()

	return nil
}

func (processor *Processor) ConsumeMetrics(ctx context.Context, ms pmetric.Metrics) error {
	resourceMetricsSlice := ms.ResourceMetrics()
	for p := 0; p < resourceMetricsSlice.Len(); p++ {
		resourceMetrics := resourceMetricsSlice.At(p)
		resource := resourceMetrics.Resource()
		scopeMetricsSlice := resourceMetrics.ScopeMetrics()
		for q := 0; q < scopeMetricsSlice.Len(); q++ {
			scopeMetrics := scopeMetricsSlice.At(q)
			metricSlice := scopeMetrics.Metrics()
			for r := 0; r < metricSlice.Len(); r++ {
				metric := metricSlice.At(r)
				summary := metric.Summary()
				summaryDataPointSlice := summary.DataPoints()
				for s := 0; s < summaryDataPointSlice.Len(); s++ {
					dataPoint := summaryDataPointSlice.At(s)
					container := data.NewContainer(&resource, &metric, &dataPoint)
					processor.Display(&container)
					processor.Update(&container)
					processor.Translate(&container)
				}
			}
		}
	}

	return processor.next.ConsumeMetrics(ctx, ms)
}

func (processor *Processor) Update(container *data.Container) {
	key := container.AsKey()

	sum := container.DataPoint.Sum()
	count := container.DataPoint.Count()

	if record, ok := processor.cache.GetRecord(key); !ok {
		record := data.NewRecord(sum, count)
		processor.cache.PutRecord(key, &record)
	} else {
		container.DataPoint.SetSum(record.IncSum(sum))
		container.DataPoint.SetCount(record.IncCount(count))
	}
}

func (processor *Processor) Translate(container *data.Container) {
	nameSpace := container.GetMetricNamespace()
	metricName := container.GetMetricName()
	canonicalName := processor.mapping.GetMetricName(nameSpace, metricName)
	container.Metric.SetName(canonicalName)
}

func (processor *Processor) Display(container *data.Container) {
	resourceAttributes := container.ResourceAttributes()
	processor.logger.Info(resourceAttributes)
	metricAttributes := container.MetricAttributes()
	processor.logger.Info(metricAttributes)
	dataPointAttributes := container.DataPointAttributes()
	processor.logger.Info(dataPointAttributes)
}

func (processor *Processor) Shutdown(ctx context.Context) error {
	if processor.cancel != nil {
		processor.cancel()
	}

	return nil
}

func (processor *Processor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}
