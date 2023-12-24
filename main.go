// Lintex: A linter for LaTeX, powered by Tree-sitter.
package main

import (
	"lintex/cmd"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	cmd.Execute()
}
