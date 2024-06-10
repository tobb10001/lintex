package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"lintex/files"
	"lintex/rules"
)

func TestCitationTilde(t *testing.T) {
	rule := rules.CitationTilde
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
			name:  "No tilde, no space",
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
			name:  "citeauthor",
			input: []byte(`Allowed with \citeauthor{myself}.`),
		},
		{
			name:  "Citeauthor",
			input: []byte(`Allowed with \Citeauthor{myself}.`),
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
