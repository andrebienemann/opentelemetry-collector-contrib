package data

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type Metrics struct {
	Slice *pmetric.MetricSlice
}

func (metrics *Metrics) Apply(f func(*Metric)) {
	for i := 0; i < metrics.Slice.Len(); i++ {
		metric := metrics.Slice.At(i)
		summary := metric.Summary()
		summaryDataPointSlice := summary.DataPoints()
		dataPoints := DataPoints{&summaryDataPointSlice}
		f(&Metric{&metric, &dataPoints})
	}
}

type Metric struct {
	metric     *pmetric.Metric
	dataPoints *DataPoints
}

func (metric *Metric) DataPoints() *DataPoints {
	return metric.dataPoints
}

func (metric *Metric) MetricAttributes() string {
	metricName := metric.metric.Name()
	metricType := metric.metric.Type()
	metricUnit := metric.metric.Unit()
	return fmt.Sprintf("name: %s, type: %s, unit: %s", metricName, metricType, metricUnit)
}

func (metric *Metric) GetName() string {
	return metric.metric.Name()
}

func (metric *Metric) SetName(name string) {
	metric.metric.SetName(name)
}

func (metric *Metric) SetDescription(description string) {
	metric.metric.SetDescription(description)
}
