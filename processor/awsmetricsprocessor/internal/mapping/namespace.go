package mapping

import (
	"fmt"
)

type Namespace struct {
	Name       string               `json:"name"`
	Prefix     string               `json:"prefix"`
	Metrics    map[string]Metric    `json:"metrics"`
	Dimensions map[string]Dimension `json:"dimensions"`
}

type Metric struct {
	Name string `json:"name"`
	Help string `json:"help"`
}

type Dimension struct {
	Name string `json:"name"`
}

func (namespace *Namespace) GetMetricName(metricName string) (string, error) {
	if metric, ok := namespace.Metrics[metricName]; ok {
		return fmt.Sprintf("%s_%s", namespace.Prefix, metric.Name), nil
	}

	return "", fmt.Errorf("Unknown metric %s", metricName)
}

func (namespace *Namespace) GetMetricHelp(metricName string) (string, error) {
	if metric, ok := namespace.Metrics[metricName]; ok {
		return metric.Help, nil
	}

	return "", fmt.Errorf("Unknown metric %s", metricName)
}

func (namespace *Namespace) GetDimensionName(dimensionName string) (string, error) {
	if dimension, ok := namespace.Dimensions[dimensionName]; ok {
		return dimension.Name, nil
	}

	return "", fmt.Errorf("Unknown dimension %s", dimensionName)
}
