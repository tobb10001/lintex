package tslatex

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
)

func GetMatches(tree *sitter.Node, pattern []byte, source []byte) (*sitter.Query, []*sitter.QueryMatch, error) {
	cursor := sitter.NewQueryCursor()
	defer cursor.Close()

	query, err := sitter.NewQuery(pattern, GetLanguage())
	if err != nil {
		return nil, nil, err
	}
	cursor.Exec(query, tree)

	var matches []*sitter.QueryMatch

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		match = cursor.FilterPredicates(match, source)
		if len(match.Captures) > 0 {
			matches = append(matches, match)
		}
	}

	return query, matches, nil
}

func GetTree(source []byte) (*sitter.Node, error) {

	parser := sitter.NewParser()
	lang := GetLanguage()
	parser.SetLanguage(lang)

	return sitter.ParseCtx(context.Background(), source, lang)
}
