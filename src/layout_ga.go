package main

import (
	"math"
	"math/rand"
)

func NewEvolutionaryLayoutEngine() LayoutEngine {
	return evolutionaryLayoutEngine{}
}

type evolutionaryLayoutEngine struct {
}

const (
	layoutGaPoolSize      = 100
	layoutGaMaxIterations = 100

	// Operator probabilities
	probabilityMutation  = -1.0 // Negative -> approx. one mutation per genome
	probabilityCrossOver = 0.6

	// Fitness function weights
	weightEdgeCrossings       = 0.3
	weightEdgeOverlaps        = 0.05
	weightSymmetricConnectors = 0.15
	weightEmptyRowsAndColumns = 0.3
	weightNodeTypes           = 0.2
)

func (e evolutionaryLayoutEngine) layOut(diagram *Diagram) *DiagramLayout {
	size := calcGridSize(diagram)
	ga := NewGeneticAlgorithm[*DiagramGene](
		layoutGaPoolSize,
		func() *Genome[*DiagramGene] { return createDiagramGenome(size, diagram) },
		[]Operator[*DiagramGene]{
			NewBinaryOperator[*DiagramGene](probabilityCrossOver, crossOverDiagram),
			NewUnaryOperator[*DiagramGene](probabilityMutation, mutateDiagram),
		},
		NewWeightedAverage[*DiagramGene](partialFitnessFunctions()),
		NewElitistSelection[*DiagramGene](3, NewRouletteWheelSelection[*DiagramGene]()),
		NewOrTermination[*DiagramGene](
			NewMaxIterations[*DiagramGene](layoutGaMaxIterations),
			NewNoProgressMade[*DiagramGene](5),
		),
	)
	ga.run()
	return toLayout(ga.best())
}

func toLayout(genome *Genome[*DiagramGene]) *DiagramLayout {
	// TODO: Implement
	return nil
}

func calcGridSize(diagram *Diagram) Size {
	sideSize := 2 * int(math.Ceil(math.Sqrt(3.0*float64(len(diagram.Shapes)))))
	return Size{sideSize, sideSize}
}

func mutateDiagram(g *Genome[*DiagramGene], position int) *DiagramGene {
	// TODO: Implement
	return g.Genes[position]
}

func crossOverDiagram(parent1 *Genome[*DiagramGene], parent2 *Genome[*DiagramGene]) (*Genome[*DiagramGene], *Genome[*DiagramGene]) {
	// TODO: Implement
	return parent1, parent2
}

func partialFitnessFunctions() []PartialFitnessFunction[*DiagramGene] {
	return []PartialFitnessFunction[*DiagramGene]{
		{weightEdgeCrossings, calcEdgeCrossings},
		{weightEdgeOverlaps, calcEdgeOverlaps},
		{weightSymmetricConnectors, calcSymmetricConnectors},
		{weightEmptyRowsAndColumns, calcEmptyRowsAndColumns},
		{weightNodeTypes, calcNodeTypes},
	}
}

func calcEdgeCrossings(genome *Genome[*DiagramGene]) Fitness {
	// TODO: Implement
	return 0.0
}

func calcEdgeOverlaps(genome *Genome[*DiagramGene]) Fitness {
	// TODO: Implement
	return 0.0
}

func calcSymmetricConnectors(genome *Genome[*DiagramGene]) Fitness {
	// TODO: Implement
	return 0.0
}

func calcEmptyRowsAndColumns(genome *Genome[*DiagramGene]) Fitness {
	// TODO: Implement
	return 0.0
}

func calcNodeTypes(genome *Genome[*DiagramGene]) Fitness {
	// TODO: Implement
	return 0.0
}

type connector struct {
	side  Side
	index int
}

type miniGridIndex struct {
}

type path struct {
	from  connector
	to    connector
	bends []miniGridIndex
}

type DiagramGene struct {
	context         *diagramContext
	shapeIndex      int
	gridPosition    int
	connectionIndex int
	path            *path
}

func (d *DiagramGene) clone() *DiagramGene {
	clone := *d
	return &clone
}

type diagramContext struct {
	size    Size
	diagram *Diagram
}

func createDiagramGenome(size Size, diagram *Diagram) *Genome[*DiagramGene] {
	context := &diagramContext{size, diagram}
	genes := make([]*DiagramGene, len(diagram.Shapes)+len(diagram.Connections))
	positions := zeroTo(size.Width * size.Height)
	max := int32(len(positions))
	for index := 0; index < len(diagram.Shapes); index++ {
		pos := rand.Int31n(max)
		shapeGene := DiagramGene{context, index, positions[pos], -1, nil}
		genes[index] = &shapeGene
		max--
		positions[max], positions[pos] = positions[pos], positions[max]
	}
	for index := 0; index < len(diagram.Connections); index++ {
		connection := diagram.Connections[index]
		path := defaultPathBetween(connection.Start, connection.End)
		connectionGene := DiagramGene{context, -1, -1, index, path}
		genes[index] = &connectionGene
	}
	return &Genome[*DiagramGene]{genes, 0.0}
}

func defaultPathBetween(start *Shape, end *Shape) *path {
	// TODO: Implement
	return nil
}

func zeroTo(size int) []int {
	result := make([]int, size)
	for i := 0; i < size; i++ {
		result[i] = i
	}
	return result
}
