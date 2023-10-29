package rules

import (
	sitter "github.com/smacker/go-tree-sitter"
)

func CaptionTrailingPeriod() Rule {
	return Rule{
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
					return &Range{Start: capture.Node.StartPoint(), End: capture.Node.EndPoint()}, nil
				}
			}
			return nil, ApplyRuleError{"Could not find a capture for the `@caption` predicate..."}
		},
	}
}
