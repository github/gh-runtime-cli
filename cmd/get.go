package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

type getCmdFlags struct {
	app string
}

type serverResponse struct {
	AppUrl string `json:"app_url"`
}

func init() {
	getCmdFlags := getCmdFlags{}
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get details of a GitHub Runtime app",
		Long: heredoc.Doc(`
			Get details of a GitHub Runtime app
		`),
		Example: heredoc.Doc(`
			$ gh runtime get --app my-app
			# => Retrieves details of the app named 'my-app'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if getCmdFlags.app == "" {
				return fmt.Errorf("--app flag is required")
			}

			getUrl := fmt.Sprintf("runtime/%s/deployment", getCmdFlags.app)
			client, err := api.DefaultRESTClient()
			if err != nil {
				return fmt.Errorf("failed creating REST client: %v", err)
			}

			response := serverResponse{}
			err = client.Get(getUrl, &response)
			if err != nil {
				return fmt.Errorf("retrieving app details: %v", err)
			}

			fmt.Printf("%s\n", response.AppUrl)
			return nil
		},
	}

	getCmd.Flags().StringVarP(&getCmdFlags.app, "app", "a", "", "The app to retrieve details for")
	rootCmd.AddCommand(getCmd)
}
