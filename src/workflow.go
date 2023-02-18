package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"sort"
)

type Step struct {
	node             *yaml.Node
	Description      string
	SubWorkflowId    string
	PerformerId      string
	Performer        interface{}
	FormId           string
	Form             *Form
	Command          string
	Event            string
	ExternalSystemId string
	ExternalSystem   *ExternalSystem
	ServiceId        string
	Service          *Service
	View             string
}

func (s *Step) getDescription() string {
	return s.Description
}

func (s *Step) setDescription(description string) {
	s.Description = description
}

func (s *Step) read(node *yaml.Node) []Issue {
	s.node = node
	fields, issue := toMap(node)
	if issue != nil {
		return []Issue{*issue}
	}
	issues := make([]Issue, 0)
	issues = append(issues, setDescription(fields, s)...)
	subWorkflowId, found, issue := stringFieldOf(fields, "workflow")
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		s.SubWorkflowId = subWorkflowId
	} else {
		issues = append(issues, s.readPerformer(node, fields)...)
	}
	return issues
}

func (s *Step) readPerformer(node *yaml.Node, fields map[string]*yaml.Node) []Issue {
	issues := make([]Issue, 0)
	performerId, found, issue := stringFieldOf(fields, "performer")
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		s.PerformerId = performerId
		issues = append(issues, s.readPerformance(node, fields)...)
	} else {
		issues = append(issues, *NodeError("step needs a performer", node))
	}
	return issues
}

func (s *Step) readPerformance(node *yaml.Node, fields map[string]*yaml.Node) []Issue {
	issues := make([]Issue, 0)
	formId, formFound, issue := stringFieldOf(fields, "form")
	if issue != nil {
		issues = append(issues, *issue)
	}
	view, viewFound, issue := stringFieldOf(fields, "view")
	if issue != nil {
		issues = append(issues, *issue)
	}
	command, commandFound, issue := stringFieldOf(fields, "command")
	if issue != nil {
		issues = append(issues, *issue)
	}
	serviceId, serviceFound, issue := stringFieldOf(fields, "service")
	if issue != nil {
		issues = append(issues, *issue)
	}
	externalSystemId, externalSystemFound, issue := stringFieldOf(fields, "externalSystem")
	if issue != nil {
		issues = append(issues, *issue)
	}
	event, eventFound, issue := stringFieldOf(fields, "event")
	if issue != nil {
		issues = append(issues, *issue)
	}
	issue = needExactlyOneOf(node, []bool{commandFound, eventFound, externalSystemFound, formFound, serviceFound, viewFound})
	if issue != nil {
		return append(issues, *issue)
	}
	s.Command = command
	s.Event = event
	s.ExternalSystemId = externalSystemId
	s.FormId = formId
	s.ServiceId = serviceId
	s.View = view
	return issues
}

func needExactlyOneOf(node *yaml.Node, values []bool) *Issue {
	numFound := 0
	for _, value := range values {
		if value {
			numFound = numFound + 1
		}
	}
	switch numFound {
	case 0:
		return NodeError("Need one of command, event, externalSystem, form, service, or view", node)
	case 1:
		return nil
	default:
		return NodeError("Need exactly one of command, event, externalSystem, form, service, or view", node)
	}
}

type Workflow struct {
	node        *yaml.Node
	Id          string
	Name        string
	Description string
	Steps       []*Step
	TopLevel    bool
}

func (w *Workflow) Print(printer *Printer) {
	printer.Print(w.Name, " (")
	switch len(w.Steps) {
	case 0:
		printer.Print("no steps")
	case 1:
		printer.Print(" 1 step")
	default:
		printer.Print(len(w.Steps), " steps")
	}
	printer.PrintLn(")")
}

func (w *Workflow) setNode(node *yaml.Node) {
	w.node = node
}

func (w *Workflow) setId(id string) {
	w.Id = id
}

func (w *Workflow) setName(name string) {
	w.Name = name
}

func (w *Workflow) getDescription() string {
	return w.Description
}

func (w *Workflow) setDescription(description string) {
	w.Description = description
}

func (w *Workflow) read(id string, node *yaml.Node) []Issue {
	var fields map[string]*yaml.Node
	fields, issues := namedObject(node, id, w)
	issues = append(issues, setDescription(fields, w)...)
	issues = append(issues, w.readSteps(fields)...)
	return issues
}

