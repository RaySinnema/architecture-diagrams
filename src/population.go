package main

type Fitness = float64

type Cloner[C any] interface {
	clone() C
}

type Genome[G Cloner[G]] struct {
	Genes   []G
	Fitness Fitness
}

func (g Genome[G]) clone() Genome[G] {
	genes := make([]G, len(g.Genes))
	for i := 0; i < len(genes); i++ {
		genes[i] = g.Genes[i].clone()
	}
	return Genome[G]{genes, g.Fitness}
}

type Population[G Cloner[G]] struct {
	genomes []*Genome[G]
}

func (p *Population[G]) set(genomes []*Genome[G]) {
	p.genomes = genomes
}

func (p *Population[G]) add(genome *Genome[G]) {
	p.genomes = append(p.genomes, genome)
}

func (p *Population[G]) calculateFitness(fitnessFunction FitnessFunction[G]) {
	for _, genome := range p.genomes {
		genome.Fitness = fitnessFunction.calculate(genome)
	}
}

func (p *Population[G]) all() []*Genome[G] {
	return p.genomes
}

func (p *Population[G]) best() *Genome[G] {
	var result *Genome[G] = nil
	bestFitness := -1.0
	for _, genome := range p.genomes {
		if genome.Fitness > bestFitness {
			bestFitness = genome.Fitness
			result = genome
		}
	}
	return result
}

func (p *Population[G]) avgFitness() Fitness {
	totalFitness := 0.0
	for _, genome := range p.genomes {
		totalFitness += genome.Fitness
	}
	return totalFitness / Fitness(len(p.genomes))
}
