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

func setDataFlow(owner *yaml.Node, fields map[string]*yaml.Node, dataProcessor DataProcessor) []Issue {
	const defaultDataFlow = "bidirectional"
	var allowedDataFlows = []string{"send", "receive", defaultDataFlow}

	value, issue := enumFieldOf(owner, fields, "dataFlow", allowedDataFlows, defaultDataFlow)
	if issue != nil {
		return []Issue{*issue}
	}
	for index, dataFlow := range allowedDataFlows {
		if dataFlow == value {
			dataProcessor.setDataFlow(DataFlow(index))
		}
	}
	return []Issue{}
}
