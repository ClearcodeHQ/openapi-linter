package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/clearcodehq/openapi-linter/validate"
)

var validateCmd = &cobra.Command{
	Use: "validate",
	Short: "Scan all JSON files in the directory and validate all JSON Schemas found in those files.",
	SilenceUsage: true,
	Args: cobra.ExactArgs(1),
	RunE: func (cmd *cobra.Command, args []string) error {
		validationErrors, err := validate.ValidateAllSchemasInDir(args[0]);
		if err != nil {
			return fmt.Errorf("Couldn't parse some files.")
		}

		if len(validationErrors) > 0{
			displayErrors(&validationErrors)
			return fmt.Errorf("The validation has failed.")
		}
		return nil;
	},
}
func displayErrors(errors *map[string] validate.ValidationError) {
	for _, err := range *errors {
		fmt.Printf("File: %s\nJSONPath: %s\nError message:\n%s", err.FilePath, err.JsonPath, err.Err)
	}
}

func init() {
	rootCmd.AddCommand(validateCmd)
}