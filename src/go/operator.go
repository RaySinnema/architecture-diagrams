package main

import "math/rand"

type Operator[G Cloner[G]] interface {
	operateOn(genomes []*Genome[G])
}

// NewUnaryOperator Create a new unary operator.
// The operator performs the operation for every gene with the given probability.
// If the probability is negative, it will be set to 1 / number of genes, so that on average one gene per genome is
// modified.
func NewUnaryOperator[G Cloner[G]](probability float64, operator func(*Genome[G], int) G) Operator[G] {
	return unaryOperator[G]{probability, operator}
}

type unaryOperator[G Cloner[G]] struct {
	probability float64
	operator    func(*Genome[G], int) G
}

func (u unaryOperator[G]) operateOn(genomes []*Genome[G]) {
	probability := u.probability
	if probability < 0 {
		probability = 1.0 / float64(len(genomes[0].Genes))
	}
	for index := 0; index < len(genomes); index++ {
		genome := genomes[index]
		for index := range genome.Genes {
			if rand.Float64() < probability {
				genome.Genes[index] = u.operator(genome, index)
			}
		}
	}
}

func NewBinaryOperator[G Cloner[G]](probability float64, operator func(*Genome[G], *Genome[G]) (*Genome[G], *Genome[G])) Operator[G] {
	return binaryOperator[G]{probability, operator}
}

type binaryOperator[G Cloner[G]] struct {
	probability float64
	operator    func(*Genome[G], *Genome[G]) (*Genome[G], *Genome[G])
}

func (b binaryOperator[G]) operateOn(genomes []*Genome[G]) {
	for index := 0; index < len(genomes)-1; index = index + 2 {
		if rand.Float64() < b.probability {
			output1, output2 := b.operator(genomes[index], genomes[index+1])
			genomes[index] = output1
			genomes[index+1] = output2
		}
	}
}
