package rules

import sitter "github.com/smacker/go-tree-sitter"

type ApplyRuleError struct {
	message string
}

func (are ApplyRuleError) Error() string {
	return are.message
}

type Range struct {
	Start sitter.Point
	End   sitter.Point
}

type Rule struct {
	Name        string
	Description string
	Pattern     []byte
	Apply       func(*sitter.Query, *sitter.QueryMatch, []byte) (*Range, error)
}
