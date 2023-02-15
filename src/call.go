package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Call struct {
	node             *yaml.Node
	Description      string
	ExternalSystemId string          `yaml:"externalSystem,omitempty"`
	ExternalSystem   *ExternalSystem `yaml:",omitempty"`
	ServiceId        string          `yaml:"service,omitempty"`
	Service          *Service        `yaml:",omitempty"`
	DataFlow         DataFlow
	TechnologyIds    []string
	TechnologiesId   string
	Technologies     []*Technology
}

func (c *Call) getTechnologies() []*Technology {
	//TODO implement me
	panic("implement me")
}

func (c *Call) getNode() *yaml.Node {
	return c.node
}

func (c *Call) getTechnologyIds() []string {
	return c.TechnologyIds
}

func (c *Call) setTechnologyIds(technologies []string) {
	c.TechnologyIds = technologies
}

func (c *Call) getTechnologyBundleId() string {
	return c.TechnologiesId
}

func (c *Call) setTechnologyBundleId(technologyBundle string) {
	c.TechnologiesId = technologyBundle
}

func (c *Call) setTechnologies(technologies []*Technology) {
	c.Technologies = technologies
}

func (c *Call) getDescription() string {
	return c.Description
}

func (c *Call) setDescription(description string) {
	c.Description = description
}

const serviceField = "service"
const systemField = "externalSystem"

func (c *Call) setDataFlow(dataFlow DataFlow) {
	c.DataFlow = dataFlow
}

func (c *Call) read(node *yaml.Node) []Issue {
	c.node = node
	fields, issue := toMap(node)
	if issue != nil {
		return []Issue{*issue}
	}
	issues := make([]Issue, 0)
	issues = append(issues, c.readCallee(node, issue, fields)...)
	issues = append(issues, setDescription(fields, c)...)
	issues = append(issues, setDataFlow(node, fields, c)...)
	issues = append(issues, setTechnologies(fields, c)...)
	return issues
}

func (c *Call) readCallee(node *yaml.Node, issue *Issue, fields map[string]*yaml.Node) []Issue {
	issues := make([]Issue, 0)
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
	return issues
}

func (c *Call) Callee() interface{} {
	if c.Service != nil {
		return c.Service
	}
	return c.ExternalSystem
}
