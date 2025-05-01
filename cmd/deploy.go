package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"net/url"
	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

type deployCmdFlags struct {
	dir string
	app string
	revisionName string
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
			Deploys a directory to a GitHub Runtime app
		`),
		Example: heredoc.Doc(`
			$ gh runtime deploy --dir ./dist --app my-app
			# => Deploys the contents of the 'dist' directory to the app named 'my-app'
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if deployCmdFlags.dir == "" {
				fmt.Println("Error: --dir flag is required")
				return
			}
			if deployCmdFlags.app == "" {
				fmt.Println("Error: --app flag is required")
				return
			}

			if _, err := os.Stat(deployCmdFlags.dir); os.IsNotExist(err) {
				fmt.Printf("Error: directory '%s' does not exist\n", deployCmdFlags.dir)
				return
			}

			_, err := os.ReadDir(deployCmdFlags.dir)
			if err != nil {
				fmt.Printf("Error reading directory '%s': %v\n", deployCmdFlags.dir, err)
				return
			}

			// Zip the directory
			zipPath := fmt.Sprintf("%s.zip", deployCmdFlags.dir)
			err = zipDirectory(deployCmdFlags.dir, zipPath)
			if err != nil {
				fmt.Printf("Error zipping directory '%s': %v\n", deployCmdFlags.dir, err)
				return
			}
			defer os.Remove(zipPath)

			client, err := api.DefaultRESTClient()
			if err != nil {
				fmt.Println(err)
				return
			}

			deploymentsUrl := fmt.Sprintf("runtime/%s/deployment/bundle", deployCmdFlags.app)
			params := url.Values{}
			if deployCmdFlags.revisionName != "" {
				params.Add("revision_name", deployCmdFlags.revisionName)
			}
			if len(params) > 0 {
				deploymentsUrl += "?" + params.Encode()
			}

			fmt.Printf("Deploying app to %s\n", deploymentsUrl)

			// body is the full zip RAW
			body, err := os.ReadFile(zipPath)
			if err != nil {
				fmt.Printf("Error reading zip file '%s': %v\n", zipPath, err)
				return
			}

			err = client.Post(deploymentsUrl, bytes.NewReader(body), nil)
			if err != nil {
				fmt.Printf("Error deploying app: %v\n", err)
				return
			}

			fmt.Printf("Successfully deployed app\n")
		},
	}
	deployCmd.Flags().StringVarP(&deployCmdFlags.dir, "dir", "d", "", "The directory to deploy")
	deployCmd.Flags().StringVarP(&deployCmdFlags.app, "app", "a", "", "The app to deploy")
	deployCmd.Flags().StringVarP(&deployCmdFlags.revisionName, "revision-name", "r", "", "The revision name to deploy")
	rootCmd.AddCommand(deployCmd)
}
