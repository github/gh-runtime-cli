package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

// Version of the CLI app.
var Version = "0.1.0"

// The command prints out the version of the CLI app.
func init() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Get version of the CLI app",
		Long: heredoc.Doc(`
			Get version of the CLI app
		`),
		Example: heredoc.Doc(`
			$ gh runtime version
			# => Retrieves version of the CLI app
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%s\n", Version)

			return nil
		},
	}

	rootCmd.AddCommand(versionCmd)
}
