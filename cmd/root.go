package cmd

import (
        "os"

        "github.com/spf13/cobra"
)

var cfgFile string

// rootCMD represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
        Use:   "openapi-linter",
        Short: "A tool that scans the provided directory and validates all JSON Schema objects inside it.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCMD.
func Execute() {
        if err := rootCmd.Execute(); err != nil {
                os.Exit(1)
        }
}

func init() {
        cobra.OnInitialize()

        // Cobra also supports local flags, which will only run
        // when this action is called directly.
        // rootCMD.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}