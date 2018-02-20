package cmd

import (
	"fmt"
	"os"

	"github.com/spencercdixon/rql/repl"
	"github.com/spencercdixon/rql/rql"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "rql",
	Short: "Console for RQL DB",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("rql (%s)\n", Version)
		fmt.Println(`Type '\?' for help`)
		println()

		var dbLoc string
		if len(args) > 0 {
			dbLoc = args[0]
		} else {
			dbLoc = "tmp"
		}

		db := rql.New(dbLoc)
		repl.Start(db, os.Stdin, os.Stdout)
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
