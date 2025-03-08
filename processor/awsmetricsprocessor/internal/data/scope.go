package data

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type Scopes struct {
	Slice *pmetric.ScopeMetricsSlice
}

func (scopes *Scopes) Apply(f func(*Scope)) {
	for i := 0; i < scopes.Slice.Len(); i++ {
		scopeMetrics := scopes.Slice.At(i)
		scope := scopeMetrics.Scope()
		metricSlice := scopeMetrics.Metrics()
		metrics := Metrics{&metricSlice}
		f(&Scope{&scope, &metrics})
	}
}

type Scope struct {
	scope   *pcommon.InstrumentationScope
	metrics *Metrics
}

func (scope *Scope) Metrics() *Metrics {
	return scope.metrics
}
