package main

type ArchitectureModel struct {
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
