package main

import (
	"path/filepath"
	"strings"
	"unicode"
)

func friendly(value string) string {
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