func (w *Workflow) readSteps(fields map[string]*yaml.Node) []Issue {
	stepNodes, found, issue := sequenceFieldOf(fields, "steps")
	if issue != nil {
		return []Issue{*issue}
	}
	issues := make([]Issue, 0)
	steps := make([]*Step, 0)
	if found {
		for _, stepNode := range stepNodes {
			step := Step{}
			steps = append(steps, &step)
			issues = append(issues, step.read(stepNode)...)
		}
	}
	w.Steps = steps
	return issues
}

type WorkflowReader struct {
}

func (w WorkflowReader) read(node *yaml.Node, _ string, model *ArchitectureModel) []Issue {
	if node == nil {
		return []Issue{}
	}
	workflowsById, issue := toMap(node)
	if issue != nil {
		return []Issue{*issue}
	}
	issues := make([]Issue, 0)
	workflows := make([]*Workflow, 0)
	for id, workflowNode := range workflowsById {
		workflow := Workflow{TopLevel: true}
		workflows = append(workflows, &workflow)
		issues = append(issues, workflow.read(id, workflowNode)...)
	}
	sort.Slice(workflows, func(i, j int) bool {
		return workflows[i].Name < workflows[j].Name
	})
	model.Workflows = workflows
	return issues
}

type WorkflowCollector struct {
}

func (w WorkflowCollector) connect(model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	for _, workflow := range model.Workflows {
		for _, step := range workflow.Steps {
			issues = append(issues, w.connectStep(step, model)...)
		}
	}
	for _, workflow := range model.Workflows {
		issues = append(issues, w.connectSubWorkflows(workflow, model)...)
	}
	return issues
}

func (w WorkflowCollector) connectStep(step *Step, model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	issues = append(issues, w.connectTarget(step, model)...)
	issues = append(issues, w.connectPerformer(step, model)...)
	return issues
}

func (w WorkflowCollector) connectTarget(step *Step, model *ArchitectureModel) []Issue {
	if step.FormId != "" {
		form, issue := findForm(step.node, step.FormId, model)
		if issue != nil {
			return []Issue{*issue}
		}
		step.Form = form
		return []Issue{}
	}
	if step.ServiceId != "" {
		service, issue := findService(step.node, step.ServiceId, model)
		if issue != nil {
			return []Issue{*issue}
		}
		step.Service = service
		return []Issue{}
	}
	if step.ExternalSystemId != "" {
		externalSystem, issue := findExternalSystem(step.node, step.ExternalSystemId, model)
		if issue != nil {
			return []Issue{*issue}
		}
		step.ExternalSystem = externalSystem
		return []Issue{}
	}
	return []Issue{}
}

func findExternalSystem(node *yaml.Node, id string, model *ArchitectureModel) (*ExternalSystem, *Issue) {
	externalSystem, found := model.findExternalSystemById(id)
	if found {
		return externalSystem, nil
	}
	return nil, NodeError(fmt.Sprintf("Unknown external system '%v'", id), node)
}

func findService(node *yaml.Node, id string, model *ArchitectureModel) (*Service, *Issue) {
	service, found := model.findServiceById(id)
	if found {
		return service, nil
	}
	return nil, NodeError(fmt.Sprintf("Unknown service '%v'", id), node)
}

func findForm(node *yaml.Node, id string, model *ArchitectureModel) (*Form, *Issue) {
	for _, service := range model.Services {
		form, found := service.findFormById(id)
		if found {
			return form, nil
		}
	}
	return nil, NodeError(fmt.Sprintf("Unknown form '%v'", id), node)
}

func findPersona(node *yaml.Node, id string, model *ArchitectureModel) (*Persona, *Issue) {
	persona, found := model.findPersonaById(id)
	if found {
		return persona, nil
	}
	return nil, NodeError(fmt.Sprintf("Unknown form '%v'", id), node)
}

