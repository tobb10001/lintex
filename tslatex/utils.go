package tslatex

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
)

func GetCaptures(tree *sitter.Node, pattern []byte, predicate string, source []byte) ([]sitter.QueryCapture, error) {
	cursor := sitter.NewQueryCursor()
	defer cursor.Close()

	query, err := sitter.NewQuery(pattern, GetLanguage())
	if err != nil {
		return nil, err
	}
	cursor.Exec(query, tree)

	var captures []sitter.QueryCapture

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		match = cursor.FilterPredicates(match, source)
		for _, capture := range match.Captures {
			if query.CaptureNameForId(capture.Index) == predicate {
				captures = append(captures, capture)
			}
		}
	}

	return captures, nil
}

func GetTree(source []byte) (*sitter.Node, error) {

	parser := sitter.NewParser()
	lang := GetLanguage()
	parser.SetLanguage(lang)

	return sitter.ParseCtx(context.Background(), source, lang)
}
