package main

type dfdExporter struct {
}

func NewDfdExporter() TextExporter {
	return dfdExporter{}
}

func (d dfdExporter) export(model ArchitectureModel, printer *Printer) error {
	d.printPersonas(model.Personas, printer)
	d.printExternalSystems(model.ExternalSystems, printer)
	d.printServices(model.Services, printer)
	d.printDatabases(model.Databases, printer)
	d.printQueues(model.Queues, printer)
	return nil
}

func (d dfdExporter) printPersonas(personas []*Persona, printer *Printer) {
	for _, persona := range personas {
		printer.PrintLn(persona.Id, ": ", persona.Name)
		for _, use := range persona.Uses {
			d.printPersonaUse(persona, use, printer)
		}
	}
}

func (d dfdExporter) printPersonaUse(persona *Persona, use *Used, printer *Printer) {
	printer.Print(persona.Id, d.dataFlowOf(use.DataFlow))
	if use.ExternalSystem != nil {
		printer.Print(use.ExternalSystem.Id)
	} else if use.Form != nil {
		printer.Print(use.Form.ImplementedBy.Id)
	} else if use.View != nil {
		printer.Print(use.View.On.Id, "_db")
	} else {
		panic(use)
	}
	printer.NewLine()
}

func (d dfdExporter) dataFlowOf(dataFlow DataFlow) string {
	switch dataFlow {
	case Bidirectional:
		return " <-> "
	case Receive:
		return " <- "
	case Send:
		return " -> "
	default:
		panic(dataFlow)
	}
}

func (d dfdExporter) printExternalSystems(externalSystems []*ExternalSystem, printer *Printer) {
	for _, externalSystem := range externalSystems {
		printer.PrintLn(externalSystem.Id, ": ", externalSystem.Name)
		for _, call := range externalSystem.Calls {
			d.printCall(externalSystem.Id, call, printer)
		}
	}
}

func (d dfdExporter) printCall(fromId string, call *Call, printer *Printer) {
	printer.Print(fromId, d.dataFlowOf(call.DataFlow))
	if call.ExternalSystem != nil {
		printer.Print(call.ExternalSystem.Id)
	} else if call.Service != nil {
		printer.Print(call.Service.Id)
	} else {
		panic(*call)
	}
	d.printTechnologies(call.Technologies, printer)
	printer.NewLine()
}

func (d dfdExporter) printTechnologies(technologies []*Technology, printer *Printer) {
	if len(technologies) == 0 {
		return
	}
	printer.Print(": ")
	prefix := ""
	for _, technology := range technologies {
		printer.Print(prefix, technology.Name)
		prefix = ", "
	}
}

func (d dfdExporter) printServices(services []*Service, printer *Printer) {
	for _, service := range services {
		printer.PrintLn(service.Id, ": ", service.Name, " { shape: circle }")
		for _, call := range service.Calls {
			d.printCall(service.Id, call, printer)
		}
		for _, dataStore := range service.DataStores {
			d.printDataStoreUse(service.Id, dataStore, printer)
		}
	}
}

func (d dfdExporter) printDataStoreUse(fromId string, use *DataStoreUse, printer *Printer) {
	printer.Print(fromId, d.dataFlowOf(use.DataFlow))
	if use.Database != nil {
		printer.Print(use.Database.Id, "_db")
		d.printTechnologies(use.Database.ApiTechnologies, printer)
	} else if use.Queue != nil {
		printer.Print(use.Queue.Id, "_q")
		d.printTechnologies(use.Queue.ApiTechnologies, printer)
	} else {
		panic(*use)
	}
	printer.NewLine()
}

func (d dfdExporter) printDatabases(databases []*Database, printer *Printer) {
	for _, database := range databases {
		d.printDataStore(database.DataStore, "_db", printer)
	}
}

func (d dfdExporter) printDataStore(dataStore DataStore, suffix string, printer *Printer) {
	printer.PrintLn(dataStore.Id, suffix, ": ", dataStore.Name, " {")
	printer.Start()
	printer.PrintLn("shape: image")
	printer.PrintLn("icon: https://github.com/RemonSinnema/architecture-diagrams/raw/main/static/data-store.png")
	printer.End()
	printer.PrintLn("}")
}

func (d dfdExporter) printQueues(queues []*DataStore, printer *Printer) {
	for _, queue := range queues {
		d.printDataStore(*queue, "_q", printer)
	}
}
