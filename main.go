package main

import (
	"context"
	"fmt"
	"os"

	"lintex/reader"
	"lintex/rules"
	"lintex/tslatex"

	sitter "github.com/smacker/go-tree-sitter"
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

	parser := sitter.NewParser()
	lang := tslatex.GetLanguage()
	parser.SetLanguage(lang)

	tree, err := sitter.ParseCtx(context.Background(), source, lang)
	if err != nil {
		panic(err)
	}

	cursor := sitter.NewQueryCursor()
	defer cursor.Close()

	for _, rule := range rules.GetRules() {

		query, err := sitter.NewQuery(rule.Pattern, lang)
		if err != nil {
			panic(err)
		}
		cursor.Exec(query, tree)

		for {
			match, ok := cursor.NextMatch()
			if !ok {
				break
			}
			match = cursor.FilterPredicates(match, source)
			for _, c := range match.Captures {
				if query.CaptureNameForId(c.Index) != rule.Predicate { continue }
				fmt.Println(rule.Name)
				fmt.Println(c.Node.Content(source))
				fmt.Println(rule.Description)
				fmt.Println("")
			}
		}
	}
}
