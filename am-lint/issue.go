package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Level int64

const (
	Error Level = iota
	Warning
)

func (l Level) String() string {
	switch l {
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	default:
		panic(fmt.Sprintf("Unknown level %v", int64(l)))
	}
}

type Issue struct {
	Level        Level
	Message      string
	Line, Column int
}

func (i Issue) String() string {
	return fmt.Sprintf("[%v, %v]: %v - %v", i.Line, i.Column, i.Level, i.Message)
}

func FileError(message string) *Issue {
	return &Issue{Level: Error, Message: message}
}

func NodeError(message string, node *yaml.Node) *Issue {
	return &Issue{Level: Error, Message: message, Line: node.Line, Column: node.Column}
}

func NeedTypeError(field string, node *yaml.Node, expectedType string) *Issue {
	return NodeError(fmt.Sprintf("%v must be a %v, not a %v", field, expectedType, kindToString(node.Kind)), node)
}

func kindToString(kind yaml.Kind) string {
	switch kind {
	case yaml.MappingNode:
		return "map"
	case yaml.SequenceNode:
		return "sequence"
	case yaml.ScalarNode:
		return "scalar"
	default:
		return fmt.Sprintf("Kind(%d)", kind)
	}
}

func NodeWarning(message string, node *yaml.Node) *Issue {
	return &Issue{Level: Warning, Message: message, Line: node.Line, Column: node.Column}
}
