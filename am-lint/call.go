package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Call struct {
	node             *yaml.Node
	Description      string
	ExternalSystemId string `yaml:"externalSystem,omitempty"`
	ServiceId        string `yaml:"service,omitempty"`
	DataFlow         string
	ExternalSystem   *ExternalSystem
}

const serviceField = "service"
const systemField = "externalSystem"
const dataFlowField = "dataFlow"

func (c *Call) read(node *yaml.Node, issues []Issue) []Issue {
	c.node = node
	fields, issue := toMap(node)
	if issue != nil {
		return append(issues, *issue)
	}
	serviceName, serviceFound, issue := stringFieldOf(fields, serviceField)
	if issue != nil {
		issues = append(issues, *issue)
	}
	systemName, systemFound, issue := stringFieldOf(fields, systemField)
	if issue != nil {
		issues = append(issues, *issue)
	}
	if serviceFound && systemFound {
		issues = append(issues, *NodeError(fmt.Sprintf("A call may be to either a %v or to an %v. Split the call into two to call both.", serviceField, systemField), node))
	} else if serviceFound {
		c.ServiceId = serviceName
	} else if systemFound {
		c.ExternalSystemId = systemName
	} else {
		issues = append(issues, *NodeError(fmt.Sprintf("One of %v or %v is required", serviceField, systemField), node))
	}
	description, found, issue := stringFieldOf(fields, "description")
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		c.Description = description
	}
	dataFlow, found, issue := stringFieldOf(fields, dataFlowField)
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		allowed := []string{"send", "receive", "bidirectional"}
		if hasDifferentValueThan(dataFlow, allowed) {
			issues = append(issues, *NodeError(fmt.Sprintf("Invalid %v: must be one of %v", dataFlowField, stringsIn(allowed)), fields[dataFlowField]))
		} else {
			c.DataFlow = dataFlow
		}
	} else {
		c.DataFlow = "bidirectional"
	}
	return issues
}

func stringsIn(values []string) string {
	result := ""
	for index, value := range values {
		if index == 0 {
			result = fmt.Sprintf("'%v'", value)
		} else if index == len(values)-1 {
			result = fmt.Sprintf("%v, or '%v'", result, value)
		} else {
			result = fmt.Sprintf("%v, '%v'", result, value)
		}
	}
	return result
}

func (c *Call) Callee() *ExternalSystem {
	if c.ExternalSystem != nil {
		return c.ExternalSystem
	}
	return nil
}
