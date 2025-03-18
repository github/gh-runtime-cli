package cmd

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

type deployCmdFlags struct {
	dir string
	app string
}

type createDeployReq struct {
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

			files, err := os.ReadDir(deployCmdFlags.dir)
			if err != nil {
				fmt.Printf("Error reading directory '%s': %v\n", deployCmdFlags.dir, err)
				return
			}

			for _, file := range files {
				deploymentsUrl := fmt.Sprintf("runtime/%s/deployment", deployCmdFlags.app)
				fmt.Printf("Deploying %s to %s\n", file.Name(), deploymentsUrl)

				// TODO: make request to create deployment
			}
		},
	}
	deployCmd.Flags().StringVarP(&deployCmdFlags.dir, "dir", "d", "", "The directory to deploy")
	deployCmd.Flags().StringVarP(&deployCmdFlags.app, "app", "a", "", "The app to deploy")
	rootCmd.AddCommand(deployCmd)
}
