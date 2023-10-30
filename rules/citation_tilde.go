package rules

import (
	"strings"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
)

var citationTildeOnce sync.Once
var citationTilde *Rule

func CitationTilde() *Rule {
	citationTildeOnce.Do(func() {
		citationTilde = &Rule{
			name:        "Citation Tilde",
			description: "A citation must be preceded by a word, that ends in a tilde to prevent a linebreak in between.",
			pattern: []byte(`
				(text
				  word: (word) @word 
				  .
				  word: (citation) @cite
				)
			`),
			apply: func(query *sitter.Query, match *sitter.QueryMatch, input []byte) (*Range, error) {
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
				if word.Node.EndPoint().Row < cite.Node.StartPoint().Row || cite.Node.StartPoint().Column-word.Node.EndPoint().Column != 0 {
					return &Range{word.Node.StartPoint(), cite.Node.EndPoint()}, nil
				}
				return nil, nil
			},
		}
	})
	return citationTilde
}
