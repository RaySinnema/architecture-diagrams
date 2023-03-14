package main

type FitnessFunction[G Cloner[G]] interface {
	calculate(genome *Genome[G]) Fitness
}

type PartialFitnessFunction[G Cloner[G]] struct {
	Weight    float64
	Calculate func(genome *Genome[G]) Fitness
}

func NewWeightedAverage[G Cloner[G]](partialFitnessFunctions []PartialFitnessFunction[G]) FitnessFunction[G] {
	sum := 0.0
	for _, partialFitnessFunction := range partialFitnessFunctions {
		sum += partialFitnessFunction.Weight
	}
	for _, partialFitnessFunction := range partialFitnessFunctions {
		partialFitnessFunction.Weight /= sum
	}
	return weightedAverage[G]{partialFitnessFunctions}
}

type weightedAverage[G Cloner[G]] struct {
	partialFitnessFunctions []PartialFitnessFunction[G]
}

func (w weightedAverage[G]) calculate(genome *Genome[G]) float64 {
	result := 0.0
	for _, partialFitnessFunction := range w.partialFitnessFunctions {
		result += partialFitnessFunction.Weight * partialFitnessFunction.Calculate(genome)
	}
	return result
}

func NewIncrementalFitnessFunction[G Cloner[G]](geneFitnessFunction func(G) Fitness) FitnessFunction[G] {
	return incrementalFitnessFunction[G]{geneFitnessFunction}
}

type incrementalFitnessFunction[G Cloner[G]] struct {
	geneFitnessFunction func(G) Fitness
}

func (i incrementalFitnessFunction[G]) calculate(genome *Genome[G]) Fitness {
	result := 0.0
	for _, gene := range genome.Genes {
		result += i.geneFitnessFunction(gene)
	}
	return result / Fitness(len(genome.Genes))
}
