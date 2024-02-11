// This module provides functionality to define rules as TOML files instead of Go code.
package rules

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"lintex/files"
	"lintex/tslatex"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
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

func TomlGetLocal(directory string) ([]TomlRule, error) {
	filesystem := os.DirFS(directory)
	return TomlRulesFromFS(filesystem, "local/")
}

func TomlGetVendored() ([]TomlRule, error) {
	filesystem, err := fs.Sub(vendored_toml, "toml")
	if err != nil {
		return nil, err
	}
	return TomlRulesFromFS(filesystem, "vendored/")
}

func checkInput(input []byte, rule TomlRule, expectedViolations int) error {
	tree, err := tslatex.GetTree(input)
	if err != nil {
		return err
	}
	violations, err := ApplyRule(files.File{Path: "testfile", Tree: tree, Source: input}, rule)
	if err != nil {
		return err
	}
	if expectedViolations != len(violations) {
		return fmt.Errorf(
			"Wrong number of violations: want=%d, got=%d",
			expectedViolations,
			len(violations),
		)
	}
	return nil
}

type TomlTestError struct {
	Rule     *TomlRule
	Location string
	Err      error
}

func TestTomlRule(rule TomlRule) []TomlTestError {
	tests := rule.Tests()
	var errors []TomlTestError
	log.Info().
		Int("obediences", len(tests.Obediences)).
		Int("violations", len(tests.Violations)).
		Str("name", rule.Name()).
		Msg("Testing rule.")

	for i, obedience := range rule.Tests().Obediences {
		err := checkInput(obedience.Input, rule, 0)
		if err != nil {
			errors = append(
				errors,
				TomlTestError{
					Rule:     &rule,
					Location: fmt.Sprintf("obedience #%d", i),
					Err:      err,
				},
			)
		}
	}
	for i, violation := range rule.Tests().Violations {
		err := checkInput(violation.Input, rule, 1)
		if err != nil {
			errors = append(
				errors,
				TomlTestError{
					Rule:     &rule,
					Location: fmt.Sprintf("violation #%d", i),
					Err:      err,
				},
			)
		}
	}
	return errors
}
