package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Call struct {
	Description        string
	ExternalSystemName string `yaml:"externalSystem,omitempty"`
	ServiceName        string `yaml:"service,omitempty"`
	DataFlow           string
}

const serviceField = "service"
const systemField = "externalSystem"
const dataFlowField = "dataFlow"

func (c *Call) read(node *yaml.Node, issues []Issue) []Issue {
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
		issues = append(issues, *NodeError(fmt.Sprintf("Only one of '%v' and '%v' is allowed. If both are called, use two calls", serviceField, systemField), node))
	} else if serviceFound {
		c.ServiceName = serviceName
	} else if systemFound {
		c.ExternalSystemName = systemName
	} else {
		issues = append(issues, *NodeError(fmt.Sprintf("One of '%v' or '%v' is required", serviceField, systemField), node))
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
			issues = append(issues, *NodeError(fmt.Sprintf("Invalid '%v': must be one of %v", dataFlowField, allowed), fields[dataFlowField]))
		} else {
			c.DataFlow = dataFlow
		}
	} else {
		c.DataFlow = "bidirectional"
	}
	return issues
}

func (c *Call) Callee() string {
	if c.ServiceName != "" {
		return c.ServiceName
	}
	if c.ExternalSystemName != "" {
		return c.ExternalSystemName
	}
	return "???"
}
