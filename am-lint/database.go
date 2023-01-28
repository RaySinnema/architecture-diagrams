package main

import "gopkg.in/yaml.v3"

type DatabaseReader struct {
}

func (d DatabaseReader) read(node *yaml.Node, _ string, model *ArchitectureModel) []Issue {
	databases, issues := DataStoreReader{}.read(node)
	model.Databases = databases
	return issues
}

type DatabaseConnector struct {
}

func (d DatabaseConnector) connect(model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	for _, database := range model.Databases {
		issues = append(issues, connectTechnologies(database, model)...)
		issues = append(issues, connectTechnologies(ApiTechnologies{database}, model)...)
	}
	return issues
}
