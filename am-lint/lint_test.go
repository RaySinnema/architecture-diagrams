package main

import (
	"os"
	"strings"
	"testing"
)

func TestInvalidYaml(t *testing.T) {
	invalidYaml := `version 1.0`

	model, issues := LintText(invalidYaml)

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

	model, _ := LintText(yamlWithVersion)

	if model == nil {
		t.Errorf("Missing model")
	}
	if model.version != "1.0" {
		t.Errorf("Incorrect version: '%v'", model.version)
	}
}

func TestDefaultVersion(t *testing.T) {
	yamlWithoutVersion := ``

	model, _ := LintText(yamlWithoutVersion)

	if model.version != "1.0" {
		t.Errorf("Incorrect version: '%v'", model.version)
	}
}

func TestIncorrectVersion(t *testing.T) {
	yamlWithIncorrectVersion := `version: ape`

	model, issues := LintText(yamlWithIncorrectVersion)

	if !hasIssue(issues, hasError("Version must be a semantic version as defined by https://semver.org")) {
		t.Errorf("Accepts incorrect version")
	}
	if model.version != "" {
		t.Errorf("Model has version: '%v'", model.version)
	}
}

func TestFutureVersion(t *testing.T) {
	yamlWithFutureVersion := `version: 3.14`

	model, issues := LintText(yamlWithFutureVersion)

	if !hasIssue(issues, hasError("Undefined version")) {
		t.Errorf("Accepts future version")
	}
	if model.version != "" {
		t.Errorf("Model has version: '%v'", model.version)
	}
}

func TestSystemName(t *testing.T) {
	yamlWithFutureVersion := `system:
  name: foo`

	model, _ := LintText(yamlWithFutureVersion)

	if model.system.name != "foo" {
		t.Errorf("Model has system: '%v'", model.system.name)
	}
}

func TestDefaultSystemName(t *testing.T) {
	fileName := "a-system-named-foo.yaml"
	if err := os.WriteFile(fileName, []byte(`version: 1.0`), 0644); err != nil {
		t.Errorf("Failed to create test file")
	}
	defer func() {
		if err := os.Remove(fileName); err != nil {
			t.Errorf("Failed to delete test file")
		}
	}()

	model, _ := LintFile(fileName)

	if model.system.name != "A system named foo" {
		t.Errorf("Model has system: '%v'", model.system.name)
	}
}

func TestEmptySystemName(t *testing.T) {
	yamlWithFutureVersion := `system:
  name: ''`

	_, issues := LintText(yamlWithFutureVersion)

	if !hasIssue(issues, hasError("Missing system name")) {
		t.Errorf("No error about missing system name")
	}
}

func TestNonStringSystem(t *testing.T) {
	yamlWithFutureVersion := `system:
  name: 3.14`

	model, _ := LintText(yamlWithFutureVersion)

	if model.system.name != "3.14" {
		t.Errorf("Model has system: '%v'", model.system.name)
	}
}
