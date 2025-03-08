package data

import (
	"fmt"
	"sort"
	"strings"

	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

type Resources struct {
	slice *pmetric.ResourceMetricsSlice
}

func From(metrics *pmetric.Metrics) *Resources {
	resourceMetricsSlice := metrics.ResourceMetrics()
	return &Resources{&resourceMetricsSlice}
}

func (resources *Resources) Apply(f func (*Resource)) {
	for i := 0; i < resources.slice.Len(); i++ {
		resourceMetrics := resources.slice.At(i)
		resource := resourceMetrics.Resource()
		scopeMetricsSlice := resourceMetrics.ScopeMetrics()
		f(&Resource{&resource, &Scopes{&scopeMetricsSlice}})
	}
}

type Resource struct {
	resource *pcommon.Resource
	scopes   *Scopes
}

func (resource *Resource) ResourceAttributes() string {
	resourceAttributes := resource.resource.Attributes()

	values := make([]string, 0, resourceAttributes.Len())

	for key, value := range resourceAttributes.AsRaw() {
		values = append(values, fmt.Sprintf("%s: %s", key, value))
	}

	sort.Strings(values)

	return strings.Join(values, ", ")
}

func (resource *Resource) GetNamespace() string {
	resourceAttributes := resource.resource.Attributes()
	serviceNamespace, _ := resourceAttributes.Get("service.namespace")
	serviceName, _ := resourceAttributes.Get("service.name")
	return fmt.Sprintf("%s/%s", serviceNamespace.AsString(), serviceName.AsString())
}

func (resource *Resource) Scopes() *Scopes {
	return resource.scopes
}
