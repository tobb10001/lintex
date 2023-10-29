package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lintex/rules"
	"lintex/tslatex"
)

func TestCitationTilde(t *testing.T) {
	rule := rules.CitationTilde()
	tests := []struct {
		name  string
		input []byte
		error bool
	}{
		{
			name:  "Correct citation",
			input: []byte(`Correct~\cite{sth}.`),
			error: false,
		},
		{
			name:  "Additional space",
			input: []byte(`Correct~ \cite{sth}.`),
			error: true,
		},
		{
			name:  "Missing tilde",
			input: []byte(`Correct \cite{sth}.`),
			error: true,
		},
		{
			name: "Line break",
			input: []byte(`Correct
\cite{sth}.`),
			error: true,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			tree, err := tslatex.GetTree(testcase.input)
			assert.NoError(t, err)
			violations, err := rules.ApplyRule(tree, testcase.input, &rule)
			assert.NoError(t, err)
			if testcase.error {
				assert.Equal(t, 1, len(violations))
			} else {
				assert.Equal(t, 0, len(violations))
			}
		})
	}
}
