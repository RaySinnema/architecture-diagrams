package main

import (
	"os"
	"strings"
	"testing"
)

func TestInvalidYaml(t *testing.T) {
	invalidYaml := `version 1.0`

	model, issues := LintText(invalidYaml)

	if !hasIssue(issues, hasError("Invalid YAML")) {
		t.Errorf("No error about invalid YAML")
	}
	if model != nil {
		t.Errorf("Model from invalid YAML")
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
		return issue.Level == Error && strings.Contains(issue.Message, message)
	}
}

func TestUnknownTopLevelElement(t *testing.T) {
	yamlWithUnknownTopLevelElement := `ape: bear`

	_, issues := LintText(yamlWithUnknownTopLevelElement)

	if !hasIssue(issues, hasWarning("")) {
		t.Errorf("No warning about unknown top-level element")
	}
}

func hasWarning(message string) func(issue Issue) bool {
	return func(issue Issue) bool {
		return issue.Level == Warning && strings.Contains(issue.Message, message)
	}
}

func TestUnknownFile(t *testing.T) {
	fileName := "a-system-named-foo.yaml"
	if err := os.WriteFile(fileName, []byte(`version: 1.0`), 0200); err != nil {
		t.Errorf("Failed to create test file")
	}
	defer func() {
		if err := os.Remove(fileName); err != nil {
			t.Errorf("Failed to delete test file")
		}
	}()

	_, issues := LintFile(fileName)

	if len(issues) == 0 {
		t.Errorf("No error when unable to read file")
	}
}

func TestVersion(t *testing.T) {
	yamlWithVersion := `version: 1.0`

	model, _ := LintText(yamlWithVersion)

	if model == nil {
		t.Fatalf("Missing model")
	}
	if model.Version != "1.0" {
		t.Errorf("Incorrect version: '%v'", model.Version)
	}
}

func TestDefaultVersion(t *testing.T) {
	yamlWithoutVersion := ``

	model, _ := LintText(yamlWithoutVersion)

	if model == nil {
		t.Fatalf("Missing model")
	}
	if model.Version != "1.0.0" {
		t.Errorf("Incorrect version: '%v'", model.Version)
	}
}

func TestIncorrectVersion(t *testing.T) {
	yamlWithIncorrectVersion := `version: ape\n`

	model, issues := LintText(yamlWithIncorrectVersion)

	if !hasIssue(issues, hasError("Version must be a semantic version as defined by https://semver.org")) {
		t.Errorf("No error about incorrect version")
	}
	if model.Version != "" {
		t.Errorf("Model has version: '%v'", model.Version)
	}
}

func TestIncorrectVersionStructure(t *testing.T) {
	yamlWithIncorrectVersionStructure := `version:
  ape: bear
`

	model, issues := LintText(yamlWithIncorrectVersionStructure)

	if !hasIssue(issues, hasError("version must be a string, not a map")) {
		t.Errorf("Missing error")
	}
	if model.Version != "" {
		t.Errorf("Model has version: '%v'", model.Version)
	}
}

func TestFutureVersion(t *testing.T) {
	yamlsWithFutureVersion := []string{
		`version: 3.14`,
		`version: 1.2`,
		`version: 1.0.1`,
	}

	for _, yamlWithFutureVersion := range yamlsWithFutureVersion {
		model, issues := LintText(yamlWithFutureVersion)

		if !hasIssue(issues, hasError("Undefined version")) {
			t.Errorf("Accepts future version")
		}
		if model.Version != "" {
			t.Errorf("Model has version: '%v'", model.Version)
		}
	}
}

func TestSystemName(t *testing.T) {
	yamlWithFutureVersion := `system:
  name: foo`

	model, _ := LintText(yamlWithFutureVersion)

	if model.System.Name != "foo" {
		t.Errorf("Model has system: '%v'", model.System.Name)
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

	if model.System.Name != "A system named foo" {
		t.Errorf("Model has system: '%v'", model.System.Name)
	}
}

func TestEmptySystemName(t *testing.T) {
	yamlWithFutureVersion := `system:
  name: ''`

	model, _ := LintText(yamlWithFutureVersion)

	if model.System.Name != "" {
		t.Errorf("System name should be empty, but is '%v'", model.Version)
	}
}

func TestNonStringSystem(t *testing.T) {
	yamlWithFutureVersion := `system:
  name: 3.14`

	model, _ := LintText(yamlWithFutureVersion)

	if model.System.Name != "3.14" {
		t.Errorf("Model has system: '%v'", model.System.Name)
	}
}

func TestInvalidSystemStructure(t *testing.T) {
	yamlWithFutureVersion := `system: foo`

	_, issues := LintText(yamlWithFutureVersion)

	if !hasIssue(issues, hasError("Expected a map")) {
		t.Errorf("No error on invalid system structure")
	}
}

func TestInvalidSystemNameStructure(t *testing.T) {
	yamlWithFutureVersion := `system:
  name:
    bar: baz`

	_, issues := LintText(yamlWithFutureVersion)

	if !hasIssue(issues, hasError("name must be a string")) {
		t.Errorf("No error on invalid system name structure")
	}
}

func TestPersonas(t *testing.T) {
	yamlWithPersonas := `personas:
  dev:
    name: Developer
    uses:
      - externalSystem: slack
        description: Reads notifications
      - form: subscriptions
        description: Maintains subscriber
  cs:
    name: Customer Support
    uses:
      - externalSystem: jira
        description: Updates issues
`

	model, _ := LintText(yamlWithPersonas)

	if len(model.Personas) != 2 {
		t.Fatalf("Incorrect number of personas: %v", len(model.Personas))
	}
	if model.Personas[0].Name != "Customer Support" {
		t.Errorf("Personas not sorted: incorrect name for 1st persona: %v", model.Personas[0].Name)
	}
}

func TestDefaultPersonaName(t *testing.T) {
	yamlWithPersonas := `personas:
  dev:
    uses:
      - externalSystem: slack
        description: Reads notifications
`

	model, _ := LintText(yamlWithPersonas)

	if len(model.Personas) != 1 {
		t.Fatalf("Incorrect number of personas: %v", len(model.Personas))
	}
	if model.Personas[0].Name != "Dev" {
		t.Errorf("Incorrect name: '%v'", model.Personas[0].Name)
	}
}

func TestInvalidPersona(t *testing.T) {
	yamlWithInvalidPersonas := `personas:
  dev:
    foo:
      bar: baz
`

	_, issues := LintText(yamlWithInvalidPersonas)

	if !hasIssue(issues, hasError("Invalid persona")) {
		t.Fatalf("No error on invalid personas")
	}
}

func TestExternalSystem(t *testing.T) {
	yamlWithExternalSystems := `externalSystems:
  broker:
    name: Privacy Broker
    type: central
  localPlatform:
    name: Local platform
    type: local
    calls:
      - service: api
        description: Sends SAR / DDR / CID
        dataFlow: send

`

	model, _ := LintText(yamlWithExternalSystems)

	if len(model.ExternalSystems) != 2 {
		t.Fatalf("Incorrect number of external systems: %v", len(model.ExternalSystems))
	}
	if model.ExternalSystems[0].Name != "Local platform" {
		t.Errorf("External systems not sorted: incorrect name for 1st external system: %v", model.ExternalSystems[0].Name)
	}
}
