package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of docsql",
	Long:  `All software has versions. This is docsql's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("docsql v1.0.1")
	},
}
