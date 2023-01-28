package main

import (
	"os"
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
	definition := `personas:
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

externalSystems:
  slack:
    name: Slack

services:
  console:
    forms:
      - subscriptions
`

	model, _ := LintText(definition)

	if len(model.Personas) != 2 {
		t.Fatalf("Incorrect number of personas: %v", len(model.Personas))
	}
	dev := model.Personas[1]
	if dev.Name != "Developer" {
		t.Fatalf("Personas not sorted: incorrect name for 2nd persona: %v", model.Personas[0].Name)
	}
	if dev.Uses[0].ExternalSystem != model.ExternalSystems[0] {
		t.Errorf("Missing use of external system: %+v", dev.Uses[0])
	}
	if dev.Uses[1].Used() != model.Services[0].Forms[0] {
		t.Errorf("Missing use of service: %+v", dev.Uses[1])
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
`, error: "A persona must use either a form or an external system"},
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
`, error: "Unknown form 'bar'"},
		{definition: `personas:
  foo:
    uses:
      - form: bar
        externalSystem: baz
`, error: "A persona may use either a form or an external system"},
		{definition: `personas:
  foo:
    uses:
      - description: bar
`, error: "A persona must use either a form or an external system"},
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
	used := model.Personas[0].Uses[0]
	if used.Used() != model.ExternalSystems[0] {
		t.Errorf("Persona use isn't linked up: %+v", *model)
	}
	if used.DataFlow != Bidirectional {
		t.Errorf("Persona use has wrong direction of data flow: %+v", *model)
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
        technologies:
          - http

services:
  api:
    name: API

technologies:
  http:
    name: HTTP
`

	model, issues := LintText(yamlWithExternalSystems)

	if len(model.ExternalSystems) != 2 {
		t.Fatalf("Incorrect number of external systems: %+v", issues)
	}
	if model.ExternalSystems[0].Name != "Local platform" {
		t.Fatalf("External systems not sorted: incorrect name for 1st external system: %v", model.ExternalSystems[0].Name)
	}
	lp := model.ExternalSystems[0]
	if lp.Type != "local" {
		t.Errorf("Invalid type: %+v", lp)
	}
	if len(lp.Calls) != 1 {
		t.Fatalf("# calls: %v", len(lp.Calls))
	}
	call := lp.Calls[0]
	if call.Callee() != model.Services[0] {
		t.Errorf("Invalid call: '%+v'", call)
	}
	if call.DataFlow != Send {
		t.Errorf("Invalid call direction: %v", call.DataFlow)
	}
	if len(call.Technologies) != 1 {
		t.Errorf("Technologies not resolved")
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
      - service: bar
`, error: "Unknown service 'bar'"},
		{definition: `externalSystems:
  foo:
    calls:
      - externalSystem: bar
`, error: "Unknown external system 'bar'"},
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

technologyBundles:
  jsonOverHttp:
    - json
    - http
  server:
    - java
    - spring
    - docker

technologies:
  json:
    name: JSON
  http:
    name: HTTP 1.1
  java:
    name: Java 17
  spring:
    name: Spring Boot 2
  docker:
    name: Docker
`

	model, issues := LintText(definition)

	if model == nil {
		t.Fatalf("Invalid model: %+v", issues)
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
	if len(api.Technologies) != 3 {
		t.Errorf("Service technologies not resolved: %+v", api)
	}
	if api.State != Ok {
		t.Errorf("Invalid api state: '%v'", api.State)
	}

	form := model.Services[1]
	if len(form.Forms) != 1 {
		t.Fatalf("Invalid # forms: %+v", form)
	}
	if len(form.Calls[0].Technologies) != 2 {
		t.Errorf("Invalid form technologies: %+v", form)
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
`, error: "A dataStore must be a database or a queue"},
		{definition: `services:
  foo:
    dataStores:
      - database: bar
        queue: baz
`, error: "A dataStore can be either a database or a queue, but not both"},
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
    technologies:
      - bar
`, error: "Unknown technology 'bar'"},
		{definition: `services:
  foo:
    calls: bar
`, error: "calls must be a sequence"},
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

func TestDatabase(t *testing.T) {
	definition := `databases:
  subscriptions:
    name: Subscriptions & in-flight requests
    technologies: cloudMySql
    apiTechnologies: sql
  requestIdsMap:
    name: Request IDs map
    description: Stores mapping between IDs from different systems for the same request
    technologies:
      - cloudsql
      - mysql
    apiTechnologies:
       - sql

technologyBundles:
  cloudMySql:
    - cloudsql
    - mysql

technologies:
  cloudsql:
    name: Google CloudSQL
  mysql:
    name: MySQL
  sql:
    name: SQL
