package main

import "fmt"

type eventModelExporter struct {
	workflowId string
}

func NewEventModelExporter(workflowId string) TextExporter {
	return eventModelExporter{workflowId}
}

func (e eventModelExporter) export(model ArchitectureModel, printer *Printer) error {
	workflow, found := model.findWorkflowById(e.workflowId)
	if !found {
		return fmt.Errorf("unknown workflow '%v'", e.workflowId)
	}
	e.printWorkflow(workflow, &model, printer)
	return nil
}

func (e eventModelExporter) printWorkflow(workflow *Workflow, model *ArchitectureModel, printer *Printer) {
	printer.PrintLn("direction: right")
	printer.PrintLn(workflow.Name, ": {")
	printer.Start()
	e.printLanes(workflow, model, printer)
	printer.End()
	printer.PrintLn("}")
}

func (e eventModelExporter) printLanes(workflow *Workflow, model *ArchitectureModel, printer *Printer) {
	// TODO: Implement
}
