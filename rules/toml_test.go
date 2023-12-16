// This file is for testing the *vendored* TOML rules.
package rules_test

import (
	"testing"

	"lintex/files"
	"lintex/rules"
	"lintex/tslatex"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTomlRules(t *testing.T) {
	ruls, err := rules.TomlGetVendored()
	require.NoError(t, err)

	for _, rule := range ruls {
		t.Run(rule.Name(), func(t *testing.T) {
			for _, obedience := range rule.Tests().Obediences {
				t.Run(obedience.Name, func(t *testing.T){
					tree, err := tslatex.GetTree(obedience.Input)
					require.NoError(t, err)
					violations, err := rules.ApplyRule(
						files.File{Path: "testfile", Tree: tree, Source: obedience.Input},
						rule,
					)
					assert.Equal(t, 0, len(violations), "Expected to see no violation.")
				})
			}
			for _, violation := range rule.Tests().Violations {
				t.Run(violation.Name, func(t *testing.T) {
					tree, err := tslatex.GetTree(violation.Input)
					require.NoError(t, err)
					violations, err := rules.ApplyRule(
						files.File{Path: "testfile", Tree: tree, Source: violation.Input},
						rule,
					)
					assert.Equal(t, 1, len(violations), "Expected to see one violation.")
				})
			}
		})
	}
}
