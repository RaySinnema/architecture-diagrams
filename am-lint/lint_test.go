package main

import (
	"strings"
	"testing"
)

func TestInvalidYaml(t *testing.T) {
	invalidYaml := `version 1.0`

	model, issues := Lint(invalidYaml)

	if model != nil {
		t.Errorf("Got model from invalid YAML")
	}
	if !hasIssue(issues, hasError("Invalid YAML")) {
		t.Errorf("No error about invalid YAML")
	}
}

func hasIssue(issues []Issue, test func(issue Issue) bool) bool {
	for _, issue := range issues {
		if test(issue) {
			return true
		}
	}
	return false
}

func hasError(message string) func(issue Issue) bool {
	return func(issue Issue) bool {
		return issue.level == Error && strings.Contains(issue.message, message)
	}
}

func TestVersion(t *testing.T) {
	yamlWithVersion := `version: 1.0`

	model, _ := Lint(yamlWithVersion)

	if model == nil {
		t.Errorf("Missing model")
	}
	if model.version != "1.0" {
		t.Errorf("Incorrect version: '%v'", model.version)
	}
}

func TestDefaultVersion(t *testing.T) {
	yamlWithoutVersion := ``

	model, _ := Lint(yamlWithoutVersion)

	if model.version != "1.0" {
		t.Errorf("Incorrect version: '%v'", model.version)
	}
}

func TestIncorrectVersion(t *testing.T) {
	yamlWithIncorrectVersion := `version: ape`

	model, issues := Lint(yamlWithIncorrectVersion)

	if !hasIssue(issues, hasError("Version must be a semantic version as defined by https://semver.org")) {
		t.Errorf("Accepts incorrect version")
	}
	if model.version != "" {
		t.Errorf("Model has version: '%v'", model.version)
	}
}

func TestFutureVersion(t *testing.T) {
	yamlWithFutureVersion := `version: 3.14`

	model, issues := Lint(yamlWithFutureVersion)

	if !hasIssue(issues, hasError("Undefined version")) {
		t.Errorf("Accepts future version")
	}
	if model.version != "" {
		t.Errorf("Model has version: '%v'", model.version)
	}
}
