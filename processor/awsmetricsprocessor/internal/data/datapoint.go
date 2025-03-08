package data

import (
	"fmt"
	"sort"
	"strings"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type DataPoints struct {
	slice *pmetric.SummaryDataPointSlice
}

func (dataPoints *DataPoints) Apply(f func (*DataPoint)) {
	for i := 0; i < dataPoints.slice.Len(); i++ {
		dataPoint := dataPoints.slice.At(i)
		f(&DataPoint{&dataPoint})
	}
}

type DataPoint struct {
	dataPoint *pmetric.SummaryDataPoint
}

func (dataPoint *DataPoint) DataPointAttributes() string {
	dataPointAttributes := dataPoint.dataPoint.Attributes()

	values := make([]string, 0, dataPointAttributes.Len())

	for key, value := range dataPointAttributes.AsRaw() {
		values = append(values, fmt.Sprintf("%s: %s", key, value))
	}

	sort.Strings(values)

	return strings.Join(values, ", ")
}

func (dataPoint *DataPoint) GetSum() float64 {
	return dataPoint.dataPoint.Sum()
}

func (dataPoint *DataPoint) SetSum(sum float64) {
	dataPoint.dataPoint.SetSum(sum)
}

func (dataPoint *DataPoint) GetCount() uint64 {
	return dataPoint.dataPoint.Count()
}

func (dataPoint *DataPoint) SetCount(count uint64) {
	dataPoint.dataPoint.SetCount(count)
}
