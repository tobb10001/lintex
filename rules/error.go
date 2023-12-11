package rules

import (
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
)

var errorOnce sync.Once
var erro *NativeRule

func Error() *NativeRule {
	errorOnce.Do(func() {
		erro = &NativeRule{
			name: "Error",
			description: "A parsing error by Tree Sitter, i.e. an error in the LaTeX source.",
			pattern: []byte(`(ERROR) @error`),
			apply: func(query *sitter.Query, match *sitter.QueryMatch, _ []byte) (*Range, error) {
				return &Range{
					Start: match.Captures[0].Node.StartPoint(),
					End: match.Captures[0].Node.EndPoint(),
				}, nil
			},
		}
	})
	return erro
}
