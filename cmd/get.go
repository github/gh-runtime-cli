package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

type getCmdFlags struct {
	app string
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
		Run: func(cmd *cobra.Command, args []string) {
			if getCmdFlags.app == "" {
				fmt.Println("Error: --app flag is required")
				return
			}

			getUrl := fmt.Sprintf("runtime/%s/deployment", getCmdFlags.app)
			client, err := api.DefaultRESTClient()
			if err != nil {
				fmt.Println(err)
				return
			}

			response := json.RawMessage{}
			err = client.Get(getUrl, &response)
			if err != nil {
				fmt.Printf("Error retrieving app details: %v\n", err)
				return
			}

			fmt.Printf("App Details: %s\n", response)
		},
	}

	getCmd.Flags().StringVarP(&getCmdFlags.app, "app", "a", "", "The app to retrieve details for")
	rootCmd.AddCommand(getCmd)
}