`

	model, _ := LintText(definition)

	if len(model.Databases) != 2 {
		t.Fatalf("Invalid # databases: %+v", model.Databases)
	}
	if model.Databases[0].Name != "Request IDs map" {
		t.Fatalf("Databases not sorted")
	}
	if len(model.Databases[0].Technologies) != 2 {
		t.Errorf("Didn't resolve technologies")
	}
}

func TestInvalidDatabase(t *testing.T) {
	assertErrorsForInvalidDefinitions(t, []InvalidDefinition{
		{definition: `databases:
  - foo
  - bar
`, error: "Expected a map"},
	})
}

func TestQueue(t *testing.T) {
	definition := `queues:
  ape:
    name: Zebra
  events:
    name: Domain Events
    description: One topic per event type.
    technologies:
      - pubsub
    apiTechnologies: grpc

technologies:
  pubsub:
    name: Google Pub/Sub
  grpc:
    name: gRPC
`

	model, _ := LintText(definition)

	if len(model.Queues) != 2 {
		t.Fatalf("Invalid # queues: %+v", model.Queues)
	}
	if model.Queues[0].Name != "Domain Events" {
		t.Fatalf("Queues not sorted")
	}
	if len(model.Queues[0].Technologies) != 1 {
		t.Fatalf("Technologies not resolved")
	}
	if len(model.Queues[0].ApiTechnologies) != 1 {
		t.Fatalf("API technologies not resolved")
	}
}

func TestTechnologies(t *testing.T) {
	definition := `technologies:
  cloudSql:
    name: CloudSQL
    quadrant: platforms
    ring: hold
    description: Replace with AWS technology
  adr:
    name: Architecture Decision Records
    quadrant: techniques
`

	model, _ := LintText(definition)

	if len(model.Technologies) != 2 {
		t.Fatalf("Invalid # technologies: %+v", model.Queues)
	}
	if model.Technologies[0].Name != "Architecture Decision Records" {
		t.Fatalf("Queues not technologies")
	}
	adr := model.Technologies[0]
	if adr.Quadrant != Techniques {
		t.Errorf("Invalid quadrant: %v", adr.Quadrant)
	}
	if adr.Ring != Adopt {
		t.Errorf("Invalid ring: %v", adr.Ring)
	}
}

func TestInvalidTechnology(t *testing.T) {
	assertErrorsForInvalidDefinitions(t, []InvalidDefinition{
		{definition: `technologies:
  foo:
    ring: adopt
`, error: "Missing required field quadrant"},
		{definition: `technologies:
  foo:
    quadrant: tools
    ring: bar
`, error: "Invalid ring: must be one of 'trial', 'assess', 'adopt', or 'hold'"},
	})
}

func TestTechnologyBundle(t *testing.T) {
	definition := `technologyBundles:
  serverWithUi:
    - server
    - ui
    - thymeleaf
  ui:
    - thymeleaf
  server:
    - java
    - spring
    - docker

technologies:
  docker:
    name: Docker
    quadrant: platform
    ring: adopt
  java:
    name: Java 17
    quadrant: languagesAndFrameworks
    ring: adopt
  spring:
    name: Spring Boot 2
    quadrant: languagesAndFrameworks
    ring: hold
    description: Upgrade to version 3
  thymeleaf:
    name: Thymeleaf
    quadrant: languagesAndFrameworks
    ring: adopt
`

	model, _ := LintText(definition)

	if len(model.TechnologyBundles) != 3 {
		t.Fatalf("Invalid # technologies: %+v", model.TechnologyBundles)
	}
	if model.TechnologyBundles[0].Id != "server" {
		t.Fatalf("Technology bundles not sorted: %+v", model.TechnologyBundles)
	}
	server := model.TechnologyBundles[0]
	if len(server.TechnologyIds) != 3 {
		t.Errorf("Invalid # technologies: %v", server.TechnologyIds)
	}
	if len(server.Technologies) != 3 {
		t.Errorf("Didn't resolve all technologies for bundle")
	}
	if len(model.TechnologyBundles[1].Technologies) != 4 {
		t.Errorf("Didn't recursively resolve all technologies for bundle: %v", model.TechnologyBundles[1].Technologies)
	}
}

func TestInvalidTechnologyBundle(t *testing.T) {
	assertErrorsForInvalidDefinitions(t, []InvalidDefinition{
		{definition: `technologyBundles:
  foo:
    bar: baz
`, error: "technologies must be a sequence"},
		{definition: `technologyBundles:
  foo:
    - bar: baz
`, error: "technology must be a string"},
		{definition: `technologyBundles:
  foo:
    - bar
`, error: "Unknown technology 'bar'"},
	})
}

func TestModelValidation(t *testing.T) {
	assertWarningsForInvalidDefinitions(t, []InvalidDefinition{
		{definition: ``, error: "At least one persona is required"},
	})
}

func assertWarningsForInvalidDefinitions(t *testing.T, cases []InvalidDefinition) {
	for _, c := range cases {
		_, issues := LintText(c.definition)

		if !hasIssue(issues, hasWarning(c.error)) {
			t.Errorf("Missing warning '%v' for invalid definition '%v'\n\nInstead, got: %+v", c.error, c.definition, issues)
		}
	}
}
