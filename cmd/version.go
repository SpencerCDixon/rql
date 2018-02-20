package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.0.1"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "All software needs a version.  This is RQL's.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("RQL Version: %s\n", Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
