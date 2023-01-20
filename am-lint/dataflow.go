package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type DataFlow interface {
	setDirection(direction string)
}

const dataFlowField = "dataFlow"

func setDataFlow(fields map[string]*yaml.Node, dataFlow DataFlow) *Issue {
	direction, found, issue := stringFieldOf(fields, dataFlowField)
	if issue != nil {
		return issue
	}
	if found {
		allowed := []string{"send", "receive", "bidirectional"}
		if hasDifferentValueThan(direction, allowed) {
			return NodeError(fmt.Sprintf("Invalid %v: must be one of %v", dataFlowField, stringsIn(allowed)), fields[dataFlowField])
		}
		dataFlow.setDirection(direction)
	} else {
		dataFlow.setDirection("bidirectional")
	}
	return nil
}
