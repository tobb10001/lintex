package rules

import (
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
)

var captionTrailingPeriodOnce sync.Once
var captionTrailingPeriod *Rule

func CaptionTrailingPeriod() *Rule {
	captionTrailingPeriodOnce.Do(func() {
		captionTrailingPeriod = &Rule{
			name:        "Caption Trailing Period",
			description: "A caption should not have a trailing period, because it would end up in the ToX as well.",
			pattern: []byte(`
				(caption
				  long: (curly_group
					(text 
					  (word) @last_word (#match? @last_word "\\.$")
					  .
					)
				  )
				) @caption
			`),
			apply: func(query *sitter.Query, match *sitter.QueryMatch, _ []byte) (*Range, error) {
				for _, capture := range match.Captures {
					if query.CaptureNameForId(capture.Index) == "caption" {
						return &Range{Start: capture.Node.StartPoint(), End: capture.Node.EndPoint()}, nil
					}
				}
				return nil, ApplyRuleError{"Could not find a capture for the `@caption` predicate..."}
			},
		}
	})
	return captionTrailingPeriod
}
