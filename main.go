package main

import (
	"fmt"
	"os"

	"lintex/output"
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

		violations, err := rules.ApplyRule(tree, source, &rule)
		if err != nil {
			panic(err)
		}

		for _, violation := range violations {
			output.PrintRuleViolation(&rule, violation, source)
		}
	}
}
