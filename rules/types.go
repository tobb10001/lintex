package rules

import sitter "github.com/smacker/go-tree-sitter"

type ApplyRuleError struct {
	message string
}

func (are ApplyRuleError) Error() string {
	return are.message
}

// A section in a file.
// Meant to hold the section of a file, that violates some rule.
type Range struct {
	Start sitter.Point
	End   sitter.Point
}

// A rule.
type Rule interface {
	// Apply the rule to a given candidate match.
	//
	// Rules, that can't be represented by a Tree-sitter query entirely can run
	// additional logic to determine, if the given match actually violates the rule.
	// Thereby, the Tree-sitter query determines cases, that potentially violate a
	// rule.
	//
	// In any case, this method is responsible for determining the range, that should
	// be highlighted.
	Apply(*sitter.Query, *sitter.QueryMatch, []byte) (*Range, error)
	// Get the description of the rule.
	Description() string
	// Get the name of the rule.
	Name() string
	// Get the Tree-sitter pattern to determine violation candidates.
	Pattern() []byte
}

// A rule implemented completely in Go.
type NativeRule struct {
	name        string
	description string
	pattern     []byte
	apply       func (*sitter.Query, *sitter.QueryMatch, []byte) (*Range, error)
}

func (nr NativeRule) Apply(query *sitter.Query, match *sitter.QueryMatch, source []byte) (*Range, error) {
	return nr.apply(query, match, source)
}
func (nr NativeRule) Description() string  { return nr.description }
func (nr NativeRule) Name() string         { return nr.name }
func (nr NativeRule) Pattern() []byte      { return nr.pattern }

// A violation to a rule.
//
// This is the intermediate product produced by the linter.
// It enables to later produce different kinds of output.
type Violation struct {
	// Reference to the rule, that was violated.
	Rule   Rule
	// Range, that shoudl be associated with the violation.
	Range  *Range
	// Source code of the file, to track back the original code if needed.
	Source []byte
}
