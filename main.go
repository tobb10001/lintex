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
	if err != nil {
		panic(err)
	}
	cursor := sitter.NewQueryCursor()
	defer cursor.Close()
	cursor.Exec(query, tree)

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		match = cursor.FilterPredicates(match, source)
		for _, c := range match.Captures {
			if query.CaptureNameForId(c.Index) != "caption" { continue }
			fmt.Println(query.CaptureNameForId(c.Index))
			fmt.Println(c.Node.Content(source))
		}
	}
}
