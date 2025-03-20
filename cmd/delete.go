package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

type deleteCmdFlags struct {
	app string
}

type deleteResp struct {
}

func init() {
	deleteCmdFlags := deleteCmdFlags{}
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a GitHub Runtime app",
		Long: heredoc.Doc(`
			Delete a GitHub Runtime app
		`),
		Example: heredoc.Doc(`
			$ gh runtime delete --app my-app
			# => Deletes the app named 'my-app'
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if deleteCmdFlags.app == "" {
				fmt.Println("Error: --app flag is required")
				return
			}

			deleteUrl := fmt.Sprintf("runtime/%s/deployment", deleteCmdFlags.app)
			client, err := api.DefaultRESTClient()
			if err != nil {
				fmt.Println(err)
				return
			}
			var response string
			err = client.Delete(deleteUrl, &response)
			if err != nil {
				// print err and response
				fmt.Printf("Error deleting app: %v\n", err)
				fmt.Printf("Response: %v\n", response)
				return
			}

			fmt.Printf("App deleted: %s\n", response)
		},
	}

	deleteCmd.Flags().StringVarP(&deleteCmdFlags.app, "app", "a", "", "The app to delete")
	rootCmd.AddCommand(deleteCmd)
}
