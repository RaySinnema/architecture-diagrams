package main

func NewContextDiagrammer() Diagrammer {
	return &contextDiagrammer{}
}

type contextDiagrammer struct {
}

func (c contextDiagrammer) toDiagram(model *ArchitectureModel) *Diagram {
	// TODO: Convert Persona, etc
	return &Diagram{}
}
