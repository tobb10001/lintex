package main

import (
	"context"
	"fmt"

	"lintex/tslatex"

	sitter "github.com/smacker/go-tree-sitter"
)

func main() {
	source := []byte(`\caption{This is a caption with a period.}`)
	pattern := []byte(`
		(caption
		  long: (curly_group
			(text 
			  (word) @last_word (#match? @last_word "\\.$")
			  .
			)
		  )
		) @caption
	`)

	parser := sitter.NewParser()
	lang := tslatex.GetLanguage()
	parser.SetLanguage(lang)

	tree, err := sitter.ParseCtx(context.Background(), source, lang)
	if err != nil {
		panic(err)
	}

	query, err := sitter.NewQuery(pattern, lang)
	cursor := sitter.NewQueryCursor()
	defer cursor.Close()
	cursor.Exec(query, tree)

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		match = cursor.FilterPredicates(match, source)
		for idx, c := range match.Captures {
			if idx != 0 {
				continue
			}
			fmt.Println(c.Node.Content(source))
		}
	}
}
