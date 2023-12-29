// Functionality related to the process of parsing TOML files into TomlRules.
package rules

import (
	"io/fs"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
)

func parseRuleFS(filesystem fs.FS, path, prefix string) (*TomlRule, error) {
	var rule TomlRuleParse
	meta, err := toml.DecodeFS(filesystem, path, &rule)
	if err != nil {
		return nil, err
	}

	undecoded := meta.Undecoded()
	if len(undecoded) > 0 {
		var keys []string
		for _, key := range undecoded {
			keys = append(keys, key.String())
		}
		log.Warn().
			Type("fstype", filesystem).
			Str("file", path).
			Strs("keys", keys).
			Msg("Found unedecoded key in TOML file.")
	}
	id, _ := strings.CutSuffix(prefix+path, ".toml")
	return TomlRuleFromParse(&rule, id), nil
}

type TomlRuleParse struct {
	Name        string
	Description string
	Patterns    []string
	Capture     string
	Tests       TomlRuleTestsParse
}

func TomlRuleFromParse(trp *TomlRuleParse, id string) *TomlRule {
	var patterns [][]byte
	for _, pattern := range trp.Patterns {
		patterns = append(patterns, []byte(pattern))
	}
	return &TomlRule{
		id:          id,
		name:        trp.Name,
		description: trp.Description,
		patterns:    patterns,
		capture:     trp.Capture,
		tests:       *TomlRuleTestsFromParse(&trp.Tests),
	}
}

type TomlRuleTestsParse struct {
	Obediences []TomlRuleTestCaseParse
	Violations []TomlRuleTestCaseParse
}

func TomlRuleTestsFromParse(trtp *TomlRuleTestsParse) *TomlRuleTests {
	var obediences, violations []TomlRuleTestCase
	for _, obedience := range trtp.Obediences {
		obediences = append(obediences, *TomlRuleTestCaseFromParse(&obedience))
	}
	for _, violation := range trtp.Violations {
		violations = append(violations, *TomlRuleTestCaseFromParse(&violation))
	}
	return &TomlRuleTests{Obediences: obediences, Violations: violations}
}

type TomlRuleTestCaseParse struct {
	Name  string
	Input string
}

func TomlRuleTestCaseFromParse(trtcp *TomlRuleTestCaseParse) *TomlRuleTestCase {
	return &TomlRuleTestCase{Name: trtcp.Name, Input: []byte(trtcp.Input)}
}
