package main

import "image"

type DiagramLayout struct {
	Shapes      map[*Shape]image.Rectangle
	Connections map[*Connection][]image.Point
}

type LayoutEngine interface {
	layOut(diagram *Diagram) *DiagramLayout
}
