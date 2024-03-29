// Package for making the findings readable.
package output

import (
	"bytes"
	"fmt"

	"lintex/rules"

	"github.com/fatih/color"
)

// Output the violation of a rule in a human readable format to stdout.
func PrintRuleViolation(violation *rules.Violation) error {
	lines := bytes.Split(violation.Source, []byte("\n"))

	fmt.Printf("%s:%d:%d\n",
		violation.File,
		violation.Range.Start.Row+1,
		violation.Range.Start.Column,
	)
	id_color := color.New(color.FgGreen)
	id_color.Print(violation.Rule.ID())
	fmt.Printf(": %s\n", violation.Rule.Name())
	printSection(lines, violation.Range)
	fmt.Println(violation.Rule.Description())
	fmt.Println("")
	return nil
}

// Print a subset of the input.
//
// When printing, one line above and one line below the section in question is also
// printed to provide context.
func printSection(lines [][]byte, rang *rules.Range) {
	// The line before the match.
	if rang.Start.Row != 0 {
		fmt.Println(fmt.Sprintf("%4d: ", rang.Start.Row) + string(lines[rang.Start.Row-1][:]))
	}

	// Define the highlight color
	h := color.New(color.FgYellow)

	// The lines that contain the match.
	fmt.Printf("%4d: %s", rang.Start.Row+1, string(lines[rang.Start.Row][:rang.Start.Column]))
	if rang.Start.Row == rang.End.Row {
		h.Print(string(lines[rang.Start.Row][rang.Start.Column:rang.End.Column]))
	} else {
		h.Print(string(lines[rang.Start.Row][rang.Start.Column:]) + "\n")
		for line := rang.Start.Row + 1; line < rang.End.Row; line++ {
			fmt.Printf("%4d: ", line+1)
			h.Print(string(lines[line]) + "\n")
		}
		fmt.Printf("%4d: ", rang.End.Row+1)
		h.Print(string(lines[rang.End.Row][:rang.End.Column]))
	}
	fmt.Print(string(lines[rang.End.Row][rang.End.Column:]) + "\n")

	// The line after the match.
	if rang.End.Row < uint32(len(lines)) {
		fmt.Println(fmt.Sprintf("%4d: ", rang.End.Row+2) + string(lines[rang.End.Row+1][:]))
	}
}
