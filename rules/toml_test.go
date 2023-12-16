// This file is for testing the *vendored* TOML rules.
package rules_test

import (
	"testing"

	"lintex/rules"

	"github.com/stretchr/testify/require"
)

func TestTomlRules(t *testing.T) {
	ruls, err := rules.TomlGetVendored()
	require.NoError(t, err)

	for _, rule := range ruls {
		t.Run(rule.Name(), func(t *testing.T) { rules.TestTomlRule(t, rule) })
	}
}
