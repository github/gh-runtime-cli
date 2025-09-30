package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/github/gh-runtime-cli/internal/config"
	"github.com/spf13/cobra"
)

type deployCmdFlags struct {
	dir          string
	app          string
	revisionName string
	sha          string
	config       string
}

func zipDirectory(sourceDir, destinationZip string) error {
	zipFile, err := os.Create(destinationZip)
	if err != nil {
		return fmt.Errorf("error creating zip file '%s': %w", destinationZip, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path '%s': %w", path, err)
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("error calculating relative path for '%s': %w", path, err)
		}

		if info.IsDir() {
			if relPath == "." {
				return nil
			}
			relPath += "/"
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("error creating zip header for '%s': %w", path, err)
		}
		header.Name = relPath
		if !info.IsDir() {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("error creating zip writer for '%s': %w", path, err)
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("error opening file '%s': %w", path, err)
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return fmt.Errorf("error writing file '%s' to zip: %w", path, err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error zipping directory '%s': %w", sourceDir, err)
	}

	return nil
}

func init() {
	deployCmdFlags := deployCmdFlags{}
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy app to GitHub Runtime",
		Long: heredoc.Doc(`
			Deploys a directory to a GitHub Runtime app.
			You can specify the app name using --app flag, --config flag to read from a runtime config file,
			or it will automatically read from runtime.config.json in the current directory if it exists.
		`),
		Example: heredoc.Doc(`
			$ gh runtime deploy --dir ./dist --app my-app [--sha <sha>]
			# => Deploys the contents of the 'dist' directory to the app named 'my-app'.
			
			$ gh runtime deploy --dir ./dist --config runtime.config.json
			# => Deploys using app name from the config file.
			
			$ gh runtime deploy --dir ./dist
			# => Deploys using app name from runtime.config.json in current directory (if it exists).
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if deployCmdFlags.dir == "" {
				return fmt.Errorf("--dir flag is required")
			}

			appName := deployCmdFlags.app

			// If config file is provided, read app name from it
			if deployCmdFlags.config != "" {
				configApp, err := config.ReadRuntimeConfig(deployCmdFlags.config)
				if err != nil {
					return err
				}
				if appName == "" {
					appName = configApp
				}
			} else if appName == "" {
				// Try to read from default config file if neither --app nor --config is provided
				if _, err := os.Stat("runtime.config.json"); err == nil {
					configApp, err := config.ReadRuntimeConfig("runtime.config.json")
					if err != nil {
						return fmt.Errorf("found runtime.config.json but failed to read it: %v", err)
					}
					appName = configApp
				}
			}

			if appName == "" {
				return fmt.Errorf("--app flag is required, --config must be specified, or runtime.config.json must exist in current directory")
			}

			if _, err := os.Stat(deployCmdFlags.dir); os.IsNotExist(err) {
				return fmt.Errorf("directory '%s' does not exist", deployCmdFlags.dir)
			}

			_, err := os.ReadDir(deployCmdFlags.dir)
			if err != nil {
				return fmt.Errorf("error reading directory '%s': %v", deployCmdFlags.dir, err)
			}

			// Zip the directory
			zipPath := fmt.Sprintf("%s.zip", deployCmdFlags.dir)
			err = zipDirectory(deployCmdFlags.dir, zipPath)
			if err != nil {
				return fmt.Errorf("error zipping directory '%s': %v", deployCmdFlags.dir, err)
			}
			defer os.Remove(zipPath)

			client, err := api.DefaultRESTClient()
			if err != nil {
				return fmt.Errorf("error creating REST client: %v", err)
			}

			deploymentsUrl := fmt.Sprintf("runtime/%s/deployment/bundle", appName)
			params := url.Values{}

			if deployCmdFlags.revisionName != "" {
				params.Add("revision_name", deployCmdFlags.revisionName)
			}

			if deployCmdFlags.sha != "" {
				params.Add("revision", deployCmdFlags.sha)
			}

			if len(params) > 0 {
				deploymentsUrl += "?" + params.Encode()
			}

			fmt.Printf("Deploying app to %s\n", deploymentsUrl)

			// body is the full zip RAW
			body, err := os.ReadFile(zipPath)
			if err != nil {
				return fmt.Errorf("error reading zip file '%s': %v", zipPath, err)
			}

			err = client.Post(deploymentsUrl, bytes.NewReader(body), nil)
			if err != nil {
				return fmt.Errorf("error deploying app: %v", err)
			}

			fmt.Printf("Successfully deployed app\n")
			return nil
		},
	}
	deployCmd.Flags().StringVarP(&deployCmdFlags.dir, "dir", "d", "", "The directory to deploy")
	deployCmd.Flags().StringVarP(&deployCmdFlags.app, "app", "a", "", "The app to deploy")
	deployCmd.Flags().StringVarP(&deployCmdFlags.config, "config", "c", "", "Path to runtime config file")
	deployCmd.Flags().StringVarP(&deployCmdFlags.revisionName, "revision-name", "r", "", "The revision name to deploy")
	deployCmd.Flags().StringVarP(&deployCmdFlags.sha, "sha", "s", "", "SHA of the app being deployed")

	rootCmd.AddCommand(deployCmd)
}
