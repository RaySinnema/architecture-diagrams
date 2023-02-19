package main

type C4Exporter struct {
}

const idOfSystemOfInterest = "system"

func (c C4Exporter) export(model ArchitectureModel, printer *Printer) {
	printer.PrintLn("workspace {")
	printer.Start()
	c.printModel(model, printer)
	c.printViews(printer)
	printer.End()
	printer.PrintLn("}")
}

type usage struct {
	user        string
	used        string
	description string
	byPersona   bool
}

func (c C4Exporter) printModel(model ArchitectureModel, printer *Printer) {
	printer.PrintLn("model {")
	printer.Start()

	usages := make([]usage, 0)
	usages = append(usages, c.printPersons(model, printer)...)
	usages = append(usages, c.printSoftwareSystems(model, printer)...)
	c.printRelationships(usages, printer)

	printer.End()
	printer.PrintLn("}")
}

func (c C4Exporter) printPersons(model ArchitectureModel, printer *Printer) []usage {
	usages := make([]usage, 0)
	for _, persona := range model.Personas {
		printer.PrintLn(persona.Id, " = person \"", persona.Name, "\" {")
		printer.Start()
		c.printDescription(persona, printer)
		for _, used := range persona.Uses {
			if used.ExternalSystem != nil {
				usages = append(usages, usage{persona.Id, used.ExternalSystem.Id, used.Description, true})
			} else if used.Form != nil {
				usages = append(usages, usage{persona.Id, used.Form.ImplementedBy.Id, used.Description, true})
			} else if used.View != nil {
				usages = append(usages, usage{persona.Id, used.View.On.Id + "_db", used.Description, true})
			}
		}
		printer.End()
		printer.PrintLn("}")
	}
	return usages
}

func (c C4Exporter) printDescription(describable Describable, printer *Printer) {
	if describable.getDescription() != "" {
		printer.PrintLn("description \"", describable.getDescription(), "\"")
	}
}

func (c C4Exporter) printSoftwareSystems(model ArchitectureModel, printer *Printer) []usage {
	var usages = make([]usage, 0)
	usages = append(usages, c.printSoftwareSystemOfInterest(model, printer)...)
	usages = append(usages, c.printExternalSystems(model, printer)...)
	return usages
}

func (c C4Exporter) printSoftwareSystemOfInterest(model ArchitectureModel, printer *Printer) []usage {
	printer.PrintLn(idOfSystemOfInterest, " = softwareSystem \"", model.System.Name, "\" {")
	printer.Start()
	printer.PrintLn("tags \"System of Interest\"")
	usages := c.printContainers(model, printer)
	printer.End()
	printer.PrintLn("}")
	return usages
}

func (c C4Exporter) printContainers(model ArchitectureModel, printer *Printer) []usage {
	usages := make([]usage, 0)
	usages = append(usages, c.printServices(model.Services, printer)...)
	c.printDatabases(model.Databases, printer)
	c.printQueues(model.Queues, printer)
	return usages
}

func (c C4Exporter) printServices(services []*Service, printer *Printer) []usage {
	usages := make([]usage, 0)
	for _, service := range services {
		printer.PrintLn(service.Id, " = container \"", service.Name, "\" {")
		printer.Start()
		c.printDescription(service, printer)
		c.printTechnology(service, printer)
		printer.PrintLn("tags \"Service\" \"", service.State.String(), "\"")
		for _, call := range service.Calls {
			if call.ExternalSystemId != "" {
				usages = append(usages, usage{service.Id, call.ExternalSystemId, call.Description, false})
			} else {
				usages = append(usages, usage{service.Id, call.ServiceId, call.Description, false})
			}
		}
		for _, dataStore := range service.DataStores {
			id := dataStore.QueueId
			if id == "" {
				id = dataStore.DatabaseId + "_db"
			} else {
				id = id + "_q"
			}
			usages = append(usages, usage{service.Id, id, dataStore.Description, false})
		}
		printer.End()
		printer.PrintLn("}")
	}
	return usages
}

func (c C4Exporter) printTechnology(implementable Implementable, printer *Printer) {
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

func (c C4Exporter) printDatabases(databases []*Database, printer *Printer) {
	for _, database := range databases {
		c.printDataStore(&database.DataStore, "_db", "Database", printer)
	}
}

func (c C4Exporter) printDataStore(dataStore *DataStore, suffix string, tag string, printer *Printer) {
	printer.PrintLn(dataStore.Id, suffix, " = container \"", dataStore.Name, "\" {")
	printer.Start()
	printer.PrintLn("tags \"", tag, "\"", " \"", dataStore.State.String(), "\"")
	c.printDescription(dataStore, printer)
	c.printTechnology(dataStore, printer)
	printer.End()
	printer.PrintLn("}")
}

func (c C4Exporter) printQueues(dataStores []*DataStore, printer *Printer) {
	for _, dataStore := range dataStores {
		c.printDataStore(dataStore, "_q", "Queue", printer)
	}
}

func (c C4Exporter) printExternalSystems(model ArchitectureModel, printer *Printer) []usage {
	usages := make([]usage, 0)
	for _, externalSystem := range model.ExternalSystems {
		printer.PrintLn(externalSystem.Id, " = softwareSystem \"", externalSystem.Name, "\" {")
		printer.Start()
		printer.Print("tags \"External System\"")
		if externalSystem.Type != "" {
			printer.Print(" \"", externalSystem.Type, "\"")
		}
		printer.NewLine()
		c.printDescription(externalSystem, printer)
		for _, call := range externalSystem.Calls {
			if call.ExternalSystemId != "" {
				usages = append(usages, usage{externalSystem.Id, call.ExternalSystemId, call.Description, false})
			} else if call.ServiceId != "" {
				usages = append(usages, usage{externalSystem.Id, call.ServiceId, call.Description, false})
			}
		}
		printer.End()
		printer.PrintLn("}")
	}
	return usages
}

func (c C4Exporter) printRelationships(usages []usage, printer *Printer) {
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
		}
		printer.NewLine()
	}
}

func (c C4Exporter) printViews(printer *Printer) {
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
	c.printStyles(printer)
	printer.End()
	printer.PrintLn("}")
}

func (c C4Exporter) printStyles(printer *Printer) {
	printer.PrintLn("styles {")
	printer.Start()

	c.printElementStyles(printer)
	c.printStateStyles(printer)
	c.printRelationshipStyles(printer)

	printer.End()
	printer.PrintLn("}")
}

func (c C4Exporter) printElementStyles(printer *Printer) {
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

	printer.PrintLn("element \"central\" {")
	printer.Start()
	printer.PrintLn("background #a2c4c9")
	printer.End()
	printer.PrintLn("}")

	printer.PrintLn("element \"local\" {")
	printer.Start()
	printer.PrintLn("background #38761d")
	printer.PrintLn("color white")
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

func (c C4Exporter) printStateStyles(printer *Printer) {
	c.printStateStyle(Ok, "b6d7a8", printer)
	c.printStateStyle(Emerging, "a4c1f4", printer)
	c.printStateStyle(Review, "fff2cc", printer)
	c.printStateStyle(Revision, "ffd966", printer)
	c.printStateStyle(Legacy, "e69238", printer)
	c.printStateStyle(Deprecated, "cc0000", printer)
}

func (c C4Exporter) printStateStyle(state State, backgroundColor string, printer *Printer) {
	printer.PrintLn("element \"", state.String(), "\" {")
	printer.Start()
	printer.PrintLn("background #", backgroundColor)
	printer.End()
	printer.PrintLn("}")
}

func (c C4Exporter) printRelationshipStyles(printer *Printer) {
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
