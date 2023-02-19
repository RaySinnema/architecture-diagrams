package main

import (
	"flag"
	"fmt"
	"sort"
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
	case "c4":
		export(fileName, C4Exporter{}, output)
	case "dfd":
		export(fileName, DfdExporter{}, output)
	case "lint":
		lintFile(fileName)
	}
}

func lintFile(fileName string) {
	if fileName == "" {
		flag.PrintDefaults()
		return
	}
	_, issues := LintFile(fileName)
	if len(issues) > 0 {
		listIssues(fileName, issues)
	} else {
		fmt.Printf("%v is OK\n", fileName)
	}
}

func listIssues(fileName string, issues []Issue) {
	fmt.Printf("Issues for %v\n", fileName)
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Line < issues[j].Line
	})

	for _, issue := range issues {
		fmt.Printf("%s\n", issue)
	}
}

func export(input string, exporter TextExporter, output string) {
	if input == "" || output == "" {
		flag.PrintDefaults()
		return
	}
	model, issues := LintFile(input)
	if model == nil {
		listIssues(input, issues)
		return
	}
	err := Export(*model, exporter, output)
	if err != nil {
		fmt.Println(err)
	}
}
