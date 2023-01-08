package main

import (
	"gopkg.in/yaml.v3"
)

type ModelPartReader interface {
	read(definition *yaml.Node, fileName string, model *ArchitectureModel) []Issue
}

func stringValueOf(field string, node *yaml.Node) (string, *Issue) {
	if node.Kind == yaml.ScalarNode {
		return node.Value, nil
	}
	return "", NeedScalarError(field, node)

}
