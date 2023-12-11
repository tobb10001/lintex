// Lintex: A linter for LaTeX, powered by Tree-sitter.
package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"lintex/files"
	"lintex/output"
	"lintex/rules"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	files, err := files.GetFiles()
	if err != nil {
		log.Fatal().Msg("Error finding files.")
	}

	var violations []rules.Violation
	for _, file := range files {
		for _, rule := range rules.GetRules() {

			ranges, err := rules.ApplyRule(file, rule)
			if err != nil {
				panic(err)
			}

			for _, rang := range ranges {
				violations = append(violations, rules.Violation{
					File: file.Path, Rule: rule, Range: rang, Source: file.Source,
				})
			}
		}
	}

	for _, violation := range violations {
		output.PrintRuleViolation(&violation)
	}
}
