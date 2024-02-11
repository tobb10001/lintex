package cmd

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"lintex/files"
	"lintex/output"
	"lintex/rules"
)

func lint(*cobra.Command, []string) {
	rulez, err := rules.GetRules()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting the rules.")
	}

	log.Debug().Int("len", len(rulez)).Msg("Found rules.")

	files, err := files.GetFiles()
	if err != nil {
		log.Fatal().Err(err).Msg("Error finding files.")
	}

	_logArr := zerolog.Arr()
	for _, file := range files {
		_logArr.Str(file.Path)
	}
	log.Debug().Int("len", len(files)).Array("files", _logArr).Msg("Found files.")

	var violations []rules.Violation
	for _, file := range files {
		for _, rule := range rulez {
			log.Trace().Str("file", file.Path).Str("rule", rule.Name()).Msg("Applying rule to file.")

			ranges, err := rules.ApplyRule(file, rule)
			if err != nil {
				log.Warn().Err(err).Str("name", rule.Name()).Msg("Error applying a rule.")
			}

			log.Trace().Int("len", len(ranges)).Str("rule", rule.Name()).Str("file", file.Path).Msg("Found ranges.")

			for _, rang := range ranges {
				violations = append(violations, rules.Violation{
					File: file.Path, Rule: rule, Range: rang, Source: file.Source,
				})
			}
		}
	}

	log.Debug().Int("len", len(violations)).Msg("Before printing violations.")

	for _, violation := range violations {
		err := output.PrintRuleViolation(&violation)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to print a violation.")
		}
	}
}
