package reader_test

import (
	"lintex/reader"
	"lintex/tslatex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIncludedFiles(t *testing.T) {

	tests := []struct{
		name string
		input []byte
		expected []string
	}{
		{
			name: "No include",
			input: []byte(`This file does not have an include\cite{nobody}.`),
			expected: []string{},
		},
		{
			name: "Include and input",
			input: []byte(`\input{file_to_input} and \include{file_to_include}`),
			expected: []string{"file_to_input.tex", "file_to_include.tex"},
		},
		{
			name: "Include with path",
			input: []byte(`\input{directory/file}`),
			expected: []string{"directory/file.tex"},
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func (t *testing.T) {
			tree, err := tslatex.GetTree(testcase.input)
			assert.NoError(t, err)
			actual, err := reader.GetIncludedFiles(tree, testcase.input)
			assert.NoError(t, err)
			assert.Equal(t, testcase.expected, actual)
		})
	}

}
