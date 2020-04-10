package cmd

import (
	validate_examples "github.com/clearcodehq/openapi-linter/validate-examples"
	"fmt"
	"github.com/spf13/cobra"
)

var validateExamplesCmd = &cobra.Command{
	Use:          "validate-examples",
	Short:        "Validate if an example matches the schema defined in the API spec.",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		errors := validate_examples.ScanForExamples(args[0])
		if len(errors) > 0 {
			for _, err := range errors {
				cmd.Println(err)
			}
			return fmt.Errorf("Validation errors found.")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateExamplesCmd)
}
