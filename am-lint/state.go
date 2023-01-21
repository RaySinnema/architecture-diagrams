package main

import "gopkg.in/yaml.v3"

type State int64

const (
	Ok State = iota
	Emerging
	Review
	Revision
	Legacy
	Deprecated
)

type Evolveable interface {
	setState(state State)
}

const stateField = "state"
const defaultState = "ok"

var allowedStates = []string{defaultState, "emerging", "review", "revision", "legacy", "deprecated"}

func setState(fields map[string]*yaml.Node, e Evolveable) *Issue {
	value, issue := enumFieldOf(fields, stateField, allowedStates, defaultState)
	if issue != nil {
		return issue
	}
	for index, state := range allowedStates {
		if state == value {
			e.setState(State(index))
		}
	}
	return nil
}
