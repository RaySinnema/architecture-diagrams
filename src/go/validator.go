package main

type Validator interface {
	validate(model *ArchitectureModel) []Issue
}
