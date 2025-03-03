package data

import (
	"fmt"
	"embed"
	"encoding/json"
)

//go:embed mappings/*
var mappings embed.FS

type Mapping struct {
	data map[string]NameSpace
}

type NameSpace struct {
	Name       string               `json:"name"`
	Metrics    map[string]Metric    `json:"metrics"`
	Dimensions map[string]Dimension `json:"dimensions"`
}

type Metric struct {
	Name string `json:"name"`
}

type Dimension struct {
	Name string `json:"name"`
}

func NewMapping() Mapping {
	mapping := Mapping{make(map[string]NameSpace)}

	files, _ := mappings.ReadDir("mappings")

	var nameSpcae NameSpace

	for _, file := range files {

		data, err := mappings.ReadFile("mappings/" + file.Name())

		if err != nil {
			fmt.Sprintf("Error reading %v", err)
			continue
		}

		if err := json.Unmarshal(data, &nameSpcae); err != nil {
			fmt.Sprintf("Error parsing JSON:", err)
			continue
		}

		mapping.data[nameSpcae.Name] = nameSpcae
	}

	return mapping
}

func (mapping *Mapping) GetMetricName(nameSpcace, metricName string) string {
	ns, _ := mapping.data[nameSpcace]
	metric, _ := ns.Metrics[metricName]
	return metric.Name
}
