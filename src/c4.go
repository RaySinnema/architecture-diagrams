package main

import (
	"os"
)

const idOfSystemOfInterest = "system"

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
	printPersons(model, printer)
	printSoftwareSystems(model, printer)
	printer.End()
	printer.PrintLn("}")
}

func printPersons(model *ArchitectureModel, printer *Printer) {
	for _, persona := range model.Personas {
		printer.PrintLn(persona.Id, " = person \"", persona.Name, "\" {")
		printer.Start()
		printDescription(persona, printer)
		printedSystemOfInterest := false
		for _, used := range persona.Uses {
			if used.ExternalSystemId != "" {
				printUses(used.ExternalSystemId, used.Description, printer)
			} else if !printedSystemOfInterest {
				printedSystemOfInterest = true
				printUses(idOfSystemOfInterest, used.Description, printer)
			}
		}
		printer.End()
		printer.PrintLn("}")
	}
}

func printDescription(describable Describable, printer *Printer) {
	if describable.getDescription() != "" {
		printer.PrintLn("description \"", describable.getDescription(), "\"")
	}
}

func printUses(id string, description string, printer *Printer) {
	printer.Print("-> ", id)
	if description != "" {
		printer.Print(" \"", description, "\"")
	}
	printer.NewLine()
}

func printSoftwareSystems(model *ArchitectureModel, printer *Printer) {
	printSoftwareSystemOfInterest(model, printer)
	printExternalSystems(model, printer)
}

func printSoftwareSystemOfInterest(model *ArchitectureModel, printer *Printer) {
	printer.PrintLn(idOfSystemOfInterest, " = softwareSystem \"", model.System.Name, "\" {")
	printer.Start()
	printContainers(model, printer)
	printer.End()
	printer.PrintLn("}")
}

func printContainers(model *ArchitectureModel, printer *Printer) {
	printServices(model.Services, printer)
	printDataStores(model.Databases, printer)
	printDataStores(model.Queues, printer)
}

func printServices(services []*Service, printer *Printer) {
	for _, service := range services {
		printer.PrintLn(service.Id, " container \"", service.Name, "\" {")
		printer.Start()
		printDescription(service, printer)
		printTechnology(service, printer)
		for _, call := range service.Calls {
			if call.ExternalSystemId != "" {
				printUses(call.ExternalSystemId, call.Description, printer)
			} else {
				printUses(call.ServiceId, call.Description, printer)
			}
		}
		for _, dataStore := range service.DataStores {
			id := dataStore.QueueId
			if id == "" {
				id = dataStore.DatabaseId
			}
			printUses(id, dataStore.Description, printer)
		}
		printer.End()
		printer.PrintLn("}")
	}
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

func printExternalSystems(model *ArchitectureModel, printer *Printer) {
	for _, externalSystem := range model.ExternalSystems {
		printer.PrintLn(externalSystem.Id, " = softwareSystem \"", externalSystem.Name, "\" {")
		printer.Start()
		printDescription(externalSystem, printer)
		printedSystemOfInterest := false
		for _, call := range externalSystem.Calls {
			if call.ExternalSystemId != "" {
				printUses(call.ExternalSystemId, call.Description, printer)
			} else if !printedSystemOfInterest {
				printedSystemOfInterest = true
				printUses(idOfSystemOfInterest, call.Description, printer)
			}
		}
		printer.End()
		printer.PrintLn("}")
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
