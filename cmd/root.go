package cmd

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "runtime",
	Short: "GitHub Runtime",
	Long: heredoc.Doc(`
		Use the GitHub Runtime CLI to deploy and manage apps on GitHub Runtime
	`),
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

type exitCode int

const (
	exitOK      exitCode = 0
	exitError   exitCode = 1
	exitCancel  exitCode = 2
	exitAuth    exitCode = 4
	exitPending exitCode = 8
)

func Execute() exitCode {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}

	return exitOK
}
