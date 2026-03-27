package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/github/gh-runtime-cli/internal/config"
	"github.com/spf13/cobra"
)

type initCmdFlags struct {
	app string
	out string
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
			client, err := api.DefaultRESTClient()
			if err != nil {
				return fmt.Errorf("failed creating REST client: %v", err)
			}

			return runInit(client, initCmdFlags)
		},
	}

	initCmd.Flags().StringVarP(&initCmdFlags.app, "app", "a", "", "The app name to initialize")
	initCmd.Flags().StringVarP(&initCmdFlags.out, "out", "o", "", "The output path for the runtime.config.json file (default: runtime.config.json in current directory)")
	rootCmd.AddCommand(initCmd)
}

func runInit(client restClient, flags initCmdFlags) error {
	if flags.app == "" {
		return fmt.Errorf("--app flag is required")
	}

	getUrl := fmt.Sprintf("runtime/%s/deployment", flags.app)

	response := appResponse{}
	err := client.Get(getUrl, &response)
	if err != nil {
		return fmt.Errorf("app '%s' does not exist or is not accessible: %v", flags.app, err)
	}

	return writeRuntimeConfig(flags.app, flags.out)
}

// writeRuntimeConfig writes a runtime.config.json file for the given app.
// If outPath is empty, it defaults to "runtime.config.json" in the current directory.
func writeRuntimeConfig(app string, outPath string) error {
	configStruct := config.RuntimeConfig{
		App: app,
	}

	configPath := "runtime.config.json"
	if outPath != "" {
		configPath = outPath
		outputDir := filepath.Dir(configPath)
		if outputDir != "." {
			err := os.MkdirAll(outputDir, 0755)
			if err != nil {
				return fmt.Errorf("error creating directory '%s': %v", outputDir, err)
			}
		}
	}

	configBytes, err := json.MarshalIndent(configStruct, "", "  ")
	if err != nil {
		return fmt.Errorf("error creating configuration: %v", err)
	}

	err = os.WriteFile(configPath, configBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing configuration file: %v", err)
	}

	fmt.Printf("Successfully initialized local project for Spark app '%s' at '%s'\n", app, configPath)
	return nil
}
