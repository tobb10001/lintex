// Package to hold all logic considering rules.
//
// This includes some rules themselves.
package rules

import (
	sitter "github.com/smacker/go-tree-sitter"

	"lintex/tslatex"
)

// Apply a rule to a given syntax tree.
//
// In order to filter predicates, this method also needs access to the source.
// It returns the ranges, that violate the rule. It might return an empty slice, if
// there are no violations to the given rule.
func ApplyRule(tree *sitter.Node, source []byte, rule Rule) ([]*Range, error) {
	query, matches, err := tslatex.GetMatches(tree, rule.Pattern(), source)
	if err != nil {
		return nil, err
	}

	var violations []*Range

	for _, match := range matches {
		rang, err := rule.Apply(query, match, source)
		if err != nil {
			panic(err)
		}
		if rang != nil {
			violations = append(violations, rang)
		}
	}

	return violations, nil

}

// Optain a list of all configured rules.
func GetRules() []Rule {
	return []Rule{
		CaptionTrailingPeriod(),
		CitationTilde(),
	}
}
