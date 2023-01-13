package main

import (
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
	"unicode"
)

type Nameable interface {
	setName(name string)
}

func namedObject(id string, node *yaml.Node, nameable Nameable, issues []Issue) (map[string]*yaml.Node, []Issue) {
	fields, issue := toMap(node)
	if issue != nil {
		return nil, append(issues, *issue)
	}
	name, found, issue := stringFieldOf(fields, "name")
	if issue != nil {
		return nil, append(issues, *issue)
	} else if found {
		nameable.setName(name)
	} else {
		nameable.setName(friendlyNameFrom(id))
	}
	return fields, issues
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
