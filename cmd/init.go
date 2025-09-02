package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

type initCmdFlags struct {
	app string
	out string
}

type runtimeConfig struct {
	App string `json:"app"`
}

type appResponse struct {
	AppUrl string `json:"app_url"`
}

func init() {
	initCmdFlags := initCmdFlags{}
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a local project for GitHub Spark",
		Long: heredoc.Doc(`
			Initialize a local project to connect it to a GitHub Spark app.
			This creates a runtime.config.json configuration file that binds your local project
			to a remote Spark app. You must specify an app name to validate the app exists.
			Optionally specify an output path where the runtime.config.json file should be created.
		`),
		Example: heredoc.Doc(`
			$ gh runtime init --app my-spark-app
			# => Binds local project to the Spark app 'my-spark-app'
			
			$ gh runtime init --app my-spark-app --out ./config/runtime.config.json
			# => Creates configuration at the specified path
			
			$ gh runtime init --app my-spark-app --out ./my-config.json
			# => Creates configuration with a custom filename
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if initCmdFlags.app == "" {
				return fmt.Errorf("--app flag is required")
			}

			// Determine the identifier to use for the API call
			identifier := initCmdFlags.app

			getUrl := fmt.Sprintf("runtime/%s/deployment", identifier)

			client, err := api.DefaultRESTClient()
			if err != nil {
				return fmt.Errorf("failed creating REST client: %v", err)
			}

			response := appResponse{}
			err = client.Get(getUrl, &response)
			if err != nil {
				return fmt.Errorf("app '%s' does not exist or is not accessible: %v", identifier, err)
			}

			// Create runtime config
			config := runtimeConfig{
				App: initCmdFlags.app,
			}

			configPath := "runtime.config.json"
			if initCmdFlags.out != "" {
				configPath = initCmdFlags.out
				// Create directory if it doesn't exist
				outputDir := filepath.Dir(configPath)
				if outputDir != "." {
					err = os.MkdirAll(outputDir, 0755)
					if err != nil {
						return fmt.Errorf("error creating directory '%s': %v", outputDir, err)
					}
				}
			}

			configBytes, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return fmt.Errorf("error creating configuration: %v", err)
			}

			err = os.WriteFile(configPath, configBytes, 0644)
			if err != nil {
				return fmt.Errorf("error writing configuration file: %v", err)
			}

			fmt.Printf("Successfully initialized local project for Spark app '%s' at '%s'\n", identifier, configPath)
			return nil
		},
	}

	initCmd.Flags().StringVarP(&initCmdFlags.app, "app", "a", "", "The app name to initialize")
	initCmd.Flags().StringVarP(&initCmdFlags.out, "out", "o", "", "The output path for the runtime.config.json file (default: runtime.config.json in current directory)")
	rootCmd.AddCommand(initCmd)
}
