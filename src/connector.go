package main

type Connector interface {
	connect(model *ArchitectureModel) []Issue
}
