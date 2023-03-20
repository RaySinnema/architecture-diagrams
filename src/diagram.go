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

func (s Size) area() int {
	return s.Width * s.Height
}

type Shape struct {
	Id            string
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

func (c Connection) connectsTo(shape *Shape) bool {
	return c.Start == shape || c.End == shape
}

type Diagram struct {
	Shapes      []*Shape
	Connections []*Connection
}

type Diagrammer interface {
	toDiagram(model *ArchitectureModel) *Diagram
}
