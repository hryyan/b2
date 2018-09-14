package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "0.2.0"

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
