package main

import (
	"gopkg.in/yaml.v3"
	"sort"
)

type Used struct {
	Description        string
	ExternalSystemName *string `yaml:"externalSystem,omitempty"`
	FormName           *string `yaml:"form,omitempty"`
}

type Persona struct {
	Name string
	Uses []Used
}

type PersonaReader struct {
}

func (p PersonaReader) read(node *yaml.Node, _ string, model *ArchitectureModel) []Issue {
	if node == nil {
		return []Issue{}
	}
	personasById := make(map[string]Persona)
	err := node.Decode(&personasById)
	if err != nil {
		return []Issue{*NodeError("Invalid personas", node)}
	}
	personas := make([]Persona, 0)
	for id, persona := range personasById {
		if persona.Name == "" {
			persona.Name = friendly(id)
		}
		personas = append(personas, persona)
	}
	sort.Slice(personas, func(i, j int) bool {
		return personas[i].Name < personas[j].Name
	})
	model.Personas = personas
	return []Issue{}
}
