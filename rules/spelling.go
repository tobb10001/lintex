package rules

import (
	"fmt"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
)

type SpellingRule struct {
	correct string
	regex   string
}

func (sr SpellingRule) ID() string   { return "spelling/" + sr.correct }
func (sr SpellingRule) Name() string { return "Spelling: " + sr.correct }

func (sr SpellingRule) Apply(patternIndex int, query *sitter.Query, match *sitter.QueryMatch, _ []byte) (*Range, error) {
	log.Trace().Str("rule", sr.Name()).Msg("Applying spelling rule.")
	for _, capture := range match.Captures {
		if query.CaptureNameForId(capture.Index) == "word" {
			return &Range{Start: capture.Node.StartPoint(), End: capture.Node.EndPoint()}, nil
		}
	}
	return nil, ApplyRuleError{"Could not find a capture for the `@word` predicate."}
}

func (sr SpellingRule) Description() string {
	return "Check that '" + sr.correct + "' is spelled correct."
}

func (sr SpellingRule) Patterns() [][]byte {
	return [][]byte{
		[]byte(
			fmt.Sprintf(`(
			  (word) @word
			  (#match? @word "%s")
			  (#not-eq? @word "%s")
			)`, sr.regex, sr.correct),
		),
	}
}

func NewSpellingRule(correct, regex string) SpellingRule {
	return SpellingRule{correct: correct, regex: regex}
}
