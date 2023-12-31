package cmd

import (
	"fmt"
	"lintex/rules"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests for local TOML rules.",
	Long:  `This command runs tests for local rules defined in the TOML format.`,
	RunE:  test,
}

func test(*cobra.Command, []string) error {
	ruls, err := rules.TomlGetLocal(".lintex/rules")
	if err != nil {
		return err
	}
	hasErr := false
	colorId := color.New(color.FgGreen)
	colorErr := color.New(color.FgRed)
	for _, rule := range ruls {
		fmt.Print("Checking rule ")
		colorId.Print(rule.ID())
		fmt.Printf(": %s\n", rule.Name())
		errors := rules.TestTomlRule(rule)
		for _, err := range errors {
			hasErr = true
			colorErr.Println(fmt.Sprintf("Error at %s: %s", err.Location, err.Err))
		}
	}
	if hasErr {
		colorErr.Println("Errors in local rules detected.")
		os.Exit(1)
		return nil
	} else {
		color.Green("All checks passed.")
		return nil
	}
}

func init() {
	rootCmd.AddCommand(testCmd)
}
