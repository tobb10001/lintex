// Lintex: A linter for LaTeX, powered by Tree-sitter.
package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"lintex/output"
	"lintex/reader"
	"lintex/rules"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}

	log.Debug().Msgf("Reading Document %s", os.Args[1])

	files, notFound, err := reader.ReadDocument(os.Args[1])
	if err != nil {
		panic(err)
	}

	if len(notFound) != 0 {
		arr := zerolog.Arr()
		for _, file := range notFound {
			arr.Str(file)	
		}
		log.Warn().Array("files", arr).Msg("Some included files aren't on disk.")
	}

	log.Debug().Msgf("Found %d files in total.", len(files))

	var violations []rules.Violation

	for _, file := range files {
		for _, rule := range rules.GetRules() {

			ranges, err := rules.ApplyRule(file.Tree, file.Source, rule)
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
