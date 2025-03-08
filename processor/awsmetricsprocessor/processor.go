package awsmetricsprocessor

import (
	"fmt"
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/awsmetricsprocessor/internal/data"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/awsmetricsprocessor/internal/cache"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/awsmetricsprocessor/internal/mapping"
)

type Processor struct {
	host    component.Host
	cancel  context.CancelFunc
	logger  *zap.Logger
	next    consumer.Metrics
	config  *Config
	cache   cache.Cache
	mapping mapping.Mapping
}

func NewProcessor(logger *zap.Logger, next consumer.Metrics, config *Config) Processor {
	return Processor{logger: logger, next: next, config: config}
}

func (processor *Processor) Start(ctx context.Context, host component.Host) error {
	processor.host = host

	ctx = context.Background()
	ctx, processor.cancel = context.WithCancel(ctx)

	processor.cache = cache.NewCache()

	mapping, _ := mapping.NewMapping()
	processor.mapping = mapping

	return nil
}

func (processor *Processor) ConsumeMetrics(ctx context.Context, ms pmetric.Metrics) error {
	resources := data.From(&ms)
	resources.Apply(func (resource *data.Resource) {
		resource.Scopes().Apply(func (scope *data.Scope) {
			scope.Metrics().Apply(func (metric *data.Metric) {
				metric.DataPoints().Apply(func (dataPoint *data.DataPoint){
					processor.Display(resource, metric, dataPoint)
					processor.Update(resource, metric, dataPoint)
				})
				processor.AddHelp(resource, metric)
				processor.Translate(resource, metric)
			})
		})
	})

	return processor.next.ConsumeMetrics(ctx, ms)
}

func (processor *Processor) Update(resource *data.Resource, metric *data.Metric, dataPoint *data.DataPoint) {
	resourceAttributes := resource.ResourceAttributes()
	metricAttributes := metric.MetricAttributes()
	dataPointAttributes := dataPoint.DataPointAttributes()

	key := fmt.Sprintf("%s, %s, %s", resourceAttributes, metricAttributes, dataPointAttributes)

	sum := dataPoint.GetSum()
	count := dataPoint.GetCount()

	if record, ok := processor.cache.GetRecord(key); !ok {
		record := cache.NewRecord(sum, count)
		processor.cache.PutRecord(key, &record)
	} else {
		dataPoint.SetSum(record.IncrSum(sum))
		dataPoint.SetCount(record.IncrCount(count))
	}
}

func (processor *Processor) Translate(resource *data.Resource, metric *data.Metric) {
	namespace, err := processor.mapping.GetNamespace(resource.GetNamespace())

	if err != nil {
		processor.logger.Info(err.Error())
		return
	}

	metricName, err := namespace.GetMetricName(metric.GetName())

	if err != nil {
		processor.logger.Info(err.Error())
		return
	}

	metric.SetName(metricName)
}

func (processor *Processor) AddHelp(resource *data.Resource, metric *data.Metric) {
	namespace, err := processor.mapping.GetNamespace(resource.GetNamespace())

	if err != nil {
		processor.logger.Info(err.Error())
		return
	}

	help, err := namespace.GetMetricHelp(metric.GetName())

	if err != nil {
		processor.logger.Info(err.Error())
		return
	}

	metric.SetDescription(help)
}

func (processor *Processor) Display(resource *data.Resource, metric *data.Metric, dataPoint *data.DataPoint) {
	resourceAttributes := resource.ResourceAttributes()
	processor.logger.Info(resourceAttributes)
	metricAttributes := metric.MetricAttributes()
	processor.logger.Info(metricAttributes)
	dataPointAttributes := dataPoint.DataPointAttributes()
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
