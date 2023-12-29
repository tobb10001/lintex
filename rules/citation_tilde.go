package rules

import (
	"strings"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
)

var citationTildeOnce sync.Once
var citationTilde *NativeRule

func CitationTilde() *NativeRule {
	citationTildeOnce.Do(func() {
		citationTilde = &NativeRule{
			id:          "vendored/citation_tilde",
			name:        "Citation Tilde",
			description: "A citation must be preceded by a word, that ends in a tilde to prevent a linebreak in between.",
			patterns: [][]byte{
				[]byte(`
					(text
					  word: (word) @word
					  .
					  ;; [Cc]iteauthor is also captured as citation.
					  word: (citation) @cite (#not-match? @cite "[Cc]iteauthor")
					)
				`),
			},
			apply: func(patternIndex int, query *sitter.Query, match *sitter.QueryMatch, input []byte) (*Range, error) {
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
