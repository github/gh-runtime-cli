package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"net/url"
	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

type createCmdFlags struct {
	app                  string
	EnvironmentVariables []string
	Secrets              []string
	RevisionName 	  	 string
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
			$ gh runtime create --app my-app --env key1=value1 --env key2=value2 --secret key3=value3 --secret key4=value4
			# => Creates the app named 'my-app'
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if createCmdFlags.app == "" {
				fmt.Println("Error: --app flag is required")
				return
			}

			// Construct the request body
			requestBody := createReq{
				EnvironmentVariables: map[string]string{},
				Secrets:              map[string]string{},
			}

			for _, pair := range createCmdFlags.EnvironmentVariables {
				parts := strings.SplitN(pair, "=", 2)
				if len(parts) == 2 {
					key := parts[0]
					value := parts[1]
					requestBody.EnvironmentVariables[key] = value
				} else {
					fmt.Printf("Error: Invalid environment variable format (%s). Must be in the form 'key=value'\n", pair)
					return
				}
			}

			for _, pair := range createCmdFlags.Secrets {
				parts := strings.SplitN(pair, "=", 2)
				if len(parts) == 2 {
					key := parts[0]
					value := parts[1]
					requestBody.Secrets[key] = value
				} else {
					fmt.Printf("Error: Invalid secret format (%s). Must be in the form 'key=value'\n", pair)
					return
				}
			}

			body, err := json.Marshal(requestBody)
			if err != nil {
				fmt.Printf("Error marshalling request body: %v\n", err)
				return
			}

			createUrl := fmt.Sprintf("runtime/%s/deployment", createCmdFlags.app)
			params := url.Values{}
			if createCmdFlags.RevisionName != "" {
				params.Add("revision_name", createCmdFlags.RevisionName)
			}
			if len(params) > 0 {
				createUrl += "?" + params.Encode()
			}
			
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
	createCmd.Flags().StringSliceVarP(&createCmdFlags.EnvironmentVariables, "env", "e", []string{}, "Environment variables to set on the app in the form 'key=value'")
	createCmd.Flags().StringSliceVarP(&createCmdFlags.Secrets, "secret", "s", []string{}, "Secrets to set on the app in the form 'key=value'")
	createCmd.Flags().StringVarP(&createCmdFlags.RevisionName, "revision-name", "r", "", "The revision name to use for the app")
	rootCmd.AddCommand(createCmd)
}
