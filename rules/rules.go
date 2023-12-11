// Package to hold all logic considering rules.
//
// This includes some rules themselves.
package rules

import (
	"lintex/files"
	"lintex/tslatex"
)

// Apply a rule to a given syntax tree.
//
// In order to filter predicates, this method also needs access to the source.
// It returns the ranges, that violate the rule. It might return an empty slice, if
// there are no violations to the given rule.
func ApplyRule(file files.File, rule Rule) ([]*Range, error) {
	query, matches, err := tslatex.GetMatches(file.Tree, rule.Pattern(), file.Source)
	if err != nil {
		return nil, err
	}

	var violations []*Range

	for _, match := range matches {
		rang, err := rule.Apply(query, match, file.Source)
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
