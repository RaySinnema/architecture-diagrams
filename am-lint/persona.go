package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"sort"
)

type Used struct {
	Description        string
	ExternalSystemName string `yaml:"externalSystem,omitempty"`
	FormName           string `yaml:"form,omitempty"`
}

func (u *Used) read(node *yaml.Node, issues []Issue) []Issue {
	fields, issue := toMap(node)
	if issue != nil {
		return append(issues, *issue)
	}
	externalSystemName, foundExternalSystem, issue := stringFieldOf(fields, "externalSystem")
	if issue != nil {
		issues = append(issues, *issue)
	}
	formName, foundForm, issue := stringFieldOf(fields, "form")
	if issue != nil {
		issues = append(issues, *issue)
	}
	if foundExternalSystem && foundForm {
		issues = append(issues, *NodeError("A use may have either a form or an external system. Split the use into two to have let the persona use both.", node))
	} else if foundExternalSystem {
		u.ExternalSystemName = externalSystemName
	} else if foundForm {
		u.FormName = formName
	} else {
		issues = append(issues, *NodeError("Must use either a form or an external system", node))
	}
	description, found, issue := stringFieldOf(fields, "description")
	if issue != nil {
		issues = append(issues, *issue)
	} else if found {
		u.Description = description
	}
	return issues
}

type Persona struct {
	Name string
	Uses []Used
}

func (p *Persona) setName(name string) {
	p.Name = name
}

func (p *Persona) read(id string, node *yaml.Node, issues []Issue) []Issue {
	var fields map[string]*yaml.Node
	fields, issues = namedObject(id, node, p, issues)
	useNodes, found, issue := sequenceFieldOf(fields, "uses")
	if issue != nil {
		return append(issues, *issue)
	}
	if !found {
		return append(issues, *NodeError(fmt.Sprintf(
			"Invalid persona '%v': must use at least either one form or one external system", id), node))
	}
	uses := make([]Used, 0)
	for _, useNode := range useNodes {
		use := Used{}
		issues = use.read(useNode, issues)
		uses = append(uses, use)
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
	personas := make([]Persona, 0)
	for id, personaNode := range personasById {
		persona := Persona{}
		issues = persona.read(id, personaNode, issues)
		personas = append(personas, persona)
	}
	sort.Slice(personas, func(i, j int) bool {
		return personas[i].Name < personas[j].Name
	})
	model.Personas = personas
	return issues
}
