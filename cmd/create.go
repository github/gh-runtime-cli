package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

type createCmdFlags struct {
	app string
}

type createReq struct {
	EnvironmentVariables map[string]string `json:"environment_variables"`
	Secrets              map[string]string `json:"secrets"`
}

type createResp struct {
}

func init() {
	createCmdFlags := createCmdFlags{}
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a GitHub Runtime app",
		Long: heredoc.Doc(`
			Create a GitHub Runtime app
		`),
		Example: heredoc.Doc(`
			$ gh runtime create --app my-app
			# => Creates the app named 'my-app'
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if createCmdFlags.app == "" {
				fmt.Println("Error: --app flag is required")
				return
			}

			// Construct the request body
			requestBody := createReq{
				EnvironmentVariables: map[string]string{
					"EXAMPLE_ENV": "value1",
				},
				Secrets: map[string]string{
					"SECRET_KEY": "secret_value",
				},
			}

			body, err := json.Marshal(requestBody)
			if err != nil {
				fmt.Printf("Error marshalling request body: %v\n", err)
				return
			}

			createUrl := fmt.Sprintf("runtime/%s/deployment", createCmdFlags.app)
			client, err := api.DefaultRESTClient()
			if err != nil {
				fmt.Println(err)
				return
			}
			var response string
			err = client.Put(createUrl, bytes.NewReader(body), &response)
			if err != nil {
				fmt.Printf("Error creating app: %v\n", err)
				return
			}

			fmt.Printf("App created: %s\n", response) // TODO pretty print details
		},
	}

	createCmd.Flags().StringVarP(&createCmdFlags.app, "app", "a", "", "The app to create")
	rootCmd.AddCommand(createCmd)
}
