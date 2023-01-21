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

func (c *Call) setDirection(direction string) {
	c.DataFlow = direction
}

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
	issue = setDataFlow(fields, c)
	if issue != nil {
		issues = append(issues, *issue)
	}
	return issues
}

func (c *Call) Callee() *ExternalSystem {
	if c.ExternalSystem != nil {
		return c.ExternalSystem
	}
	return nil
}
