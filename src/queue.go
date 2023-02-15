package main

import (
	"gopkg.in/yaml.v3"
)

type QueueReader struct {
}

func (q QueueReader) read(node *yaml.Node, _ string, model *ArchitectureModel) []Issue {
	queues, issues := DataStoreReader{}.read(node)
	model.Queues = queues
	return issues
}

type QueueConnector struct {
}

func (d QueueConnector) connect(model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	for _, queue := range model.Queues {
		issues = append(issues, connectTechnologies(queue, model)...)
		issues = append(issues, connectTechnologies(ApiTechnologies{queue}, model)...)
	}
	return issues
}
