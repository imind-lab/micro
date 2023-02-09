package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "show microctl version",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("v2.0.0")
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}
