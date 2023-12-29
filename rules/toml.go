// This module provides functionality to define rules as TOML files instead of Go code.
package rules

import (
	"embed"
	"fmt"
	"io/fs"
	"lintex/files"
	"lintex/tslatex"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed toml/*
var vendored_toml embed.FS

type TomlRule struct {
	id          string
	name        string
	description string
	patterns    [][]byte
	capture     string
	tests       TomlRuleTests
}

type TomlRuleTests struct {
	Obediences []TomlRuleTestCase
	Violations []TomlRuleTestCase
}

type TomlRuleTestCase struct {
	Name  string
	Input []byte
}

func (tr TomlRule) Apply(patternIndex int, query *sitter.Query, match *sitter.QueryMatch, _ []byte) (*Range, error) {
	log.Trace().Str("rule", tr.Name()).Msg("Applying TOML rule.")
	for _, capture := range match.Captures {
		if query.CaptureNameForId(capture.Index) == tr.capture {
			return &Range{Start: capture.Node.StartPoint(), End: capture.Node.EndPoint()}, nil
		}
	}
	return nil, ApplyRuleError{fmt.Sprintf("Could not find a capture for the `@%s` predicate.", tr.capture)}
}
func (tr TomlRule) ID() string          { return tr.id }
func (tr TomlRule) Description() string { return tr.description }
func (tr TomlRule) Name() string        { return tr.name }
func (tr TomlRule) Patterns() [][]byte  { return tr.patterns }

func (tr TomlRule) Tests() TomlRuleTests { return tr.tests }

func TomlRulesFromFS(filesystem fs.FS, prefix string) ([]TomlRule, error) {
	var files []string
	error := fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ".toml" {
			files = append(files, path)
		}
		return nil
	})
	if error != nil {
		return nil, error
	}

	var rules []TomlRule
	for _, path := range files {
		log.Trace().Str("file", path).Msg("Parsing TOML rule.")
		rule, err := parseRuleFS(filesystem, path, prefix)
		if err != nil {
			return nil, err
		}
		rules = append(rules, *rule)
	}
	log.Debug().Int("len", len(rules)).Msg("Found TOML rules.")
	return rules, nil
}

func TomlGetVendored() ([]TomlRule, error) {
	filesystem, err := fs.Sub(vendored_toml, "toml")
	if err != nil {
		return nil, err
	}
	return TomlRulesFromFS(filesystem, "vendored/")
}

func checkInput(t *testing.T, input []byte, rule TomlRule, expectedViolations int) bool {
	tree, err := tslatex.GetTree(input)
	if !assert.NoError(t, err) {
		return false
	}
	violations, err := ApplyRule(files.File{Path: "testfile", Tree: tree, Source: input}, rule)
	require.NoError(t, err)
	return assert.Equal(t, expectedViolations, len(violations))
}

func TestTomlRule(t *testing.T, rule TomlRule) bool {
	result := true
	t.Logf("Rule: %s, Obediences: %d, Violations: %d.", rule.Name(), len(rule.Tests().Obediences), len(rule.Tests().Violations))

	for _, obedience := range rule.Tests().Obediences {
		result = checkInput(t, obedience.Input, rule, 0) && result
	}
	for _, violation := range rule.Tests().Violations {
		result = checkInput(t, violation.Input, rule, 1) && result
	}
	return result
}
