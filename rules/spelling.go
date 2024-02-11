package rules

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/spf13/viper"
)

type SpellcheckDefinition struct {
	Correct string
	Regex   string
}

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

func GetSpelling() ([]SpellingRule, error) {
	log.Trace().Any("viper", viper.AllSettings()).Send()
	definitionsAny := viper.Get("spellchecks")
	if definitionsAny == nil {
		// There are simply no rules specified.
		return []SpellingRule{}, nil
	}
	definitions, ok := definitionsAny.([]any)
	if !ok {
		return nil, fmt.Errorf("Couldn't read spellcheck definitions: Not an array: %+v (%T)", definitionsAny, definitions)
	}
	var res []SpellingRule
	for i, definitionAny := range definitions {
		var definition SpellcheckDefinition
		err := mapstructure.Decode(definitionAny, &definition)
		if err != nil {
			log.Warn().
				Err(err).
				Int("position", i).
				Msg("Couldn't read spellcheck definition.")
		}
		res = append(res, SpellingRule{correct: definition.Correct, regex: definition.Regex})
	}
	log.Debug().Int("len", len(res)).Msg("Created spelling rules.")
	return res, nil
}
