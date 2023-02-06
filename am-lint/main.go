package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	dir := "../examples"
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Can't read examples: %v\n", err)
		return
	}
	found := false
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		found = true
		printExample(fmt.Sprintf("%s/%s", dir, entry.Name()))
	}
	if !found {
		fmt.Println("No examples found")
	}
}

func printExample(fileName string) {
	fmt.Println(fileName)
	model, issues := LintFile(fileName)

	if len(issues) > 0 {
		fmt.Println("Issues:")
		for _, issue := range issues {
			fmt.Println(issue)
		}
		fmt.Println("")
	}
	if model != nil {
		fmt.Println(model.String())
	}
}
