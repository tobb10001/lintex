package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var WorkDir string

var rootCmd = &cobra.Command{
	Use:   "lintex",
	Short: "A linter for LaTeX, powered by Tree-Sitter.",
	Long:  `To run LinTeX, simply run lintex from your projects root directory.`,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		return os.Chdir(WorkDir)
	},
	Run: lint,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.
		PersistentFlags().
		StringVar(&WorkDir, "workdir", "", "specify the working directory")
}
