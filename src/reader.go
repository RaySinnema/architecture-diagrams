package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type ModelPartReader interface {
	read(definition *yaml.Node, fileName string, model *ArchitectureModel) []Issue
}

func mapFieldOf(fields map[string]*yaml.Node, field string) (map[string]*yaml.Node, bool, *Issue) {
	var result map[string]*yaml.Node
	var issue *Issue
	node, found := fields[field]
	if found {
		result, issue = toMap(node)
	} else {
		result = make(map[string]*yaml.Node)
		issue = nil
	}
	return result, found, issue
}

func toMap(node *yaml.Node) (map[string]*yaml.Node, *Issue) {
	result := make(map[string]*yaml.Node)
	if node == nil || node.IsZero() {
		return result, nil
	}
	if node.Kind != yaml.MappingNode {
		return nil, NodeError("Expected a map", node)
	}
	var tag string
	for index, child := range node.Content {
		if index%2 == 0 { // Key
			tag = child.Value
		} else { // Value
			result[tag] = child
		}
	}
	return result, nil
}

func sequenceFieldOf(fields map[string]*yaml.Node, field string) ([]*yaml.Node, bool, *Issue) {
	var result []*yaml.Node
	var issue *Issue
	node, found := fields[field]
	if found {
		result, issue = toSequence(node, field)
	} else {
		result = make([]*yaml.Node, 0)
		issue = nil
	}
	return result, found, issue
}

func toSequence(node *yaml.Node, field string) ([]*yaml.Node, *Issue) {
	if node.Kind == yaml.SequenceNode {
		return node.Content, nil
	}
	return []*yaml.Node{}, NeedTypeError(field, node, "sequence")
}

func stringFieldOf(fields map[string]*yaml.Node, field string) (string, bool, *Issue) {
	node, found := fields[field]
	if found {
		result, issue := toString(node, field)
		if issue != nil {
			return "", true, issue
		}
		return result, true, nil
	}
	return "", false, nil
}

func toString(node *yaml.Node, field string) (string, *Issue) {
	if node.Kind == yaml.ScalarNode {
		return node.Value, nil
	}
	return "", NeedTypeError(field, node, "string")

}

func enumFieldOf(owner *yaml.Node, fields map[string]*yaml.Node, field string, allowed []string, defaultValue string) (string, *Issue) {
	value, found, issue := stringFieldOf(fields, field)
	if issue != nil {
		return "", issue
	}
	if !found {
		if hasDifferentValueThan(defaultValue, allowed) {
			return "", NodeError(fmt.Sprintf("Missing required field %v", field), owner)
		}
		value = defaultValue
	}
	if hasDifferentValueThan(value, allowed) {
		return "", NodeError(fmt.Sprintf("Invalid %v: must be one of %v", field, stringsIn(allowed)), fields[field])
	}
	return value, nil
}

func hasDifferentValueThan(value string, allowed []string) bool {
	for _, ok := range allowed {
		if ok == value {
			return false
		}
	}
	return true
}

func stringsIn(values []string) string {
	result := ""
	for index, value := range values {
		if index == 0 {
			result = fmt.Sprintf("'%v'", value)
		} else if index == len(values)-1 {
			result = fmt.Sprintf("%v, or '%v'", result, value)
		} else {
			result = fmt.Sprintf("%v, '%v'", result, value)
		}
	}
	return result
}
