package mapping

import (
	"fmt"
	"embed"
	"encoding/json"
)

//go:embed json/*
var mappings embed.FS

type Mapping struct {
	data map[string]Namespace
}

func NewMapping() (Mapping, []error) {
	errs := make([]error, 0)
	data := make(map[string]Namespace)
	mapping := Mapping{data}

	files, _ := mappings.ReadDir("json")

	var nameSpcae Namespace

	for _, file := range files {
		fileName := file.Name()

		fileContent, err := mappings.ReadFile("json/" + fileName)

		if err != nil {
			errs = append(errs, fmt.Errorf("Error reading %s", fileName))
			continue
		}

		if err := json.Unmarshal(fileContent, &nameSpcae); err != nil {
			errs = append(errs, fmt.Errorf("Error parsing %s", fileName))
			continue
		}

		mapping.data[nameSpcae.Name] = nameSpcae
	}

	return mapping, errs
}

func (mapping *Mapping) GetNamespace(name string) (Namespace, error) {
	if namespace, ok := mapping.data[name]; ok {
		return namespace, nil
	}

	return Namespace{}, fmt.Errorf("Unknown namespace %s", name)
}
