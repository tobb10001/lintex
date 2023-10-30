package output

import (
	"fmt"

	"lintex/rules"
)

func PrintRuleViolation(violation *rules.Violation) {
	lines := getLines(violation.Source)

	fmt.Println(violation.Rule.Name())
	printSection(lines, violation.Range)
	fmt.Println(violation.Rule.Description())
	fmt.Println("")
}

func getLines(source []byte) [][]byte {
	var res [][]byte
	res = append(res, []byte(""))
	for _, byt := range source {
		if byt == '\n' {
			res = append(res, []byte(""))
		} else {
			res[len(res) - 1] = append(res[len(res) - 1], byt)
		}
	}
	return res
}

func printSection(lines [][]byte, rang *rules.Range) {
	if rang.Start.Row != 0 {
		fmt.Println(string(lines[rang.Start.Row - 1][:]))
	}
	for line := rang.Start.Row; line <= rang.End.Row; line++ {
		fmt.Println(string(lines[line][:]))
	}
	if rang.End.Row < uint32(len(lines)) {
		fmt.Println(string(lines[rang.End.Row + 1][:]))
	}
}
