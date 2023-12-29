// Package to hold all logic considering rules.
//
// This includes some rules themselves.
package rules

import (
	"lintex/files"
	"lintex/tslatex"

	"github.com/rs/zerolog/log"
)

// Apply a rule to a given syntax tree.
//
// In order to filter predicates, this method also needs access to the source.
// It returns the ranges, that violate the rule. It might return an empty slice, if
// there are no violations to the given rule.
func ApplyRule(file files.File, rule Rule) ([]*Range, error) {
	var violations []*Range
	for i, pattern := range rule.Patterns() {
		query, matches, err := tslatex.GetMatches(file.Tree, pattern, file.Source)
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			rang, err := rule.Apply(i, query, match, file.Source)
			if err != nil {
				panic(err)
			}
			if rang != nil {
				violations = append(violations, rang)
			}
		}
	}

	return violations, nil

}

func GetNativeRules() ([]Rule) {
	return []Rule{
		CitationTilde(),
	}
}

// Optain a list of all configured rules.
func GetRules() ([]Rule, error) {
	log.Debug().Msg("Getting rules...")
	toml_vendored, err := TomlGetVendored()
	if err != nil {
		return nil, err
	}
	var rules []Rule
	// rules = append(rules, toml_vendored...)
	for _, rule := range toml_vendored {
		rules = append(rules, rule)
	}
	rules = append(rules, GetNativeRules()...)
	return rules, nil
}
