package main

import "fmt"

type SystemReader struct {
}

func (_ SystemReader) read(m map[string]interface{}, fileName string, model *ArchitectureModel) []Issue {
	name := systemOf(m, fileName)
	if len(name) == 0 {
		return []Issue{*NewError("Missing system name")}
	}
	model.system.name = name
	return []Issue{}
}

func systemOf(m map[string]interface{}, fileName string) string {
	raw, exists := m["system"]
	if !exists {
		return friendly(fileName)
	}
	system, ok := raw.(map[string]interface{})
	if !ok {
		return friendly(fileName)
	}
	name, exists := system["name"]
	if !exists {
		return friendly(fileName)
	}
	return fmt.Sprintf("%v", name)
}
