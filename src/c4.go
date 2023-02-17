package main

import (
	"os"
)

const idOfSystemOfInterest = "system"

type usage struct {
	user          string
	used          string
	description   string
	bidirectional bool
	byPersona     bool
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
		for _, used := range persona.Uses {
			if used.ExternalSystem != nil {
				usages = append(usages, usage{persona.Id, used.ExternalSystem.Id, used.Description, false, true})
			} else if used.Form != nil {
				usages = append(usages, usage{persona.Id, used.Form.ImplementedBy.Id, used.Description, false, true})
			} else if used.View != nil {
				usages = append(usages, usage{persona.Id, used.View.ImplementedBy.Id + "_db", used.Description, false, true})
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
	printer.PrintLn("tags \"System of Interest\"")
	usages := printContainers(model, printer)
	printer.End()
	printer.PrintLn("}")
	return usages
}

func printContainers(model *ArchitectureModel, printer *Printer) []usage {
	usages := make([]usage, 0)
	usages = append(usages, printServices(model.Services, printer)...)
	printDatabases(model.Databases, printer)
	printQueues(model.Queues, printer)
	return usages
}

func printServices(services []*Service, printer *Printer) []usage {
	usages := make([]usage, 0)
	for _, service := range services {
		printer.PrintLn(service.Id, " = container \"", service.Name, "\" {")
		printer.Start()
		printDescription(service, printer)
		printTechnology(service, printer)
		printer.PrintLn("tags \"Service\" \"", service.State.String(), "\"")
		for _, call := range service.Calls {
			if call.ExternalSystemId != "" {
				usages = append(usages, usage{service.Id, call.ExternalSystemId, call.Description,
					call.DataFlow == Bidirectional, false})
			} else {
				usages = append(usages, usage{service.Id, call.ServiceId, call.Description,
					call.DataFlow == Bidirectional, false})
			}
		}
		for _, dataStore := range service.DataStores {
			id := dataStore.QueueId
			if id == "" {
				id = dataStore.DatabaseId + "_db"
			} else {
				id = id + "_q"
			}
			usages = append(usages, usage{service.Id, id, dataStore.Description,
				dataStore.DataFlow == Bidirectional, false})
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
		printer.Print(prefix, technology.Name)
		prefix = ", "
	}
	printer.PrintLn("\"")
}

func printDatabases(databases []*Database, printer *Printer) {
	for _, database := range databases {
		printDataStore(&database.DataStore, "_db", "Database", printer)
	}
}

func printDataStore(dataStore *DataStore, suffix string, tag string, printer *Printer) {
	printer.PrintLn(dataStore.Id, suffix, " = container \"", dataStore.Name, "\" {")
	printer.Start()
	printer.PrintLn("tags \"", tag, "\"", " \"", dataStore.State.String(), "\"")
	printDescription(dataStore, printer)
	printTechnology(dataStore, printer)
	printer.End()
	printer.PrintLn("}")
}

func printQueues(dataStores []*DataStore, printer *Printer) {
	for _, dataStore := range dataStores {
		printDataStore(dataStore, "_q", "Queue", printer)
	}
}

func printExternalSystems(model *ArchitectureModel, printer *Printer) []usage {
	usages := make([]usage, 0)
	for _, externalSystem := range model.ExternalSystems {
		printer.PrintLn(externalSystem.Id, " = softwareSystem \"", externalSystem.Name, "\" {")
		printer.Start()
		printer.PrintLn("tags \"External System\"")
		printDescription(externalSystem, printer)
		for _, call := range externalSystem.Calls {
			if call.ExternalSystemId != "" {
				usages = append(usages, usage{externalSystem.Id, call.ExternalSystemId, call.Description,
					call.DataFlow == Bidirectional, false})
			} else if call.ServiceId != "" {
				usages = append(usages, usage{externalSystem.Id, call.ServiceId, call.Description,
					call.DataFlow == Bidirectional, false})
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
		if usage.byPersona {
			printer.PrintLn(" {")
			printer.Start()
			printer.PrintLn("tags \"Using\"")
			printer.End()
			printer.Print("}")
		} else if usage.bidirectional {
			printer.NewLine()
			printer.Print(usage.used, " -> ", usage.user)
		}
		printer.NewLine()
	}
}

func printC4views(printer *Printer) {
	printer.PrintLn("views {")
	printer.Start()
	printer.PrintLn("systemContext ", idOfSystemOfInterest, " {")
	printer.Start()
	printer.PrintLn("include *")
	printer.PrintLn("autolayout")
	printer.End()
	printer.PrintLn("}")
	printer.PrintLn("container ", idOfSystemOfInterest, " {")
	printer.Start()
	printer.PrintLn("include *")
	printer.PrintLn("autolayout")
	printer.End()
	printer.PrintLn("}")
	printStyles(printer)
	printer.End()
	printer.PrintLn("}")
}

func printStyles(printer *Printer) {
	printer.PrintLn("styles {")
	printer.Start()

	printElementStyles(printer)
	printStateStyles(printer)
	printRelationshipStyles(printer)

	printer.End()
	printer.PrintLn("}")
}

func printElementStyles(printer *Printer) {
	printer.PrintLn("element \"Person\" {")
	printer.Start()
	printer.PrintLn("shape Person")
	printer.PrintLn("stroke #3966a0")
	printer.PrintLn("strokeWidth 10")
	printer.PrintLn("background GhostWhite")
	printer.End()
	printer.PrintLn("}")

	printer.PrintLn("element \"System of Interest\" {")
	printer.Start()
	printer.PrintLn("background #b6d7a8")
	printer.PrintLn("fontSize 36")
	printer.PrintLn("shape RoundedBox")
	printer.PrintLn("stroke black")
	printer.End()
	printer.PrintLn("}")

	printer.PrintLn("element \"External System\" {")
	printer.Start()
	printer.PrintLn("background #e2e2e2")
	printer.PrintLn("shape RoundedBox")
	printer.PrintLn("stroke black")
	printer.End()
	printer.PrintLn("}")

	printer.PrintLn("element \"Service\" {")
	printer.Start()
	printer.PrintLn("stroke black")
	printer.End()
	printer.PrintLn("}")

	printer.PrintLn("element \"Database\" {")
	printer.Start()
	printer.PrintLn("shape Cylinder")
	printer.PrintLn("stroke black")
	printer.End()
	printer.PrintLn("}")

	printer.PrintLn("element \"Queue\" {")
	printer.Start()
	printer.PrintLn("shape Pipe")
	printer.PrintLn("stroke black")
	printer.End()
	printer.PrintLn("}")
}

func printStateStyles(printer *Printer) {
	printStateStyle(Ok, "b6d7a8", printer)
	printStateStyle(Emerging, "a4c1f4", printer)
	printStateStyle(Review, "fff2cc", printer)
	printStateStyle(Revision, "ffd966", printer)
	printStateStyle(Legacy, "e69238", printer)
	printStateStyle(Deprecated, "cc0000", printer)
}

func printStateStyle(state State, backgroundColor string, printer *Printer) {
	printer.PrintLn("element \"", state.String(), "\" {")
	printer.Start()
	printer.PrintLn("background #", backgroundColor)
	printer.End()
	printer.PrintLn("}")
}

func printRelationshipStyles(printer *Printer) {
	printer.PrintLn("relationship \"Relationship\" {")
	printer.Start()
	printer.PrintLn("color black")
	printer.PrintLn("routing Curved")
	printer.PrintLn("style solid")
	printer.PrintLn("thickness 2")
	printer.End()
	printer.PrintLn("}")

	printer.PrintLn("relationship \"Using\" {")
	printer.Start()
	printer.PrintLn("color #60327c")
	printer.PrintLn("thickness 5")
	printer.End()
	printer.PrintLn("}")
}
