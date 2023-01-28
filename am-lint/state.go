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

type Evolvable interface {
	setState(state State)
}

func setState(owner *yaml.Node, fields map[string]*yaml.Node, e Evolvable) []Issue {
	const defaultState = "ok"
	var allowedStates = []string{defaultState, "emerging", "review", "revision", "legacy", "deprecated"}

	value, issue := enumFieldOf(owner, fields, "state", allowedStates, defaultState)
	if issue != nil {
		return []Issue{*issue}
	}
	for index, state := range allowedStates {
		if state == value {
			e.setState(State(index))
		}
	}
	return []Issue{}
}
