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
}

func (u *Used) setDescription(description string) {
	u.Description = description
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
	if foundExternalSystem && foundForm {
		issues = append(issues, *NodeError("A use may have either a form or an external system. Split the use into two to have let the persona use both.", node))
	} else if foundExternalSystem {
		u.ExternalSystemId = externalSystemId
	} else if foundForm {
		u.FormId = formId
	} else {
		issues = append(issues, *NodeError("Must use either a form or an external system", node))
	}
	return issues
}

func (u *Used) Used() *ExternalSystem {
	if u.ExternalSystem != nil {
		return u.ExternalSystem
	}
	return nil
}

type Persona struct {
	node *yaml.Node
	Id   string
	Name string
	Uses []*Used
}

func (p *Persona) setNode(node *yaml.Node) {
	p.node = node
}

func (p *Persona) setName(name string) {
	p.Name = name
}

func (p *Persona) setId(id string) {
	p.Id = id
}

func (p *Persona) read(id string, node *yaml.Node) []Issue {
	var fields map[string]*yaml.Node
	fields, issues := namedObject(id, node, p)
	useNodes, found, issue := sequenceFieldOf(fields, "uses")
	if issue != nil {
		return append(issues, *issue)
	}
	if !found {
		return append(issues, *NodeError(fmt.Sprintf(
			"Invalid persona '%v': must use at least either one form or one external system", id), node))
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

type PersonaCollector struct {
}

func (c PersonaCollector) connect(model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	for _, persona := range model.Personas {
		for _, use := range persona.Uses {
			issues = append(issues, c.connectUsed(use, model)...)
		}
	}
	return issues
}

func (c PersonaCollector) connectUsed(used *Used, model *ArchitectureModel) []Issue {
	issues := make([]Issue, 0)
	if used.ExternalSystemId != "" {
		for _, candidate := range model.ExternalSystems {
			if candidate.Id == used.ExternalSystemId {
				used.ExternalSystem = candidate
			}
		}
		if used.ExternalSystem == nil {
			issues = append(issues, *NodeError(fmt.Sprintf("Unknown external system %v", used.ExternalSystemId), used.node))
		}
	}
	return issues
}
