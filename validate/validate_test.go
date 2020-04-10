package validate_test

import (
	"github.com/clearcodehq/openapi-linter/validate"
	"fmt"
	"reflect"
	"testing"
)

func TestValidateJSONFile(t *testing.T) {
	t.Run("Empty object", func(t *testing.T) {
		// GIVEN
		jsonErrors := make(map[string]validate.ValidationError)
		jsonSchema := map[string]interface{}{}

		// WHEN
		validate.TraverseJSONObject("", "", jsonSchema, &jsonErrors)

		// THEN
		if len(jsonErrors) > 0 {
			t.Errorf("Empty object shouldn't raise an error.")
		}
	})
	t.Run("The root object is a schema", func(t *testing.T) {
		// GIVEN
		jsonErrors := make(map[string]validate.ValidationError)
		jsonSchema := map[string]interface{}{
			"$schema": "http://json-schema.org/draft-07/schema",
			"properties": map[string]interface{}{
				"evilfield": map[string]interface{}{
					"type": "boolean",
				},
			},
		}

		// WHEN
		validate.TraverseJSONObject("", "", jsonSchema, &jsonErrors)

		// THEN
		if len(jsonErrors) > 0 {
			t.Errorf("Empty object  shouldn't raise an error")
		}
	})
	t.Run("The root object is a schema and has an error", func(t *testing.T) {
		// GIVEN
		expectedErrors := map[string]validate.ValidationError{
			"xxx.json:root": validate.ValidationError{
				"xxx.json",
				"root",
				fmt.Errorf("properties.evilfield.type: Must validate at least one schema (anyOf)\nproperties.evilfield.type: properties.evilfield.type must be one of the following: \"array\", \"boolean\", \"integer\", \"null\", \"number\", \"object\", \"string\"\n"),
			},
		}
		jsonErrors := make(map[string]validate.ValidationError)
		jsonSchema := map[string]interface{}{
			"$schema": "http://json-schema.org/draft-07/schema",
			"properties": map[string]interface{}{
				"evilfield": map[string]interface{}{
					"type": "unsupported error",
				},
			},
		}

		// WHEN
		validate.TraverseJSONObject("xxx.json", "root", jsonSchema, &jsonErrors)

		// THEN
		if !reflect.DeepEqual(jsonErrors, expectedErrors) {
			t.Errorf("AssertionFail: %+v != %+v", jsonErrors, expectedErrors)
		}
	})

	t.Run("The root object is a schema and has an error", func(t *testing.T) {
		// GIVEN
		expectedErrors := map[string]validate.ValidationError{
			"xxx.json:.200": validate.ValidationError{
				"xxx.json",
				".200",
				fmt.Errorf("properties.evilfield.type: Must validate at least one schema (anyOf)\nproperties.evilfield.type: properties.evilfield.type must be one of the following: \"array\", \"boolean\", \"integer\", \"null\", \"number\", \"object\", \"string\"\n"),
			},
			"xxx.json:.400": validate.ValidationError{
				"xxx.json",
				".400",
				fmt.Errorf("properties.evilfield.type: Must validate at least one schema (anyOf)\nproperties.evilfield.type: properties.evilfield.type must be one of the following: \"array\", \"boolean\", \"integer\", \"null\", \"number\", \"object\", \"string\"\n"),
			},
		}
		jsonErrors := make(map[string]validate.ValidationError)
		jsonSchema := map[string]interface{}{
			"200": map[string]interface{}{
				"$schema": "http://json-schema.org/draft-07/schema",
				"properties": map[string]interface{}{
					"evilfield": map[string]interface{}{
						"type": "unsupported error",
					},
				},
			},
			"400": map[string]interface{}{
				"$schema": "http://json-schema.org/draft-07/schema",
				"properties": map[string]interface{}{
					"evilfield": map[string]interface{}{
						"type": "unsupported error",
					},
				},
			},
		}

		// WHEN
		validate.TraverseJSONObject("xxx.json", "", jsonSchema, &jsonErrors)

		// THEN
		if !reflect.DeepEqual(jsonErrors, expectedErrors) {
			t.Errorf("AssertionFail: %+v != %+v", jsonErrors, expectedErrors)
		}
	})
}

type MockFileInfo struct {
	fileName    string
	isDirectory bool
}

func (this MockFileInfo) Name() string {
	return this.fileName
}

func (this MockFileInfo) IsDir() bool {
	return this.isDirectory
}

func TestIsJsonFile(t *testing.T) {
	t.Run("Skip non-json files", func(t *testing.T) {
		file := MockFileInfo{
			"aaa.txt",
			false,
		}
		if validate.IsJsonFile(file) != false {
			t.Errorf("File without the .json extension isn't a json file.")
		}
	})
	t.Run("Skip directories", func(t *testing.T) {
		file := MockFileInfo{
			"aaa",
			true,
		}
		if validate.IsJsonFile(file) != false {
			t.Errorf("Exclude directories")
		}
	})
	t.Run("Skip directories with `.json` suffix", func(t *testing.T) {
		file := MockFileInfo{
			"aaa.json",
			true,
		}
		if validate.IsJsonFile(file) != false {
			t.Errorf("Exclude directories with .json suffix.")
		}
	})
	t.Run("Find a json file", func(t *testing.T) {
		file := MockFileInfo{
			"aaa.json",
			false,
		}
		if validate.IsJsonFile(file) != true {
			t.Errorf("JSON file not found.")
		}
	})
}
