package cmd

import (
	"io/fs"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "lintex",
	Short: "A linter for LaTeX, powered by Tree-Sitter.",
	Long:  `To run LinTeX, simply run lintex from your projects root directory.`,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		err := os.Chdir(viper.GetString("workdir"))
		if err != nil {
			return err
		}
		err = viper.ReadInConfig()
		if err != nil {
			_, isConfigFileNotFound := err.(viper.ConfigFileNotFoundError)
			_, isPathError := err.(*fs.PathError)
			if !(isConfigFileNotFound || isPathError) {
				log.Fatal().
					AnErr("error", err).
					Msg("Couldn't read config file.")
			}
		}
		if viper.GetBool("debug") {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
		if viper.GetBool("trace") {
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		}
		return nil
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
	// Config file
	viper.SetConfigFile("./.lintex/config.toml")

	// Flags
	rootCmd.
		PersistentFlags().
		String("workdir", ".", "specify the working directory")
	_ = viper.BindPFlag("workdir", rootCmd.PersistentFlags().Lookup("workdir"))
	rootCmd.
		PersistentFlags().
		Bool("debug", false, "set log level to 'debug'")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	rootCmd.
		PersistentFlags().
		Bool("trace", false, "set log level to 'trace'")
	_ = viper.BindPFlag("trace", rootCmd.PersistentFlags().Lookup("trace"))
}
