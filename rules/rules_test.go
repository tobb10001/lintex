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
	regex, err := regexp.Compile("vendored/([a-z_]+/)*([a-z_]+)")
	require.NoError(t, err)

	for _, rule := range ruls {
		matched := regex.MatchString(rule.ID())
		assert.Truef(t, matched, "Rule ID doesn't match the RegEx: %#v", rule.ID())
	}
}
