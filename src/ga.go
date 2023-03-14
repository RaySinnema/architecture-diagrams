package main

import (
	"fmt"
)

type GeneticAlgorithm[G Cloner[G]] interface {
	run()
	best() *Genome[G]
}

func NewGeneticAlgorithm[G Cloner[G]](
	poolSize int,
	createGenome func() *Genome[G],
	operators []Operator[G],
	fitnessFunction FitnessFunction[G],
	selection Selection[G],
	termination Termination[G],
) GeneticAlgorithm[G] {
	return &ga[G]{makePool(poolSize, createGenome), operators, fitnessFunction,
		selection, termination}
}

func makePool[G Cloner[G]](size int, createGenome func() *Genome[G]) Population[G] {
	result := Population[G]{}
	genomes := make([]*Genome[G], size)
	for i := 0; i < size; i++ {
		genomes[i] = createGenome()
	}
	result.set(genomes)
	return result
}

type ga[G Cloner[G]] struct {
	pool            Population[G]
	operators       []Operator[G]
	fitnessFunction FitnessFunction[G]
	selection       Selection[G]
	termination     Termination[G]
}

func (ga *ga[G]) best() *Genome[G] {
	return ga.pool.best()
}

func (ga *ga[G]) run() {
	ga.calculateFitness()
	ga.dumpPool("Initialized")
	for iteration := 0; !ga.termination.isMet(ga.pool); iteration++ {
		ga.refreshPool()
		ga.calculateFitness()
		ga.dumpPool(fmt.Sprintf("Iteration %d", iteration))
	}
}

func (ga *ga[G]) calculateFitness() {
	ga.pool.calculateFitness(ga.fitnessFunction)
}

func (ga *ga[G]) refreshPool() {
	newPool := ga.selection.selectFrom(ga.pool.all())
	for _, operator := range ga.operators {
		operator.operateOn(newPool)
	}
	ga.pool.set(newPool)
}

func (ga *ga[G]) dumpPool(message string) {
	fmt.Printf("\n%s:\n", message)
	total := 0.0
	count := 0.0
	max := 0.0
	for _, dna := range ga.pool.all() {
		fmt.Printf("%v => %0.4f\n", dna.Genes, dna.Fitness)
		total += dna.Fitness
		count += 1.0
		if dna.Fitness > max {
			max = dna.Fitness
		}
	}
	fmt.Printf("avg %0.4f, max %0.4f\n", total/count, max)
}
