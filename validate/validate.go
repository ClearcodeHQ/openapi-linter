package validate

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"os"
	"reflect"
)

// Store informations about a validation error
type ValidationError struct {
	FilePath string
	JsonPath string
	Err error
}

// Recursively iterates over jsonObject and tries to validate all nested objects as JSON object schemas.
// Returns the list of validation errors if found
func TraverseJSONObject(filePath string, jsonPath string, jsonObject map[string] interface{}, errors *map[string] ValidationError){
	for key, value := range jsonObject {
		childPath := jsonPath
		objectType := reflect.ValueOf(value).Kind()

		if objectType.String() == "map" {
			childPath = fmt.Sprintf("%s.%s", jsonPath, key)
			childObject := value.(map[string] interface{})

			TraverseJSONObject(filePath, childPath, childObject, errors)
		} else {
			validationKey := fmt.Sprintf("%s:%s", filePath, jsonPath)
			_, wasValidated := (*errors) [validationKey]
			if !wasValidated {
				err := ValidateSchema(&jsonObject);
				if err != nil {
					(*errors) [validationKey] = ValidationError{filePath, jsonPath, err}
				}
			}
		}
	}
}


// Interface to make testing easier
type JsonFileInfo interface {
	Name() string
	IsDir() bool
}

// Determines if file is a directory or a possible JSON file.
// Doesn't check the mimetype (yet).
func IsJsonFile(fileInfo JsonFileInfo) bool {
	if fileInfo.IsDir() == true {
		return false
	}

	if strings.HasSuffix(fileInfo.Name(), ".json") {
		return true;
	}
	return false
}

// Scans the directory and returns only JSON files
func FindJsonFiles(dir string) ([] string, error) {
	paths := [] string {}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info != nil && IsJsonFile(info) {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}
// validates all schemas within given directory
func ValidateAllSchemasInDir(dir string) (map[string]ValidationError, error) {

	jsonFiles, err := FindJsonFiles(dir)
	if err != nil {
		return map[string] ValidationError{}, err
	}

	jsonErrors := make(map[string] ValidationError)
	for _, file := range jsonFiles {
			if fileErr := ValidateJSONFile(file, &jsonErrors); fileErr != nil {
				return map[string] ValidationError{}, fileErr
			}
		}
	return jsonErrors, err
}

// Opens a file and unmarshals its content and validates it.
func ValidateJSONFile(schemaPath string, jsonErrors *map[string] ValidationError) error {
	jsonFile, err := os.Open(schemaPath)

	if err != nil {
		return err
	}
	jsonBytes, err  := ioutil.ReadAll(jsonFile);

	if err != nil {
		return err
	}

	jsonSchema := map[string] interface {}{}

	err = json.Unmarshal(jsonBytes, &jsonSchema)
	if err != nil {
		return err
	}

	TraverseJSONObject( schemaPath, "", jsonSchema, jsonErrors)

	defer jsonFile.Close()
	return nil;
}

// Check the object is a valid JSON Schema.
func ValidateSchema(schemaContent *map[string] interface{}) (error) {
	schemaLoader := gojsonschema.NewSchemaLoader()
	schemaLoader.Validate = true
	fileLoader := gojsonschema.NewGoLoader(schemaContent)
	if err := schemaLoader.AddSchemas(fileLoader); err != nil {
		return err
	}
	return nil;
}