package main

import "gopkg.in/yaml.v3"

type Describable interface {
	getDescription() string
	setDescription(description string)
}

func setDescription(fields map[string]*yaml.Node, d Describable) []Issue {
	description, found, issue := stringFieldOf(fields, "description")
	if issue != nil {
		return []Issue{*issue}
	}
	if found {
		d.setDescription(description)
	}
	return []Issue{}
}
