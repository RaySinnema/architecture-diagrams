package main

import (
	"gopkg.in/yaml.v3"
)

type DataFlow interface {
	setDirection(direction string)
}

const dataFlowField = "dataFlow"
const defaultDataFlow = "bidirectional"

var allowedDataFlows = []string{"send", "receive", defaultDataFlow}

func setDataFlow(fields map[string]*yaml.Node, dataFlow DataFlow) *Issue {
	direction, issue := enumFieldOf(fields, dataFlowField, allowedDataFlows, defaultDataFlow)
	if issue != nil {
		return issue
	}
	dataFlow.setDirection(direction)
	return nil
}
