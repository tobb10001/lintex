package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"lintex/files"
	"lintex/rules"
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
			file, err := files.NewFile("testfile", testcase.input)
			require.NoError(t, err)
			violations, err := rules.ApplyRule(file, rule)
			assert.NoError(t, err)
			if testcase.error {
				assert.Equal(t, 1, len(violations))
			} else {
				assert.Equal(t, 0, len(violations))
			}
		})
	}
}
