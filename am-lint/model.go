package main

import "gopkg.in/yaml.v3"

type ArchitectureModel struct {
	node              *yaml.Node
	Version           string
	System            System
	Personas          []*Persona
	ExternalSystems   []*ExternalSystem
	Services          []*Service
	Databases         []*DataStore
	Queues            []*DataStore
	Technologies      []*Technology
	TechnologyBundles []*TechnologyBundle
}
