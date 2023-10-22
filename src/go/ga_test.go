package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestGeneticAlgorithm(t *testing.T) {
	ga := NewGeneticAlgorithm[bit](10, createBits, makeOperators(), makeFitnessFunction(),
		makeSelection(), makeTermination())

	ga.run()

	best := ga.best().Fitness
	if math.Abs(1.0-best) > 1e-9 {
		t.Fatalf("Didn't find maximum, but %v", best)
	}
}

func createBits() *Genome[bit] {
	bits := make([]bit, 8)
	for i := 0; i < len(bits); i++ {
		if rand.Float64() < 0.5 {
			bits[i] = true
		} else {
			bits[i] = false
		}
	}
	return &Genome[bit]{bits, 0.0}
}

type bit bool

func (b bit) clone() bit {
	return b
}

func (b bit) String() string {
	return fmt.Sprintf("%5t", b)
}

func makeOperators() []Operator[bit] {
	return []Operator[bit]{
		NewBinaryOperator[bit](0.7, crossOver),
		NewUnaryOperator[bit](0.1, mutate),
	}
}

func mutate(genome *Genome[bit], index int) bit {
	return !genome.Genes[index]
}

func crossOver(parent1 *Genome[bit], parent2 *Genome[bit]) (*Genome[bit], *Genome[bit]) {
	size := len(parent1.Genes)
	child1 := make([]bit, size)
	child2 := make([]bit, size)
	crossOverPoint := int(rand.Int31n(int32(size) - 1))
	for index := 0; index < crossOverPoint; index++ {
		child1[index] = parent1.Genes[index]
		child2[index] = parent2.Genes[index]
	}
	for index := crossOverPoint; index < size; index++ {
		child1[index] = parent2.Genes[index]
		child2[index] = parent1.Genes[index]
	}
	return &Genome[bit]{child1, 0.0}, &Genome[bit]{child2, 0.0}
}

func makeFitnessFunction() FitnessFunction[bit] {
	return NewIncrementalFitnessFunction[bit](bitFitness)
}

func bitFitness(b bit) Fitness {
	if b {
		return 1.0
	}
	return 0.0
}

func makeSelection() Selection[bit] {
	return NewElitistSelection[bit](1, NewRouletteWheelSelection[bit]())
}

func makeTermination() Termination[bit] {
	return NewOrTermination(NewMaxFitnessAchieved[bit](1.0), NewMaxIterations[bit](10))
}
