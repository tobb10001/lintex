package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lintex/files"
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
			name: "No tilde, no space",
			input: []byte(`No tilde or space\cite{sth}.`),
			error: true,
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
		{
			name: "citeauthor",
			input: []byte(`Allowed with \citeauthor{myself}.`),
		},
		{
			name: "Citeauthor",
			input: []byte(`Allowed with \Citeauthor{myself}.`),
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
