package main

import (
	"gopkg.in/yaml.v3"
)

type DataFlow int64

const (
	Send = iota
	Receive
	Bidirectional
)

type DataProcessor interface {
	setDataFlow(dataFlow DataFlow)
}

const dataFlowField = "dataFlow"
const defaultDataFlow = "bidirectional"

var allowedDataFlows = []string{"send", "receive", defaultDataFlow}

func setDataFlow(fields map[string]*yaml.Node, dataProcessor DataProcessor) *Issue {
	value, issue := enumFieldOf(fields, dataFlowField, allowedDataFlows, defaultDataFlow)
	if issue != nil {
		return issue
	}
	for index, dataFlow := range allowedDataFlows {
		if dataFlow == value {
			dataProcessor.setDataFlow(DataFlow(index))
		}
	}
	return nil
}
