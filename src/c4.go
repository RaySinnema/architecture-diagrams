package main

import (
	"os"
)

const idOfSystemOfInterest = "system"

type usage struct {
	user        string
	used        string
	description string
}

func GenerateC4(model *ArchitectureModel, fileName string) error {
	printer := NewPrinter()
	printC4workspace(model, printer)
	data := []byte(printer.String())
	return os.WriteFile(fileName, data, 0666)
}

func printC4workspace(model *ArchitectureModel, printer *Printer) {
	printer.PrintLn("workspace {")
	printer.Start()
	printC4model(model, printer)
	printC4views(printer)
	printer.End()
	printer.PrintLn("}")
}

func printC4model(model *ArchitectureModel, printer *Printer) {
	printer.PrintLn("model {")
	printer.Start()

	usages := make([]usage, 0)
	usages = append(usages, printPersons(model, printer)...)
	usages = append(usages, printSoftwareSystems(model, printer)...)
	printUsages(usages, printer)

	printer.End()
	printer.PrintLn("}")
}

func printPersons(model *ArchitectureModel, printer *Printer) []usage {
	usages := make([]usage, 0)
	for _, persona := range model.Personas {
		printer.PrintLn(persona.Id, " = person \"", persona.Name, "\" {")
		printer.Start()
		printDescription(persona, printer)
		printedSystemOfInterest := false
		for _, used := range persona.Uses {
			if used.ExternalSystemId != "" {
				usages = append(usages, usage{persona.Id, used.ExternalSystemId, used.Description})
			} else if !printedSystemOfInterest {
				printedSystemOfInterest = true
				usages = append(usages, usage{persona.Id, idOfSystemOfInterest, used.Description})
			}
		}
		printer.End()
		printer.PrintLn("}")
	}
	return usages
}

func printDescription(describable Describable, printer *Printer) {
	if describable.getDescription() != "" {
		printer.PrintLn("description \"", describable.getDescription(), "\"")
	}
}

func printSoftwareSystems(model *ArchitectureModel, printer *Printer) []usage {
	var usages = make([]usage, 0)
	usages = append(usages, printSoftwareSystemOfInterest(model, printer)...)
	usages = append(usages, printExternalSystems(model, printer)...)
	return usages
}

func printSoftwareSystemOfInterest(model *ArchitectureModel, printer *Printer) []usage {
	printer.PrintLn(idOfSystemOfInterest, " = softwareSystem \"", model.System.Name, "\" {")
	printer.Start()
	usages := printContainers(model, printer)
	printer.End()
	printer.PrintLn("}")
	return usages
}

func printContainers(model *ArchitectureModel, printer *Printer) []usage {
	usages := make([]usage, 0)
	usages = append(usages, printServices(model.Services, printer)...)
	printDataStores(model.Databases, printer)
	printDataStores(model.Queues, printer)
	return usages
}

func printServices(services []*Service, printer *Printer) []usage {
	usages := make([]usage, 0)
	for _, service := range services {
		printer.PrintLn(service.Id, " = container \"", service.Name, "\" {")
		printer.Start()
		printDescription(service, printer)
		printTechnology(service, printer)
		for _, call := range service.Calls {
			if call.ExternalSystemId != "" {
				usages = append(usages, usage{service.Id, call.ExternalSystemId, call.Description})
			} else {
				usages = append(usages, usage{service.Id, call.ServiceId, call.Description})
			}
		}
		for _, dataStore := range service.DataStores {
			id := dataStore.QueueId
			if id == "" {
				id = dataStore.DatabaseId
			}
			usages = append(usages, usage{service.Id, id, dataStore.Description})
		}
		printer.End()
		printer.PrintLn("}")
	}
	return usages
}

func printTechnology(implementable Implementable, printer *Printer) {
	technologies := implementable.getTechnologies()
	if len(technologies) == 0 {
		return
	}
	printer.Print("technology \"")
	prefix := ""
	for _, technology := range technologies {
		printer.Print(prefix, technology)
		prefix = ", "
	}
	printer.PrintLn("\"")
}

func printDataStores(dataStores []*DataStore, printer *Printer) {
	for _, dataStore := range dataStores {
		printer.PrintLn(dataStore.Id, " = container \"", dataStore.Name, "\" {")
		printer.Start()
		printDescription(dataStore, printer)
		printTechnology(dataStore, printer)
		printer.End()
		printer.PrintLn("}")
	}
}

func printExternalSystems(model *ArchitectureModel, printer *Printer) []usage {
	usages := make([]usage, 0)
	for _, externalSystem := range model.ExternalSystems {
		printer.PrintLn(externalSystem.Id, " = softwareSystem \"", externalSystem.Name, "\" {")
		printer.Start()
		printDescription(externalSystem, printer)
		printedSystemOfInterest := false
		for _, call := range externalSystem.Calls {
			if call.ExternalSystemId != "" {
				usages = append(usages, usage{externalSystem.Id, call.ExternalSystemId, call.Description})
			} else if !printedSystemOfInterest {
				printedSystemOfInterest = true
				usages = append(usages, usage{externalSystem.Id, idOfSystemOfInterest, call.Description})
			}
		}
		printer.End()
		printer.PrintLn("}")
	}
	return usages
}

func printUsages(usages []usage, printer *Printer) {
	for _, usage := range usages {
		printer.Print(usage.user, " -> ", usage.used)
		if usage.description != "" {
			printer.Print(" \"", usage.description, "\"")
		}
		printer.NewLine()
	}
}

func printC4views(printer *Printer) {
	printer.PrintLn("views {")
	printer.Start()
	printer.PrintLn("systemContext ", idOfSystemOfInterest, " {")
	printer.PrintLn("}")
	printer.PrintLn("container ", idOfSystemOfInterest, " {")
	printer.PrintLn("}")
	printer.End()
	printer.PrintLn("}")
}
