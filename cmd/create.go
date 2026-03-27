package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

type createCmdFlags struct {
	app                  string
	EnvironmentVariables []string
	Secrets              []string
	RevisionName 	  	 string
	Init                 bool
}

type createReq struct {
	EnvironmentVariables map[string]string `json:"environment_variables"`
	Secrets              map[string]string `json:"secrets"`
}

type createResp struct {
	AppUrl string `json:"app_url"`
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
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.DefaultRESTClient()
			if err != nil {
				return fmt.Errorf("failed creating REST client: %v", err)
			}

			appUrl, err := runCreate(client, createCmdFlags)
			if err != nil {
				return err
			}

			fmt.Printf("App created: %s\n", appUrl)
			return nil
		},
	}

	createCmd.Flags().StringVarP(&createCmdFlags.app, "app", "a", "", "The app to create")
	createCmd.Flags().StringSliceVarP(&createCmdFlags.EnvironmentVariables, "env", "e", []string{}, "Environment variables to set on the app in the form 'key=value'")
	createCmd.Flags().StringSliceVarP(&createCmdFlags.Secrets, "secret", "s", []string{}, "Secrets to set on the app in the form 'key=value'")
	createCmd.Flags().StringVarP(&createCmdFlags.RevisionName, "revision-name", "r", "", "The revision name to use for the app")
	createCmd.Flags().BoolVar(&createCmdFlags.Init, "init", false, "Initialize a runtime.config.json file in the current directory after creating the app")
	rootCmd.AddCommand(createCmd)
}

func runCreate(client restClient, flags createCmdFlags) (string, error) {
	if flags.app == "" {
		return "", fmt.Errorf("--app flag is required")
	}

	requestBody := createReq{
		EnvironmentVariables: map[string]string{},
		Secrets:              map[string]string{},
	}

	for _, pair := range flags.EnvironmentVariables {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			requestBody.EnvironmentVariables[parts[0]] = parts[1]
		} else {
			return "", fmt.Errorf("invalid environment variable format (%s). Must be in the form 'key=value'", pair)
		}
	}

	for _, pair := range flags.Secrets {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			requestBody.Secrets[parts[0]] = parts[1]
		} else {
			return "", fmt.Errorf("invalid secret format (%s). Must be in the form 'key=value'", pair)
		}
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %v", err)
	}

	createUrl := fmt.Sprintf("runtime/%s/deployment", flags.app)
	params := url.Values{}
	if flags.RevisionName != "" {
		params.Add("revision_name", flags.RevisionName)
	}
	if len(params) > 0 {
		createUrl += "?" + params.Encode()
	}

	response := createResp{}
	err = client.Put(createUrl, bytes.NewReader(body), &response)
	if err != nil {
		return "", fmt.Errorf("error creating app: %v", err)
	}

	if flags.Init {
		err = writeRuntimeConfig(flags.app, "")
		if err != nil {
			return response.AppUrl, fmt.Errorf("error initializing config: %v", err)
		}
	}

	return response.AppUrl, nil
}
