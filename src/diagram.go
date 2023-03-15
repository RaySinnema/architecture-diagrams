package main

type Side int

const (
	Top Side = iota
	Right
	Bottom
	Left
)

type Size struct {
	Width  int
	Height int
}

type Shape struct {
	Text          string
	Size          Size
	NumConnectors map[Side][]int
}

type ConnectionSymbol = int

type Connection struct {
	Start       *Shape
	StartSymbol ConnectionSymbol
	End         *Shape
	EndSymbol   ConnectionSymbol
}

type Diagram struct {
	Shapes      []*Shape
	Connections []*Connection
}

type Diagrammer interface {
	toDiagram(model *ArchitectureModel) *Diagram
}
