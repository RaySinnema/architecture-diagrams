package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

type InvalidDefinition struct {
	definition string
	error      string
}

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
	definitions := []string{
		`version: 3.14`,
		`version: 1.2`,
		`version: 1.0.1`,
	}

	for _, yamlWithFutureVersion := range definitions {
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

func TestNonStringSystemName(t *testing.T) {
	yamlWithFutureVersion := `system:
  name: 3.14`

	model, _ := LintText(yamlWithFutureVersion)

	if model.System.Name != "3.14" {
		t.Errorf("Model has system: '%v'", model.System.Name)
	}
}

func TestInvalidSystems(t *testing.T) {
	assertErrorsForInvalidDefinitions(t, []InvalidDefinition{
		{definition: `system: foo`, error: "Expected a map"},
		{definition: `system:
  name:
    bar: baz`, error: "name must be a string"},
	})
}

func assertErrorsForInvalidDefinitions(t *testing.T, cases []InvalidDefinition) {
	for _, c := range cases {
		_, issues := LintText(c.definition)

		if !hasIssue(issues, hasError(c.error)) {
			t.Errorf("Missing error '%v' for invalid definition '%v'\n\nInstead, got: %+v", c.error, c.definition, issues)
		}
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

func TestInvalidPersonas(t *testing.T) {
	assertErrorsForInvalidDefinitions(t, []InvalidDefinition{
		{definition: `personas:
  - foo
  - bar
`, error: "Expected a map"},
		{definition: `personas:
  foo:
    uses: bar
`, error: "uses must be a sequence"},
		{definition: `personas:
  foo:
    uses:
      - bar
`, error: "Expected a map"},
		{definition: `personas:
  dev:
    foo:
      bar: baz
`, error: "Invalid persona"},
		{definition: `personas:
  foo:
    uses:
      - externalSystem:
        - bar
`, error: "externalSystem must be a string, not a sequence"},
		{definition: `personas:
  foo:
    uses:
      - form:
        - bar
`, error: "form must be a string, not a sequence"},
		{definition: `personas:
  foo:
    uses:
      - form: bar
        externalSystem: baz
`, error: "A use may have either a form or an external system"},
		{definition: `personas:
  foo:
    uses:
      - description: bar
`, error: "Must use either a form or an external system"},
		{definition: `personas:
  foo:
    uses:
      - description:
          - bar
`, error: "description must be a string, not a sequence"},
	})
}

func TestPersonaUsesExternalSystem(t *testing.T) {
	definition := `personas:
  ape:
    uses:
      - externalSystem: bear

externalSystems:
  bear:
    description: foo
`

	model, issues := LintText(definition)

	if len(issues) > 0 {
		t.Errorf("Got issues: %+v", issues)
	}
	if model.Personas[0].Uses[0].Used() != model.ExternalSystems[0] {
		t.Errorf("Persona use isn't linked up: %+v", *model)
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
		t.Fatalf("External systems not sorted: incorrect name for 1st external system: %v", model.ExternalSystems[0].Name)
	}
	lp := model.ExternalSystems[0]
	if len(lp.Calls) != 1 {
		t.Fatalf("# calls: %v", len(lp.Calls))
	}
	call := lp.Calls[0]
	if call.ServiceId != "api" {
		t.Errorf("Invalid call: '%v'", call.ServiceId)
	}
	if call.DataFlow != Send {
		t.Errorf("Invalid call direction: %v", call.DataFlow)
	}
}

func TestInvalidExternalSystems(t *testing.T) {
	assertErrorsForInvalidDefinitions(t, []InvalidDefinition{
		{definition: `externalSystems:
  - foo
  - bar
`, error: "Expected a map"},
		{definition: `externalSystems:
  foo:
    type:
      bar: baz
`, error: "type must be a string, not a map"},
		{definition: `externalSystems:
  foo:
    calls:
      bar: baz
`, error: "calls must be a sequence, not a map"},
		{definition: `externalSystems:
  foo:
    calls:
      - bar
`, error: "Expected a map"},
		{definition: `externalSystems:
  foo:
    calls:
      - service:
          - bar
`, error: "service must be a string"},
		{definition: `externalSystems:
  foo:
    calls:
      - externalSystem:
          - bar
`, error: "externalSystem must be a string"},
		{definition: `externalSystems:
  foo:
    calls:
      - externalSystem: bar
        service: baz
`, error: "A call may be to either a service or to an externalSystem"},
		{definition: `externalSystems:
  foo:
    calls:
      - description: bar
`, error: "One of service or externalSystem is required"},
		{definition: `externalSystems:
  foo:
    calls:
      - description:
          - bar
`, error: "description must be a string"},
		{definition: `externalSystems:
  foo:
    calls:
      - dataFlow:
          - bar
`, error: "dataFlow must be a string"},
		{definition: `externalSystems:
  foo:
    calls:
      - dataFlow: bar
`, error: "Invalid dataFlow: must be one of 'send', 'receive', or 'bidirectional'"},
	})
}

func TestExternalSystemCallsExternalSystem(t *testing.T) {
	definition := `externalSystems:
  ape:
    calls:
      - externalSystem: bear
  bear:
    description: foo
`

	model, _ := LintText(definition)

	call := model.ExternalSystems[0].Calls[0]
	if call.Callee() != model.ExternalSystems[1] {
		t.Errorf("Callee isn't linked up: %+v", *model)
	}
	if call.DataFlow != Bidirectional {
		t.Errorf("Invalid call direction: %v", call.DataFlow)
	}
}

func TestService(t *testing.T) {
	definition := `services:
  form:
    name: Privacy Form
    state: emerging
    technologies:
      - spring
      - java
    forms:
      - privacy
    calls:
      - service: api
        description: Calls
        dataFlow: send
        technologies: jsonOverHttp
  api:
    name: API
    dataStores:
      - queue: events
        description: Writes domain events
        dataFlow: send
    technologies: server
  subscriptions:
    dataStores:
      - queue: events
        dataFlow: receive
      - database: subscriptions
    forms:
      subscriptionsOld:
        name: Subscriptions
        state: legacy
`

	model, issues := LintText(definition)

	if model == nil {
		t.Fatalf("Invalid model: %+v", issues)
	}
	if len(issues) > 0 {
		t.Errorf("Unexpected issues: %+v", issues)
	}
	if len(model.Services) != 3 {
		t.Fatalf("Incorrect number of services: %v", len(model.Services))
	}

	api := model.Services[0]
	if api.Name != "API" {
		t.Fatalf("Services not sorted: incorrect name for 1st service: %v", api.Name)
	}
	if len(api.DataStores) == 1 {
		if api.DataStores[0].QueueId != "events" {
			t.Errorf("Invalid data store: %+v", api.DataStores[0])
		}
	} else {
		t.Errorf("Invalid # data stores: %+v", api)
	}
	if api.TechnologiesId != "server" {
		t.Errorf("Invalid technologies: %+v", api.TechnologiesId)
	}
	if api.State != Ok {
		t.Errorf("Invalid api state: '%v'", api.State)
	}

	form := model.Services[1]
	if len(form.Forms) != 1 {
		t.Fatalf("Invalid # forms: %+v", form)
	}
	if !reflect.DeepEqual(form.TechnologyIds, []string{"java", "spring"}) {
		t.Errorf("Invalid form technologies: %+v", form.TechnologyIds)
	}
	if form.State != Emerging {
		t.Errorf("Invalid form state: '%v'", form.State)
	}

	subscriptions := model.Services[2]
	if subscriptions.Name != "Subscriptions" {
		t.Fatalf("Invalid subscriptions service: '%v'", subscriptions.Name)
	}
	if subscriptions.DataStores[0].DataFlow != Receive {
		t.Errorf("Invalid subscriptions queue dataflow: %v", subscriptions.DataStores[0].DataFlow)
	}
	if len(subscriptions.Forms) != 1 {
		t.Fatalf("Invalid # subscriptions forms: %+v", subscriptions)
	}
}

func TestInvalidService(t *testing.T) {
	assertErrorsForInvalidDefinitions(t, []InvalidDefinition{
		{definition: `services:
  - foo
  - bar
`, error: "Expected a map"},
		{definition: `services:
  foo:
    dataStores: bar
`, error: "dataStores must be a sequence"},
		{definition: `services:
  foo:
    dataStores:
      - bar
`, error: "Expected a map"},
		{definition: `services:
  foo:
    dataStores:
      - queue:
          - bar
`, error: "queue must be a string"},
		{definition: `services:
  foo:
    dataStores:
      - database:
          - bar
`, error: "database must be a string"},
		{definition: `services:
  foo:
    dataStores:
      - description: bar
`, error: "A dataStore must have either a database or a queue"},
		{definition: `services:
  foo:
    dataStores:
      - database: bar
        description:
          - baz
`, error: "description must be a string"},
		{definition: `services:
  foo:
    dataStores:
      - database: bar
        dataFlow: baz
`, error: "Invalid dataFlow: must be one of 'send', 'receive', or 'bidirectional'"},
		{definition: `services:
  foo:
    technologies:
      bar: baz
`, error: "technologies must be a sequence"},
		{definition: `services:
  foo:
    technologies:
      - bar:
          - baz
`, error: "technology must be a string"},
		{definition: `services:
  foo:
    forms: bar
`, error: "forms must be a sequence"},
		{definition: `services:
  foo:
    forms:
      - bar:
        - baz
`, error: "form must be a string"},
		{definition: `services:
  foo:
    forms:
      bar:
        state:
          - baz
`, error: "state must be a string"},
		{definition: `services:
  foo:
    forms:
      bar:
        name:
          - baz
`, error: "name must be a string"},
		{definition: `services:
  foo:
    state:
      - bar
`, error: "state must be a string"},
		{definition: `services:
  foo:
    state: weird
`, error: "Invalid state: must be one of 'ok', 'emerging', 'review', 'revision', 'legacy', or 'deprecated'"},
	})
}
