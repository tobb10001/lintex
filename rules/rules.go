package rules

import (
	sitter "github.com/smacker/go-tree-sitter"

	"lintex/tslatex"
)

func ApplyRule(tree *sitter.Node, source []byte, rule Rule) ([]*Range, error) {
	query, matches, err := tslatex.GetMatches(tree, rule.Pattern(), source)
	if err != nil {
		return nil, err
	}

	var violations []*Range

	for _, match := range matches {
		rang, err := rule.Apply()(query, match, source)
		if err != nil {
			panic(err)
		}
		if rang != nil {
			violations = append(violations, rang)
		}
	}

	return violations, nil

}

func GetRules() []Rule {
	return []Rule{
		CaptionTrailingPeriod(),
		CitationTilde(),
	}
}
