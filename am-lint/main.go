package main

import (
	"fmt"
)

func main() {
	definition := `system: foo
`
	model, issues := LintText(definition)

	fmt.Println("Issues:")
	for issue := range issues {
		fmt.Println(issue)
	}
	fmt.Println("Model:")
	fmt.Println(model)
}
