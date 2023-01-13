package main

import (
	"gopkg.in/yaml.v3"
	"sort"
)

type ExternalSystem struct {
	Name  string
	Type  string
	Calls []Call
}

func (s *ExternalSystem) setName(name string) {
	s.Name = name
}

func (s *ExternalSystem) read(id string, node *yaml.Node, issues []Issue) []Issue {
	var fields map[string]*yaml.Node
	fields, issues = namedObject(id, node, s, issues)
	systemType, found, issue := stringFieldOf(fields, "type")
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		s.Type = systemType
	}
	callNodes, found, issue := sequenceFieldOf(fields, "calls")
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		calls := make([]Call, 0)
		for _, callNode := range callNodes {
			call := Call{}
			issues = call.read(callNode, issues)
			calls = append(calls, call)
		}
		s.Calls = calls
	}
	return issues
}

type ExternalSystemReader struct {
}

func (e ExternalSystemReader) read(node *yaml.Node, _ string, model *ArchitectureModel) []Issue {
	if node == nil {
		return []Issue{}
	}
	externalSystemsById, issue := toMap(node)
	if issue != nil {
		return []Issue{*issue}
	}
	issues := make([]Issue, 0)
	externalSystems := make([]ExternalSystem, 0)
	for id, externalSystemNode := range externalSystemsById {
		externalSystem := ExternalSystem{}
		issues = externalSystem.read(id, externalSystemNode, issues)
		externalSystems = append(externalSystems, externalSystem)
	}
	sort.Slice(externalSystems, func(i, j int) bool {
		return externalSystems[i].Name < externalSystems[j].Name
	})
	model.ExternalSystems = externalSystems
	return issues
}
