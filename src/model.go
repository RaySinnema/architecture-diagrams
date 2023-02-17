package main

import (
	"gopkg.in/yaml.v3"
)

type ArchitectureModel struct {
	node              *yaml.Node
	Version           string
	System            System
	Personas          []*Persona
	ExternalSystems   []*ExternalSystem
	Services          []*Service
	Databases         []*Database
	Queues            []*DataStore
	Technologies      []*Technology
	TechnologyBundles []*TechnologyBundle
	Workflows         []*Workflow
}

func (model ArchitectureModel) String() string {
	printer := NewPrinter()
	model.writeSystem(printer)
	printer.Start()

	model.writePersonas(printer)
	model.writeExternalSystems(printer)
	model.writeServices(printer)
	model.writeDatabases(printer)
	model.writeQueues(printer)
	model.writeWorkflows(printer)

	printer.End()
	return printer.String()
}

func (model ArchitectureModel) writeSystem(printer *Printer) {
	printer.PrintLn(model.System.Name, ":")
}

func (model ArchitectureModel) writePersonas(printer *Printer) {
	if len(model.Personas) == 0 {
		return
	}
	printer.NewLine()
	printer.PrintLn("Personas:")
	for _, persona := range model.Personas {
		printer.Print("- ")
		persona.Print(printer)
	}
}

func (model ArchitectureModel) writeExternalSystems(printer *Printer) {
	if len(model.ExternalSystems) == 0 {
		return
	}
	printer.NewLine()
	printer.PrintLn("External systems:")
	for _, externalSystem := range model.ExternalSystems {
		printer.Print("- ")
		externalSystem.Print(printer)
	}
}

func (model ArchitectureModel) writeServices(printer *Printer) {
	if len(model.Services) == 0 {
		return
	}
	printer.NewLine()
	printer.PrintLn("Services:")
	for _, service := range model.Services {
		printer.Print("- ")
		service.Print(printer)
	}
}

func (model ArchitectureModel) writeDatabases(printer *Printer) {
	if len(model.Databases) == 0 {
		return
	}
	printer.NewLine()
	printer.PrintLn("Databases:")
	for _, database := range model.Databases {
		printer.Print("- ")
		database.Print(printer)
	}
}

func (model ArchitectureModel) writeQueues(printer *Printer) {
	if len(model.Queues) == 0 {
		return
	}
	printer.NewLine()
	printer.PrintLn("Queues:")
	for _, queue := range model.Queues {
		printer.Print("- ")
		queue.Print(printer)
	}
}

func (model ArchitectureModel) writeWorkflows(printer *Printer) {
	if len(model.Workflows) == 0 {
		return
	}
	printer.NewLine()
	printer.PrintLn("Workflows:")
	for _, workflow := range model.Workflows {
		if workflow.TopLevel {
			printer.Print("- ")
			workflow.Print(printer)
		}
	}
}

func (model ArchitectureModel) findPersonaById(id string) (*Persona, bool) {
	for _, candidate := range model.Personas {
		if candidate.Id == id {
			return candidate, true
		}
	}
	return nil, false
}

func (model ArchitectureModel) findExternalSystemById(id string) (*ExternalSystem, bool) {
	for _, candidate := range model.ExternalSystems {
		if candidate.Id == id {
			return candidate, true
		}
	}
	return nil, false
}

func (model ArchitectureModel) findServiceById(id string) (*Service, bool) {
	for _, candidate := range model.Services {
		if candidate.Id == id {
			return candidate, true
		}
	}
	return nil, false
}

func (model ArchitectureModel) findTechnologyBundleById(id string) (*TechnologyBundle, bool) {
	for _, candidate := range model.TechnologyBundles {
		if candidate.Id == id {
			return candidate, true
		}
	}
	return nil, false
}

func (model ArchitectureModel) findWorkflowById(id string) (*Workflow, bool) {
	for _, candidate := range model.Workflows {
		if candidate.Id == id {
			return candidate, true
		}
	}
	return nil, false
}

func (model ArchitectureModel) findDatabaseById(id string) (*Database, bool) {
	for _, candidate := range model.Databases {
		if candidate.Id == id {
			return candidate, true
		}
	}
	return nil, false
}
