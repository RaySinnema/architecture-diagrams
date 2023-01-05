package main

type ModelPartReader interface {
	read(definition map[string]interface{}, fileName string, model *ArchitectureModel) []Issue
}
