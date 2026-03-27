package cmd

import (
	"fmt"
	"net/url"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

type deleteCmdFlags struct {
	app string
	revisionName string
}

func init() {
	deleteCmdFlags := deleteCmdFlags{}
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a GitHub Runtime app",
		Long: heredoc.Doc(`
			Delete a GitHub Runtime app
		`),
		Example: heredoc.Doc(`
			$ gh runtime delete --app my-app
			# => Deletes the app named 'my-app'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.DefaultRESTClient()
			if err != nil {
				return fmt.Errorf("failed creating REST client: %v", err)
			}

			response, err := runDelete(client, deleteCmdFlags)
			if err != nil {
				return err
			}

			fmt.Printf("App deleted: %s\n", response)
			return nil
		},
	}

	deleteCmd.Flags().StringVarP(&deleteCmdFlags.app, "app", "a", "", "The app to delete")
	deleteCmd.Flags().StringVarP(&deleteCmdFlags.revisionName, "revision-name", "r", "", "The revision name to use for the app")
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(client restClient, flags deleteCmdFlags) (string, error) {
	if flags.app == "" {
		return "", fmt.Errorf("--app flag is required")
	}

	deleteUrl := fmt.Sprintf("runtime/%s/deployment", flags.app)
	params := url.Values{}
	if flags.revisionName != "" {
		params.Add("revision_name", flags.revisionName)
	}
	if len(params) > 0 {
		deleteUrl += "?" + params.Encode()
	}

	var response string
	err := client.Delete(deleteUrl, &response)
	if err != nil {
		return response, fmt.Errorf("error deleting app: %v", err)
	}

	// Actual response on success is empty body so return the ID
	return flags.app, nil
}
