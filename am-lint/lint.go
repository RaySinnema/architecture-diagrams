package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var readers = []ModelPartReader{VersionReader{}, SystemReader{}}

func LintText(text string) (*ArchitectureModel, []Issue) {
	model, issues := lint(text, "")
	return model, issues
}

func lint(definition string, fileName string) (*ArchitectureModel, []Issue) {
	m := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(definition), &m); err != nil {
		return nil, []Issue{*NewError("Invalid YAML")}
	}
	issues := make([]Issue, 0)
	model := ArchitectureModel{}

	for _, reader := range readers {
		issues = append(issues, reader.read(m, fileName, &model)...)
	}

	return &model, issues
}

func LintFile(fileName string) (*ArchitectureModel, []Issue) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, []Issue{*NewError(fmt.Sprintf("Couldn't read file %s: %v", fileName, err))}
	}
	model, issues := lint(string(bytes), fileName)
	return model, issues
}
