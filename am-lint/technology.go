package main

import (
	"gopkg.in/yaml.v3"
	"sort"
)

type Implementable interface {
	setTechnologyIds(technologies []string)
	setTechnologyBundleId(technologyBundle string)
}

func setTechnologies(fields map[string]*yaml.Node, implementable Implementable) []Issue {
	issues := make([]Issue, 0)
	const technologiesField = "technologies"
	technologyBundle, found, issue := stringFieldOf(fields, technologiesField)
	if issue == nil {
		if found {
			implementable.setTechnologyBundleId(technologyBundle)
		}
		return issues
	}
	technologiesNodes, _, issue := sequenceFieldOf(fields, technologiesField)
	if issue == nil {
		technologies := make([]string, 0)
		for _, technologyNode := range technologiesNodes {
			technology, issue := toString(technologyNode, "technology")
			if issue == nil {
				technologies = append(technologies, technology)
			} else {
				issues = append(issues, *issue)
			}
		}
		sort.Slice(technologies, func(i, j int) bool {
			return technologies[i] < technologies[j]
		})
		implementable.setTechnologyIds(technologies)
	} else {
		issues = append(issues, *issue)
	}

	return issues
}
