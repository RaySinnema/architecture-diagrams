package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type System struct {
	Name string
}

type SystemReader struct {
}

func (_ SystemReader) read(node *yaml.Node, fileName string, model *ArchitectureModel) []Issue {
	content := make(map[string]string)
	if node != nil {
		err := node.Decode(&content)
		if err != nil {
			return []Issue{*NodeError(fmt.Sprintf("invalid system: %v", err), node)}
		}
	}
	name, found := content["name"]
	if found {
		model.System.Name = name
	} else {
		model.System.Name = friendly(fileName)
	}
	return []Issue{}
}
