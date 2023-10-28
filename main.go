package main

import (
	"fmt"
	"os"

	"lintex/reader"
	"lintex/rules"
	"lintex/tslatex"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}
	source, err := reader.Read(os.Args[1])
	if err != nil {
		panic(err)
	}

	tree, err := tslatex.GetTree(source)
	if err != nil {
		panic(err)
	}

	for _, rule := range rules.GetRules() {

		captures, err := tslatex.GetCaptures(tree, rule.Pattern, rule.Predicate, source)
		if err != nil {
			panic(err)
		}

		for _, capture := range captures {
			fmt.Println(rule.Name)
			fmt.Println(capture.Node.Content(source))
			fmt.Println(rule.Description)
			fmt.Println("")
		}
	}
}
