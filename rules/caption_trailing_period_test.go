package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lintex/rules"
	"lintex/tslatex"
)

func TestCaptionTrailingPeriod(t *testing.T) {
	rule := rules.CaptionTrailingPeriod()
	tests := []struct {
		name  string
		input []byte
		error bool
	}{
		{
			name:  "With trailing period",
			input: []byte(`\caption{This caption has a trailing period.}`),
			error: true,
		},
		{
			name:  "Without trailing period",
			input: []byte(`\caption{This caption has a trailing period}`),
			error: false,
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
