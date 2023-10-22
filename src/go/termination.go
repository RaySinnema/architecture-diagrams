package main

import "math"

type Termination[G Cloner[G]] interface {
	isMet(population Population[G]) bool
}

func NewMaxIterations[G Cloner[G]](max int) Termination[G] {
	return &maxIterations[G]{max, 0}
}

type maxIterations[G Cloner[G]] struct {
	maxIterations    int
	currentIteration int
}

func (mi *maxIterations[G]) isMet(_ Population[G]) bool {
	result := mi.currentIteration >= mi.maxIterations
	if !result {
		mi.currentIteration++
	}
	return result
}

func NewOrTermination[G Cloner[G]](terminations ...Termination[G]) Termination[G] {
	return &orTermination[G]{terminations}
}

type orTermination[G Cloner[G]] struct {
	terminations []Termination[G]
}

func (o *orTermination[G]) isMet(population Population[G]) bool {
	for _, termination := range o.terminations {
		if termination.isMet(population) {
			return true
		}
	}
	return false
}

func NewMaxFitnessAchieved[G Cloner[G]](target Fitness) Termination[G] {
	return &maxFitnessAchived[G]{target}
}

type maxFitnessAchived[G Cloner[G]] struct {
	target Fitness
}

func (m maxFitnessAchived[G]) isMet(population Population[G]) bool {
	return math.Abs(m.target-population.best().Fitness) < 1e-6
}

func NewNoProgressMade[G Cloner[G]](iterations int) Termination[G] {
	return &noProgressMade[G]{0, make([]Fitness, iterations)}
}

type noProgressMade[G Cloner[G]] struct {
	index     int
	fitnesses []Fitness
}

func (f *noProgressMade[G]) isMet(population Population[G]) bool {
	f.fitnesses[f.index] = population.avgFitness()
	if f.index < len(f.fitnesses) {
		f.index++
		return false
	}
	f.index = (f.index + 1) % len(f.fitnesses)
	prev := f.fitnesses[0]
	const maxFitnessDelta = 1e-6
	for i := 1; i < len(f.fitnesses); i++ {
		if math.Abs(prev-f.fitnesses[i]) > maxFitnessDelta {
			return false
		}
		prev = f.fitnesses[i]
	}
	return true
}
