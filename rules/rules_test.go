package rules_test

import (
	"regexp"
	"testing"

	"lintex/rules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNativeRuleID(t *testing.T) {
	ruls := rules.GetNativeRules()

	for _, rule := range ruls {
		matched, err := regexp.MatchString("vendored/([a-z_]+/)*([a-z_]+)", rule.ID())
		require.NoError(t, err)
		assert.Truef(t, matched, "Rule ID doesn't match the RegEx: %#v", rule.ID())
	}
}
