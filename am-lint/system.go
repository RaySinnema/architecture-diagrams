package main

import (
	"gopkg.in/yaml.v3"
)

type System struct {
	Name string
}

type SystemReader struct {
}

const nameField = "name"

func (_ SystemReader) read(node *yaml.Node, fileName string, model *ArchitectureModel) []Issue {
	if node != nil {
		details, issue := toMap(node)
		if issue != nil {
			return []Issue{*issue}
		}
		name, found, issue := stringFieldOf(details, nameField)
		if issue != nil {
			return []Issue{*issue}
		}
		if found {
			model.System.Name = name
			return []Issue{}
		}
	}
	model.System.Name = friendlyNameFrom(fileName)
	return []Issue{}
}
