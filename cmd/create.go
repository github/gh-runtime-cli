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
	name                 string
	visibility           string
	org                  string
	environmentVariables []string
	secrets              []string
	revisionName         string
	init                 bool
}

type createReq struct {
	Name                 string            `json:"friendly_name,omitempty"`
	Visibility           string            `json:"visibility,omitempty"`
	OrganizationLogin    string            `json:"organization_login,omitempty"`
	EnvironmentVariables map[string]string `json:"environment_variables"`
	Secrets              map[string]string `json:"secrets"`
}

type createResp struct {
	AppUrl string `json:"app_url"`
	ID     string `json:"id"`
}

func init() {
	createCmdFlags := createCmdFlags{}
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a GitHub Runtime app",
		Long: heredoc.Doc(`
			Create a GitHub Runtime app.
		`),
		Example: heredoc.Doc(`
			$ gh runtime create --app my-app --env key1=value1 --env key2=value2 --secret key3=value3 --secret key4=value4
			# => Creates the app with the ID 'my-app'

			$ gh runtime create --name my-new-app
			# => Creates a new app with the given name

			$ gh runtime create --app my-app --visibility only_owner
			# => Creates the app visible only to the owner

			$ gh runtime create --app my-app --visibility github
			# => Creates the app visible to all GitHub users

			$ gh runtime create --app my-app --visibility selected_orgs --org my-org
			# => Creates the app visible to 'my-org' organization
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.DefaultRESTClient()
			if err != nil {
				return fmt.Errorf("failed creating REST client: %v", err)
			}

			resp, err := runCreate(client, createCmdFlags)
			if err != nil {
				return err
			}

			fmt.Printf("App created: %s\n", resp.AppUrl)
			if resp.ID != "" {
				fmt.Printf("ID: %s\n", resp.ID)
			}
			return nil
		},
	}

	createCmd.Flags().StringVarP(&createCmdFlags.app, "app", "a", "", "The app ID to create")
	createCmd.Flags().StringVarP(&createCmdFlags.name, "name", "n", "", "The name for the app")
	createCmd.Flags().StringVarP(&createCmdFlags.visibility, "visibility", "v", "", "The visibility of the app (e.g. 'only_owner', 'github', or 'selected_orgs')")
	createCmd.Flags().StringVarP(&createCmdFlags.org, "org", "o", "", "The organization login to grant access (only valid with --visibility=selected_orgs)")
	createCmd.Flags().StringSliceVarP(&createCmdFlags.environmentVariables, "env", "e", []string{}, "Environment variables to set on the app in the form 'key=value'")
	createCmd.Flags().StringSliceVarP(&createCmdFlags.secrets, "secret", "s", []string{}, "Secrets to set on the app in the form 'key=value'")
	createCmd.Flags().StringVarP(&createCmdFlags.revisionName, "revision-name", "r", "", "The revision name to use for the app")
	createCmd.Flags().BoolVar(&createCmdFlags.init, "init", false, "Initialize a runtime.config.json file in the current directory after creating the app")
	rootCmd.AddCommand(createCmd)
}

func runCreate(client restClient, flags createCmdFlags) (createResp, error) {
	if flags.app == "" && flags.name == "" {
		return createResp{}, fmt.Errorf("either --app or --name flag is required")
	}

	if flags.org != "" && flags.visibility != "selected_orgs" {
		return createResp{}, fmt.Errorf("--org can only be used with --visibility=selected_orgs")
	}

	if flags.visibility == "selected_orgs" && flags.org == "" {
		return createResp{}, fmt.Errorf("--org is required when --visibility=selected_orgs")
	}

	requestBody := createReq{
		Name:                 flags.name,
		Visibility:           flags.visibility,
		OrganizationLogin:    flags.org,
		EnvironmentVariables: map[string]string{},
		Secrets:              map[string]string{},
	}

	for _, pair := range flags.environmentVariables {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			requestBody.EnvironmentVariables[parts[0]] = parts[1]
		} else {
			return createResp{}, fmt.Errorf("invalid environment variable format (%s). Must be in the form 'key=value'", pair)
		}
	}

	for _, pair := range flags.secrets {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			requestBody.Secrets[parts[0]] = parts[1]
		} else {
			return createResp{}, fmt.Errorf("invalid secret format (%s). Must be in the form 'key=value'", pair)
		}
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return createResp{}, fmt.Errorf("error marshalling request body: %v", err)
	}

	var createUrl string
	if flags.app != "" {
		createUrl = fmt.Sprintf("runtime/%s/deployment", flags.app)
	} else {
		createUrl = "runtime"
	}

	params := url.Values{}
	if flags.revisionName != "" {
		params.Add("revision_name", flags.revisionName)
	}
	if len(params) > 0 {
		createUrl += "?" + params.Encode()
	}

	response := createResp{}
	err = client.Put(createUrl, bytes.NewReader(body), &response)
	if err != nil {
		return createResp{}, fmt.Errorf("error creating app: %v", err)
	}

	if flags.init {
		if response.ID == "" {
			return response, fmt.Errorf("error initializing config: server did not return an app ID")
		}
		err = writeRuntimeConfig(response.ID, "")
		if err != nil {
			return response, fmt.Errorf("error initializing config: %v", err)
		}
	}

	return response, nil
}
