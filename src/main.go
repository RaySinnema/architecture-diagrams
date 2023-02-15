package main

import (
	"flag"
	"fmt"
)

func main() {
	var command string
	var fileName string
	var output string

	flag.StringVar(&command, "c", "lint", "Command.")
	flag.StringVar(&fileName, "f", "", "Name of model file")
	flag.StringVar(&output, "o", "", "Name of output file")
	flag.Parse()

	switch command {
	case "lint":
		if fileName == "" {
			flag.PrintDefaults()
		} else {
			lintFile(fileName)
		}
	case "c4":
		if fileName == "" || output == "" {
			flag.PrintDefaults()
		} else {
			c4(fileName, output)
		}
	}
}

func lintFile(fileName string) {
	_, issues := LintFile(fileName)
	if len(issues) > 0 {
		listIssues(fileName, issues)
	} else {
		fmt.Printf("%v is OK\n", fileName)
	}
}

func listIssues(fileName string, issues []Issue) {
	fmt.Printf("Issues for %v\n", fileName)
	for _, issue := range issues {
		fmt.Printf("%s\n", issue)
	}
}

func c4(input string, output string) {
	model, issues := LintFile(input)
	if model == nil {
		listIssues(input, issues)
	} else {
		GenerateC4(model, output)
	}
}
