// Package for making the findings readable.
package output

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"lintex/rules"

	"github.com/rs/zerolog/log"
)

// Output the violation of a rule in a human readable format to stdout.
func PrintRuleViolation(violation *rules.Violation) {
	lines := bytes.Split(violation.Source, []byte("\n"))

	wd, err := os.Getwd()
	if err != nil {
		panic("Cannot get working directory???")
	}
	relPath, err := filepath.Rel(wd, violation.File)
	if err != nil {
		log.Error().Str("abspath", violation.File).Msg("Couldn't convert abspath to relative.")
		relPath = violation.File
	}

	fmt.Printf("%s:%d:%d\n",
		relPath,
		violation.Range.Start.Row+1,
		violation.Range.Start.Column,
	)
	fmt.Println(violation.Rule.Name())
	printSection(lines, violation.Range)
	fmt.Println(violation.Rule.Description())
	fmt.Println("")
}

// Print a subset of the input.
//
// When printing, one line above and one line below the section in question is also
// printed to provide context.
func printSection(lines [][]byte, rang *rules.Range) {
	if rang.Start.Row != 0 {
		fmt.Println(string(lines[rang.Start.Row-1][:]))
	}
	for line := rang.Start.Row; line <= rang.End.Row; line++ {
		fmt.Println(string(lines[line][:]))
	}
	if rang.End.Row < uint32(len(lines)) {
		fmt.Println(string(lines[rang.End.Row+1][:]))
	}
}
