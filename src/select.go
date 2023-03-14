package main

import (
	"math/rand"
	"sort"
)

type Selection[G Cloner[G]] interface {
	selectFrom(genomes []*Genome[G]) []*Genome[G]
}

func NewRouletteWheelSelection[G Cloner[G]]() Selection[G] {
	return rouletteWheelSelection[G]{}
}

type rouletteWheelSelection[G Cloner[G]] struct {
}

func (r rouletteWheelSelection[G]) selectFrom(genomes []*Genome[G]) []*Genome[G] {
	totalFitness := 0.0
	for _, genome := range genomes {
		totalFitness += genome.Fitness
	}
	probabilities := make([]float64, len(genomes))
	for index, genome := range genomes {
		probabilities[index] = genome.Fitness / totalFitness
	}
	return r.selectRandomly(genomes, probabilities)
}

func (r rouletteWheelSelection[G]) selectRandomly(genomes []*Genome[G], probabilities []float64) []*Genome[G] {
	result := make([]*Genome[G], len(genomes))
	for index := 0; index < len(genomes); index++ {
		clone := genomes[r.randomIndex(probabilities)].clone()
		result[index] = &clone
	}
	return result
}

func (r rouletteWheelSelection[G]) randomIndex(probabilities []float64) int {
	target := rand.Float64()
	reached := 0.0
	for result := 0; result < len(probabilities); result++ {
		reached += probabilities[result]
		if reached >= target {
			return result
		}
	}
	return len(probabilities) - 1
}

func NewElitistSelection[G Cloner[G]](elites int, remainder Selection[G]) Selection[G] {
	return elitistSelection[G]{elites, remainder}
}

type elitistSelection[G Cloner[G]] struct {
	elites    int
	remainder Selection[G]
}

func (e elitistSelection[G]) selectFrom(genomes []*Genome[G]) []*Genome[G] {
	result := make([]*Genome[G], len(genomes))
	for index, genome := range genomes {
		clone := genome.clone()
		result[index] = &clone
	}
	sort.Slice(result, func(i int, j int) bool {
		return result[i].Fitness > result[j].Fitness
	})
	selected := e.remainder.selectFrom(genomes)
	for i := e.elites; i < len(result); i++ {
		clone := selected[i-e.elites].clone()
		result[i] = &clone
	}
	return result
}
