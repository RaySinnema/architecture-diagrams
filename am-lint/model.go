package main

type ArchitectureModel struct {
	Version         string
	System          System
	Personas        []*Persona
	ExternalSystems []*ExternalSystem
}
