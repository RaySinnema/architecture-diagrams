package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type State int64

const (
	Ok State = iota
	Emerging
	Review
	Revision
	Legacy
	Deprecated
)

func (s State) String() string {
	switch s {
	case Ok:
		return "OK"
	case Emerging:
		return "emerging"
	case Review:
		return "review"
	case Revision:
		return "revision"
	case Legacy:
		return "legacy"
	case Deprecated:
		return "deprecated"
	default:
		panic(fmt.Sprintf("Unknown state: %v", int64(s)))
	}
}

func (s State) Print(printer *Printer) {
	if s != Ok {
		printer.Print(" (", s, ")")
	}

}

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
