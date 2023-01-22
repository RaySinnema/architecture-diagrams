package main

import (
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
	"unicode"
)

type Nameable interface {
	setNode(node *yaml.Node)
	setId(id string)
	setName(name string)
}

func namedObject(node *yaml.Node, id string, nameable Nameable) (map[string]*yaml.Node, []Issue) {
	nameable.setNode(node)
	nameable.setId(id)
	fields, issue := toMap(node)
	if issue != nil {
		return nil, []Issue{*issue}
	}
	return fields, setName(fields, nameable, id)
}

func setName(fields map[string]*yaml.Node, nameable Nameable, id string) []Issue {
	name, found, issue := stringFieldOf(fields, "name")
	if issue != nil {
		return []Issue{*issue}
	}
	if found {
		nameable.setName(name)
	} else {
		nameable.setName(friendlyNameFrom(id))
	}
	return []Issue{}
}

func friendlyNameFrom(value string) string {
	ext := filepath.Ext(value)
	name := strings.TrimSuffix(value, ext)
	name = strings.Replace(name, "-", " ", -1)
	if len(name) > 0 {
		runes := []rune(name)
		runes[0] = unicode.ToUpper(runes[0])
		name = string(runes)
	}
	return name
}
