package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"sort"
)

type ExternalSystem struct {
	node  *yaml.Node
	Id    string
	Name  string
	Type  string
	Calls []*Call
}

func (es *ExternalSystem) setNode(node *yaml.Node) {
	es.node = node
}

func (es *ExternalSystem) setId(id string) {
	es.Id = id
}

func (es *ExternalSystem) setName(name string) {
	es.Name = name
}

func (es *ExternalSystem) read(id string, node *yaml.Node, issues []Issue) []Issue {
	var fields map[string]*yaml.Node
	fields, issues = namedObject(id, node, es, issues)
	systemType, found, issue := stringFieldOf(fields, "type")
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		es.Type = systemType
	}
	callNodes, found, issue := sequenceFieldOf(fields, "calls")
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		calls := make([]*Call, 0)
		for _, callNode := range callNodes {
			call := Call{}
			issues = call.read(callNode, issues)
			calls = append(calls, &call)
		}
		es.Calls = calls
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
	externalSystems := make([]*ExternalSystem, 0)
	for id, externalSystemNode := range externalSystemsById {
		externalSystem := ExternalSystem{}
		issues = externalSystem.read(id, externalSystemNode, issues)
		externalSystems = append(externalSystems, &externalSystem)
	}
	sort.Slice(externalSystems, func(i, j int) bool {
		return externalSystems[i].Name < externalSystems[j].Name
	})
	model.ExternalSystems = externalSystems
	return issues
}

type ExternalSystemConnector struct {
}

func (c ExternalSystemConnector) connect(model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	for _, externalSystem := range model.ExternalSystems {
		for _, call := range externalSystem.Calls {
			issues = append(issues, connect(call, model)...)
		}
	}
	return issues
}

func connect(call *Call, model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	if call.ExternalSystemId != "" {
		for _, candidate := range model.ExternalSystems {
			if candidate.Id == call.ExternalSystemId {
				call.ExternalSystem = candidate
			}
		}
		if call.ExternalSystem == nil {
			issues = append(issues, *NodeError(fmt.Sprintf("Unknown external system %v", call.ExternalSystemId), call.node))
		}
	}
	return issues
}
