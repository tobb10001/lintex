// This file is for testing the *vendored* TOML rules.
package rules_test

import (
	"regexp"
	"testing"

	"lintex/rules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTomlRules(t *testing.T) {
	ruls, err := rules.TomlGetVendored()
	require.NoError(t, err)

	for _, rule := range ruls {
		t.Run(rule.Name(), func(t *testing.T) { rules.TestTomlRule(t, rule) })
	}
}

func TestTomlRuleID(t *testing.T) {
	ruls, err := rules.TomlGetVendored()
	require.NoError(t, err)

	for _, rule := range ruls {
		matched, err := regexp.MatchString("vendored/([a-z_]+/)*([a-z_]+)", rule.ID())
		require.NoError(t, err)
		assert.Truef(t, matched, "Rule ID doesn't match the RegEx: %#v", rule.ID())
	}
}
