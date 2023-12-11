package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"lintex/files"
	"lintex/rules"
	"lintex/tslatex"
)

func TestError(t *testing.T) {
	rule := rules.Error()
	tests := []struct {
		name string
		input []byte
		error bool
	}{
		{
			name: "Input empty brackets",
			input: []byte(`\input{}`),
			error: true,
		},
		{
			name: "Input empty brackets",
			input: []byte(`\input{foo}`),
			error: false,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			tree, err := tslatex.GetTree(testcase.input)
			require.NoError(t, err)
			violations, err := rules.ApplyRule(
				files.File{Path: "testfile", Tree: tree, Source: testcase.input},
				rule,
			)
			require.NoError(t, err)
			if testcase.error {
				assert.Equal(t, 1, len(violations))
			} else {
				assert.Equal(t, 0, len(violations))
			}
		})
	}
}
