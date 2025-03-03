package data

import (
	"fmt"
	"sort"
	"strings"

	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

type Container struct {
	Resource  *pcommon.Resource
	Metric    *pmetric.Metric
	DataPoint *pmetric.SummaryDataPoint
}

func NewContainer(resource *pcommon.Resource, metric *pmetric.Metric, dataPoint *pmetric.SummaryDataPoint) Container {
	return Container {
		resource,
		metric,
		dataPoint,
	}
}

func (container *Container) AsKey() string {
	resourceAttributes := container.ResourceAttributes()
	metricAttributes := container.MetricAttributes()
	dataPointAttributes := container.DataPointAttributes()

	return fmt.Sprintf("%s, %s, %s", resourceAttributes, metricAttributes, dataPointAttributes)
}

func (container *Container) ResourceAttributes() string {
	resourceAttributes := container.Resource.Attributes()

	values := make([]string, 0, resourceAttributes.Len())

	for key, value := range resourceAttributes.AsRaw() {
		values = append(values, fmt.Sprintf("%s: %s", key, value))
	}

	sort.Strings(values)

	return strings.Join(values, ", ")
}

func (container *Container) MetricAttributes() string {
	metricName := container.Metric.Name()
	metricType := container.Metric.Type()
	metricUnit := container.Metric.Unit()

	return fmt.Sprintf("name: %s, type: %s, unit: %s", metricName, metricType, metricUnit)
}

func (container *Container) DataPointAttributes() string {
	dataPointAttributes := container.DataPoint.Attributes()

	values := make([]string, 0, dataPointAttributes.Len())

	for key, value := range dataPointAttributes.AsRaw() {
		values = append(values, fmt.Sprintf("%s: %s", key, value))
	}

	sort.Strings(values)

	return strings.Join(values, ", ")
}

func (container *Container) GetMetricNamespace() string {
	resourceAttributes := container.Resource.Attributes()
	serviceNamespace, _ := resourceAttributes.Get("service.namespace")
	serviceName, _ := resourceAttributes.Get("service.name")
	return fmt.Sprintf("%s/%s", serviceNamespace.AsString(), serviceName.AsString())
}

func (container *Container) GetMetricName() string {
	return container.Metric.Name()
}