func (w WorkflowCollector) connectPerformer(step *Step, model *ArchitectureModel) []Issue {
	if step.FormId != "" {
		return w.connectFormPerformer(step, model)
	}
	if step.View != "" {
		return w.connectViewPerformer(step, model)
	}
	if step.Command != "" {
		return w.connectCommandPerformer(step, model)
	}
	if step.ServiceId != "" {
		return w.connectServicePerformer(step, model)
	}
	if step.ExternalSystemId != "" {
		return w.connectExternalSystemPerformer(step, model)
	}
	if step.Event != "" {
		return w.connectEventPerformer(step, model)
	}
	return []Issue{}
}

func (w WorkflowCollector) connectFormPerformer(step *Step, model *ArchitectureModel) []Issue {
	persona, issue := findPersona(step.node, step.PerformerId, model)
	if issue != nil {
		return []Issue{*issue}
	}
	step.Performer = persona
	return []Issue{}
}

func (w WorkflowCollector) connectViewPerformer(step *Step, model *ArchitectureModel) []Issue {
	service, found := model.findServiceById(step.PerformerId)
	if found {
		_, found = service.findDatabaseView(step.View)
		if found {
			step.Performer = service
		} else {
			return []Issue{*NodeError(fmt.Sprintf("Service '%v' doesn't have a database with view '%v'", step.PerformerId, step.View), step.node)}
		}
	} else {
		persona, found := model.findPersonaById(step.PerformerId)
		if found {
			step.Performer = persona
		} else {
			return []Issue{*NodeError(fmt.Sprintf("Unknown service or persona '%v'", step.PerformerId), step.node)}
		}
	}
	return []Issue{}
}

func (w WorkflowCollector) connectCommandPerformer(step *Step, model *ArchitectureModel) []Issue {
	form, issue := findForm(step.node, step.PerformerId, model)
	if issue == nil {
		step.Performer = form
		return []Issue{}
	}
	service, issue := findService(step.node, step.PerformerId, model)
	if issue == nil {
		step.Performer = service
		return []Issue{}
	}
	externalSystem, issue := findExternalSystem(step.node, step.PerformerId, model)
	if issue == nil {
		step.Performer = externalSystem
		return []Issue{}
	}
	return []Issue{*NodeError(fmt.Sprintf("Unknown form, service, or external system '%v'",
		step.PerformerId), step.node)}
}

func (w WorkflowCollector) connectServicePerformer(step *Step, model *ArchitectureModel) []Issue {
	service, issue := findService(step.node, step.PerformerId, model)
	if issue == nil {
		step.Performer = service
		return []Issue{}
	}
	externalSystem, issue := findExternalSystem(step.node, step.PerformerId, model)
	if issue == nil {
		step.Performer = externalSystem
		return []Issue{}
	}
	return []Issue{*NodeError(fmt.Sprintf("Unknown service or external system '%v'",
		step.PerformerId), step.node)}
}

func (w WorkflowCollector) connectExternalSystemPerformer(step *Step, model *ArchitectureModel) []Issue {
	persona, issue := findPersona(step.node, step.PerformerId, model)
	if issue == nil {
		step.Performer = persona
		return []Issue{}
	}
	service, issue := findService(step.node, step.PerformerId, model)
	if issue == nil {
		step.Performer = service
		return []Issue{}
	}
	return []Issue{*NodeError(fmt.Sprintf("Unknown persona or service '%v'",
		step.PerformerId), step.node)}
}

func (w WorkflowCollector) connectEventPerformer(step *Step, model *ArchitectureModel) []Issue {
	service, issue := findService(step.node, step.PerformerId, model)
	if issue == nil {
		step.Performer = service
		return []Issue{}
	}
	return []Issue{*NodeError(fmt.Sprintf("Unknown service '%v'", step.PerformerId), step.node)}
}

func (w WorkflowCollector) connectSubWorkflows(workflow *Workflow, model *ArchitectureModel) []Issue {
	steps := make([]*Step, 0)
	for _, step := range workflow.Steps {
		if step.SubWorkflowId == "" {
			steps = append(steps, step)
		} else {
			workflow, found := model.findWorkflowById(step.SubWorkflowId)
			if !found {
				return []Issue{*NodeError(fmt.Sprintf("Unknown workflow '%v'", step.SubWorkflowId), step.node)}
			}
			workflow.TopLevel = false
			w.connectSubWorkflows(workflow, model)
			steps = append(steps, workflow.Steps...)
		}
	}
	workflow.Steps = steps
	return []Issue{}
}
