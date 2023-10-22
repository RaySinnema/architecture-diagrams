package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"sort"
)

type Used struct {
	node             *yaml.Node
	Description      string
	ExternalSystemId string `yaml:"externalSystem,omitempty"`
	ExternalSystem   *ExternalSystem
	FormId           string `yaml:"form,omitempty"`
	Form             *Form
	ViewId           string
	View             *View
	DataFlow         DataFlow
}

func (u *Used) getDescription() string {
	return u.Description
}

func (u *Used) setDescription(description string) {
	u.Description = description
}

func (u *Used) setDataFlow(dataFlow DataFlow) {
	u.DataFlow = dataFlow
}

func (u *Used) read(node *yaml.Node) []Issue {
	u.node = node
	fields, issue := toMap(node)
	if issue != nil {
		return []Issue{*issue}
	}
	issues := make([]Issue, 0)
	issues = append(issues, u.readUsed(node, issue, fields)...)
	issues = append(issues, setDescription(fields, u)...)
	issues = append(issues, setDataFlow(node, fields, u)...)
	return issues
}

func (u *Used) readUsed(node *yaml.Node, issue *Issue, fields map[string]*yaml.Node) []Issue {
	issues := make([]Issue, 0)
	externalSystemId, foundExternalSystem, issue := stringFieldOf(fields, "externalSystem")
	if issue != nil {
		issues = append(issues, *issue)
	}
	formId, foundForm, issue := stringFieldOf(fields, "form")
	if issue != nil {
		issues = append(issues, *issue)
	}
	viewId, foundView, issue := stringFieldOf(fields, "view")
	if issue != nil {
		issues = append(issues, *issue)
	}
	if foundExternalSystem && foundForm {
		issues = append(issues, *NodeError("A persona may use either a form or an external system. Split the into two to have let the persona use both.", node))
	} else if foundExternalSystem {
		u.ExternalSystemId = externalSystemId
	} else if foundForm {
		u.FormId = formId
	} else if foundView {
		u.ViewId = viewId
	} else {
		issues = append(issues, *NodeError("A persona must use either a form, a view, or an external system", node))
	}
	return issues
}

func (u *Used) Used() interface{} {
	if u.ExternalSystem != nil {
		return u.ExternalSystem
	}
	return u.Form
}

type Persona struct {
	node        *yaml.Node
	Id          string
	Name        string
	Description string
	Uses        []*Used
}

func (p *Persona) Print(printer *Printer) {
	printer.PrintLn(p.Name)
}

func (p *Persona) setNode(node *yaml.Node) {
	p.node = node
}

func (p *Persona) setId(id string) {
	p.Id = id
}

func (p *Persona) setName(name string) {
	p.Name = name
}

func (p *Persona) getDescription() string {
	return p.Description
}

func (p *Persona) setDescription(description string) {
	p.Description = description
}

func (p *Persona) read(id string, node *yaml.Node) []Issue {
	var fields map[string]*yaml.Node
	fields, issues := namedObject(node, id, p)
	issues = append(issues, setDescription(fields, p)...)
	useNodes, found, issue := sequenceFieldOf(fields, "uses")
	if issue != nil {
		return append(issues, *issue)
	}
	if !found {
		return append(issues, *NodeError("A persona must use either a form, a view, or an external system", node))
	}
	uses := make([]*Used, 0)
	for _, useNode := range useNodes {
		use := Used{}
		issues = append(issues, use.read(useNode)...)
		uses = append(uses, &use)
	}
	p.Uses = uses
	return issues
}

type PersonaReader struct {
}

func (p PersonaReader) read(node *yaml.Node, _ string, model *ArchitectureModel) []Issue {
	if node == nil {
		return []Issue{}
	}
	personasById, issue := toMap(node)
	if issue != nil {
		return []Issue{*issue}
	}
	issues := make([]Issue, 0)
	personas := make([]*Persona, 0)
	for id, personaNode := range personasById {
		persona := Persona{}
		personas = append(personas, &persona)
		issues = append(issues, persona.read(id, personaNode)...)
	}
	sort.Slice(personas, func(i, j int) bool {
		return personas[i].Name < personas[j].Name
	})
	model.Personas = personas
	return issues
}

type PersonaConnector struct {
}

func (c PersonaConnector) connect(model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	for _, persona := range model.Personas {
		for _, use := range persona.Uses {
			issues = append(issues, c.connectUsed(use, model)...)
		}
	}
	return issues
}

func (c PersonaConnector) connectUsed(used *Used, model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	if used.ExternalSystemId != "" {
		for _, candidate := range model.ExternalSystems {
			if candidate.Id == used.ExternalSystemId {
				used.ExternalSystem = candidate
			}
		}
		if used.ExternalSystem == nil {
			issues = append(issues, *NodeError(fmt.Sprintf("Unknown external system '%v'", used.ExternalSystemId), used.node))
		}
	}
	if used.FormId != "" {
		for _, service := range model.Services {
			for _, candidate := range service.Forms {
				if candidate.Id == used.FormId {
					used.Form = candidate
				}
			}
		}
		if used.Form == nil {
			issues = append(issues, *NodeError(fmt.Sprintf("Unknown form '%v'", used.FormId), used.node))
		}
	}
	if used.ViewId != "" {
		for _, database := range model.Databases {
			for _, candidate := range database.Views {
				if candidate.Id == used.ViewId {
					used.View = candidate
				}
			}
		}
		if used.View == nil {
			issues = append(issues, *NodeError(fmt.Sprintf("Unknown view '%v'", used.ViewId), used.node))
		}
	}
	return issues
}

type PersonaValidator struct {
}

func (v PersonaValidator) validate(model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	if len(model.Personas) == 0 {
		issues = append(issues, *NodeWarning("At least one persona is required", model.node))
	}
	return issues
}
