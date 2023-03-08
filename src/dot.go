package main

import "fmt"

type dotExporter struct {
}

func NewDotExporter() TextExporter {
	return dotExporter{}
}

func (d dotExporter) export(model ArchitectureModel, printer *Printer) error {
	printer.PrintLn("digraph {")
	printer.Start()
	printer.PrintLn("splines=ortho")
	printer.PrintLn()
	d.printModel(&model, printer)
	printer.End()
	printer.PrintLn("}")
	return nil
}

func (d dotExporter) printModel(model *ArchitectureModel, printer *Printer) {
	d.printPersonas(model.Personas, printer)
	d.printExternalSystems(model.ExternalSystems, printer)
	d.printServices(model.Services, printer)
	d.printDatabase(model.Databases, printer)
	d.printQueue(model.Queues, printer)
}

func (d dotExporter) printPersonas(personas []*Persona, printer *Printer) {
	for _, persona := range personas {
		printer.PrintLn(persona.Id, "[shape=polygon,sides=5,color=\"#3966a0\",label=\"", persona.Name, "\"]")
		for _, use := range persona.Uses {
			target := d.using(use)
			if target != "" {
				printer.PrintLn(persona.Id, " -> ", target, " [dir=", d.directionOf(use.DataFlow), "]")
			}
		}
	}
	printer.PrintLn()
}

func (d dotExporter) using(use *Used) string {
	if use.Form != nil {
		return use.Form.ImplementedBy.Id
	}
	if use.View != nil {
		return fmt.Sprintf("%s_db", use.View.On.Id)
	}
	if use.ExternalSystem != nil {
		return use.ExternalSystem.Id
	}
	return ""
}

func (d dotExporter) directionOf(dataFlow DataFlow) string {
	switch dataFlow {
	case Send:
		return "forward"
	case Receive:
		return "back"
	case Bidirectional:
		return "both"
	default:
		return "none"
	}
}

func (d dotExporter) printExternalSystems(externalSystems []*ExternalSystem, printer *Printer) {
	for _, externalSystem := range externalSystems {
		printer.PrintLn(externalSystem.Id, "[shape=rectangle,style=\"rounded,filled\",fillcolor=\"#e2e2e2\",label=\"",
			externalSystem.Name, "\"]")
		for _, call := range externalSystem.Calls {
			target := d.calling(call)
			if target != "" {
				printer.PrintLn(externalSystem.Id, " -> ", target, " [dir=", d.directionOf(call.DataFlow), "]")
			}
		}
	}
	printer.PrintLn()
}

func (d dotExporter) calling(call *Call) string {
	if call.Service != nil {
		return call.Service.Id
	}
	if call.ExternalSystem != nil {
		return call.ExternalSystem.Id
	}
	return ""
}

func (d dotExporter) printServices(services []*Service, printer *Printer) {
	for _, service := range services {
		printer.PrintLn(service.Id, "[shape=box,style=filled,fillcolor=\"#", d.colorOf(service.State), "\",label=\"",
			service.Name, "\"]")
		for _, call := range service.Calls {
			target := d.calling(call)
			if target != "" {
				printer.PrintLn(service.Id, " -> ", target, " [dir=", d.directionOf(call.DataFlow), "]")
			}
		}
		for _, dataStore := range service.DataStores {
			target := d.storingIn(dataStore)
			if target != "" {
				printer.PrintLn(service.Id, " -> ", target, " [dir=", d.directionOf(dataStore.DataFlow), "]")
			}
		}
	}
	printer.PrintLn()
}

func (d dotExporter) colorOf(state State) string {
	switch state {
	case Ok:
		return "b6d7a8"
	case Emerging:
		return "a4c1f4"
	case Review:
		return "fff2cc"
	case Revision:
		return "ffd966"
	case Legacy:
		return "e69238"
	case Deprecated:
		return "cc0000"
	default:
		return "ghostwhite"
	}
}

func (d dotExporter) storingIn(store *DataStoreUse) string {
	if store.Database != nil {
		return fmt.Sprintf("%s_db", store.Database.Id)
	}
	if store.Queue != nil {
		return fmt.Sprintf("%s_q", store.Queue.Id)
	}
	return ""
}

func (d dotExporter) printDatabase(databases []*Database, printer *Printer) {
	for _, database := range databases {
		printer.PrintLn(database.Id, "_db [shape=cylinder,style=filled,fillcolor=\"#", d.colorOf(database.State),
			"\",label=\"", database.Name, "\"]")
	}
}

func (d dotExporter) printQueue(queues []*DataStore, printer *Printer) {
	for _, queue := range queues {
		printer.PrintLn(queue.Id, "_q [shape=parallelogram,style=filled,fillcolor=\"#", d.colorOf(queue.State),
			"\",label=\"", queue.Name, "\"]")
	}
}
