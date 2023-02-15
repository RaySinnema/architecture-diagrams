package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"sort"
)

type ExternalSystem struct {
	node        *yaml.Node
	Id          string
	Name        string
	Description string
	Type        string
	Calls       []*Call
}

func (es *ExternalSystem) Print(printer *Printer) {
	printer.Print(es.Name)
	if es.Type != "" {
		printer.Print(" (", es.Type, ")")
	}
	printer.NewLine()
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

func (es *ExternalSystem) getDescription() string {
	return es.Description
}

func (es *ExternalSystem) setDescription(description string) {
	es.Description = description
}

func (es *ExternalSystem) read(id string, node *yaml.Node) []Issue {
	var fields map[string]*yaml.Node
	fields, issues := namedObject(node, id, es)
	issues = append(issues, setDescription(fields, es)...)
	issues = append(issues, es.readType(fields)...)
	issues = append(issues, es.readCalls(fields)...)
	return issues
}

func (es *ExternalSystem) readType(fields map[string]*yaml.Node) []Issue {
	systemType, found, issue := stringFieldOf(fields, "type")
	if issue != nil {
		return []Issue{*issue}
	}
	if found {
		es.Type = systemType
	}
	return []Issue{}
}

func (es *ExternalSystem) readCalls(fields map[string]*yaml.Node) []Issue {
	callNodes, found, issue := sequenceFieldOf(fields, "calls")
	if issue != nil {
		return []Issue{*issue}
	}
	issues := make([]Issue, 0)
	calls := make([]*Call, 0)
	if found {
		for _, callNode := range callNodes {
			call := Call{}
			calls = append(calls, &call)
			issues = append(issues, call.read(callNode)...)
		}
	}
	es.Calls = calls
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
		externalSystems = append(externalSystems, &externalSystem)
		issues = append(issues, externalSystem.read(id, externalSystemNode)...)
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
			issues = append(issues, c.connectCall(call, model)...)
			issues = append(issues, connectTechnologies(call, model)...)
		}
	}
	return issues
}

func (c ExternalSystemConnector) connectCall(call *Call, model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	if call.ExternalSystemId != "" {
		for _, candidate := range model.ExternalSystems {
			if candidate.Id == call.ExternalSystemId {
				call.ExternalSystem = candidate
			}
		}
		if call.ExternalSystem == nil {
			issues = append(issues, *NodeError(fmt.Sprintf("Unknown external system '%v'", call.ExternalSystemId), call.node))
		}
	}
	if call.ServiceId != "" {
		for _, candidate := range model.Services {
			if candidate.Id == call.ServiceId {
				call.Service = candidate
			}
		}
		if call.Service == nil {
			issues = append(issues, *NodeError(fmt.Sprintf("Unknown service '%v'", call.ServiceId), call.node))
		}
	}
	return issues
}
