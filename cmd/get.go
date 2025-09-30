package cmd

import (
	"fmt"
	"net/url"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/github/gh-runtime-cli/internal/config"
	"github.com/spf13/cobra"
)

type getCmdFlags struct {
	app          string
	revisionName string
	config       string
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
			Get details of a GitHub Runtime app.
			You can specify the app name using --app flag, --config flag to read from a runtime config file,
			or it will automatically read from runtime.config.json in the current directory if it exists.
		`),
		Example: heredoc.Doc(`
			$ gh runtime get --app my-app
			# => Retrieves details of the app named 'my-app'
			
			$ gh runtime get --config runtime.config.json
			# => Retrieves details using app name from the config file.
			
			$ gh runtime get
			# => Retrieves details using app name from runtime.config.json in current directory (if it exists).
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName, err := config.ResolveAppName(getCmdFlags.app, getCmdFlags.config)
			if err != nil {
				return err
			}

			getUrl := fmt.Sprintf("runtime/%s/deployment", appName)
			params := url.Values{}
			if getCmdFlags.revisionName != "" {
				params.Add("revision_name", getCmdFlags.revisionName)
			}
			if len(params) > 0 {
				getUrl += "?" + params.Encode()
			}

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
	getCmd.Flags().StringVarP(&getCmdFlags.config, "config", "c", "", "Path to runtime config file")
	getCmd.Flags().StringVarP(&getCmdFlags.revisionName, "revision-name", "r", "", "The revision name to use for the app")
	rootCmd.AddCommand(getCmd)
}
