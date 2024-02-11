package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lintex/files"
	"lintex/rules"
	"lintex/tslatex"
)

func TestSpellingRule(t *testing.T) {
	rule := rules.NewSpellingRule("virtualization", "virtuali[sz]ation")
	tests := []struct {
		name  string
		input []byte
		error bool
	}{
		{
			name:  "correct spelling",
			input: []byte("virtualization"),
			error: false,
		},
		{
			name:  "incorrect spelling",
			input: []byte("virtualisation"),
			error: true,
		},
		{
			name:  "something else",
			input: []byte("something else"),
			error: false,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			tree, err := tslatex.GetTree(testcase.input)
			assert.NoError(t, err)
			violations, err := rules.ApplyRule(
				files.File{Path: "testfile", Tree: tree, Source: testcase.input},
				rule,
			)
			assert.NoError(t, err)
			if testcase.error {
				assert.Equal(t, 1, len(violations))
			} else {
				assert.Equal(t, 0, len(violations))
			}
		})
	}
}
