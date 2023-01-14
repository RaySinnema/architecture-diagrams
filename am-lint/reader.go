package main

import (
	"gopkg.in/yaml.v3"
)

type ModelPartReader interface {
	read(definition *yaml.Node, fileName string, model *ArchitectureModel) []Issue
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

func toString(node *yaml.Node, field string) (string, *Issue) {
	if node.Kind == yaml.ScalarNode {
		return node.Value, nil
	}
	return "", NeedTypeError(field, node, "string")

}

func hasDifferentValueThan(value string, allowed []string) bool {
	for _, ok := range allowed {
		if ok == value {
			return false
		}
	}
	return true
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

func sequenceFieldOf(fields map[string]*yaml.Node, field string) ([]*yaml.Node, bool, *Issue) {
	result := make([]*yaml.Node, 0)
	node, found := fields[field]
	if found {
		result, issue := toSequence(node)
		if issue != nil {
			return result, false, issue
		}
		return result, true, nil
	}
	return result, false, nil
}

func toSequence(node *yaml.Node) ([]*yaml.Node, *Issue) {
	if node.Kind == yaml.SequenceNode {
		return node.Content, nil
	}
	return []*yaml.Node{}, NodeError("Expected a sequence", node)
}
