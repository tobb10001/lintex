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

type ApplyRuleFunc func(*sitter.Query, *sitter.QueryMatch, []byte) (*Range, error)

type Rule interface {
	Apply() ApplyRuleFunc
	Description() string
	Name() string
	Pattern() []byte
}

type NativeRule struct {
	name        string
	description string
	pattern     []byte
	apply       ApplyRuleFunc
}

func (nr NativeRule) Apply() ApplyRuleFunc { return nr.apply }
func (nr NativeRule) Description() string  { return nr.description }
func (nr NativeRule) Name() string         { return nr.name }
func (nr NativeRule) Pattern() []byte      { return nr.pattern }

type Violation struct {
	Rule   Rule
	Range  *Range
	Source []byte
}
