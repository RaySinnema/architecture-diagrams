package main

import (
	"gopkg.in/yaml.v3"
)

type System struct {
	Name string
}

func (s *System) setNode(_ *yaml.Node) {
	// Nothing to do
}

func (s *System) setId(_ string) {
	// Nothing to do
}

func (s *System) setName(name string) {
	s.Name = name
}

type SystemReader struct {
}

func (_ SystemReader) read(node *yaml.Node, fileName string, model *ArchitectureModel) []Issue {
	fields, issue := toMap(node)
	if issue != nil {
		return []Issue{*issue}
	}
	return setName(fields, &model.System, fileName)
}
