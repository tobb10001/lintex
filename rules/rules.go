package rules

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

type Range struct {
	start sitter.Point
	end   sitter.Point
}

type Rule struct {
	Name        string
	Description string
	Pattern     []byte
	Apply       func(*sitter.Query, *sitter.QueryMatch, []byte) (*Range, error)
}

type ApplyRuleError struct {
	message string
}

func (are ApplyRuleError) Error() string {
	return are.message
}

func GetRules() []Rule {
	return []Rule{
		{
			Name:        "Caption Trailing Period",
			Description: "A caption should not have a trailing period, because it would end up in the ToX as well.",
			Pattern: []byte(`
				(caption
				  long: (curly_group
					(text 
					  (word) @last_word (#match? @last_word "\\.$")
					  .
					)
				  )
				) @caption
			`),
			Apply: func(query *sitter.Query, match *sitter.QueryMatch, _ []byte) (*Range, error) {
				for _, capture := range match.Captures {
					if query.CaptureNameForId(capture.Index) == "caption" {
						return &Range{start: capture.Node.StartPoint(), end: capture.Node.EndPoint()}, nil
					}
				}
				return nil, ApplyRuleError{"Could not find a capture for the `@caption` predicate..."}
			},
		},
		{
			Name:        "Citation After Tilde",
			Description: "A citation must be preceded by a word, that ends in a tilde to prevent a linebreak in between.",
			Pattern: []byte(`
				(text
				  word: (word) @word 
				  .
				  word: (citation) @cite
				)
			`),
			Apply: func(query *sitter.Query, match *sitter.QueryMatch, input []byte) (*Range, error) {
				var word, cite sitter.QueryCapture
				for _, capture := range match.Captures {
					capture_name := query.CaptureNameForId(capture.Index)
					if capture_name == "word" {
						word = capture
					} else if capture_name == "cite" {
						cite = capture
					}
				}
				if !strings.HasSuffix(word.Node.Content(input), "~") {
					return &Range{word.Node.StartPoint(), cite.Node.EndPoint()}, nil
				}
				if word.Node.EndPoint().Row < cite.Node.StartPoint().Row || cite.Node.StartPoint().Column-word.Node.EndPoint().Column != 1 {
					return &Range{word.Node.StartPoint(), cite.Node.EndPoint()}, nil
				}
				return nil, nil
			},
		},
	}
}
