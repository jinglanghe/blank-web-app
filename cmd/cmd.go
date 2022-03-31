package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "blankWebApp",
	Short: "blankWebApp management",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute the current command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
