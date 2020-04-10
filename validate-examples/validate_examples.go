package validate_examples

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/PaesslerAG/jsonpath"
	"github.com/bmatcuk/doublestar"
	"github.com/xeipuuv/gojsonschema"

	"io/ioutil"
	"path/filepath"
	"strings"
)

type Example struct {
	schemaPath  string
	examplePath string
}

const arrayError = "Invalid type. Expected: object, given: array"

// Find all examples and their respective schemas.
// Unfortunately, `jsonpath` seems not to implement JsonPath filters on recursive traversal queries.
// This introduces an layer of complexity related to the dynamic traversal of the JSON nodes.
// The function can be simplified if you'll find a lib that implements queries like:
// $..[?(@.examples)] <- find all nodes that have the examples property
// Or if you find any library that has a decent implementation of JSON Walker
// Related: https://github.com/PaesslerAG/jsonpath/issues/24
func FindExamples(jsonObject map[string]interface{}, cb func(Example, error)) {
	jsonNodes, _ := jsonpath.Get("$..*", jsonObject)
	nodes := sortNodes(jsonNodes.([]interface{}))

	for _, node := range nodes {

		schema, _ := jsonpath.Get(`$["schema"]["$ref"]`, node)
		example, _ := jsonpath.Get(`$["example"]["$ref"]`, node)

		schemaStr, schemaOk := schema.(string)
		exampleStr, exampleOk := example.(string)

		if !schemaOk || !exampleOk {
			if exampleOk {
				cb(Example{}, fmt.Errorf("can't find schema of the example"))
			}
			continue
		}

		cb(Example{schemaStr, exampleStr}, nil)
	}
}

// Deterministic sort json nodes.
// Return only nodes that are map[string]interface{}.
func sortNodes(jsonNodes []interface{}) []map[string]interface{} {

	nodes := map[string]interface{}{}
	keys := []string(nil)

	for _, node := range jsonNodes {
		key := fmt.Sprint(node)
		keys = append(keys, key)
		nodes[key] = node
	}

	sort.Strings(keys)

	sortedNodes := []map[string]interface{}(nil)
	for _, key := range keys {
		if node, ok := nodes[key].(map[string]interface{}); ok {
			sortedNodes = append(sortedNodes, node)
		}
	}

	return sortedNodes
}

// Read JSON object from a file and Unmarshal it as a generic map.
func GetObjectFromFile(filePath string) (map[string]interface{}, error) {
	fileObject := map[string]interface{}{}
	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't read the file:%s:%s", filePath, err)
	}

	err = json.Unmarshal(jsonBytes, &fileObject)

	if err != nil {
		return nil, fmt.Errorf("can't unmarshal contents: %s: %+v", filePath, err)
	}
	return fileObject, nil
}

// Generator to retrieve contents of next json files and pass them to a scan function.
func ScanJSONFiles(mainPath string, scanFunction func(string)) {
	jsonFiles, _ := doublestar.Glob(mainPath + "/**/*.json")
	for _, jsonFile := range jsonFiles {
		if strings.Contains(jsonFile, ".partial.json") {
			continue
		}
		file, err := os.Stat(jsonFile)

		if err != nil || !file.IsDir() {
			scanFunction(jsonFile)
		}
	}
}

func ScanForExamples(rootPath string) []error {
	var errors []error
	ScanJSONFiles(rootPath, func(jsonPath string) {
		jsonObject, err := GetObjectFromFile(jsonPath)

		if err != nil {
			errors = append(errors, err)
			return
		}
		FindExamples(jsonObject, func(example Example, parseErr error) {
			if example.examplePath == "" || example.schemaPath == "" {
				errors = append(errors, parseErr)
				return
			}

			examplePath := filepath.Join(filepath.Dir(jsonPath), example.examplePath)
			schemaPath := filepath.Join(filepath.Dir(jsonPath), example.schemaPath)

			// Windows
			examplePath = strings.ReplaceAll(examplePath, `\`, `/`)
			schemaPath = strings.ReplaceAll(schemaPath, `\`, `/`)

			if parseErr != nil {
				errors = append(errors, parseErr)
				return
			}

			exampleLoader, exampleLoaderErr := GetReferenceLoader(examplePath)
			exampleSchemaLoader, exampleSchemaLoaderErr := GetReferenceLoader(schemaPath)

			if exampleLoaderErr != nil {
				errors = append(errors, fmt.Errorf("[example=%s, schema=%s] %s", example.examplePath, example.schemaPath, exampleLoaderErr))
				return
			}

			if exampleSchemaLoaderErr != nil {
				errors = append(errors, fmt.Errorf("[example=%s, schema=%s] %s", example.examplePath, example.schemaPath, exampleSchemaLoaderErr))
				return
			}

			result, valErr := gojsonschema.Validate(*exampleSchemaLoader, *exampleLoader)
			if valErr != nil {
				errors = append(errors, fmt.Errorf("%s: %s", example.examplePath, valErr))
				return
			}

			// Handle example array.
			if len(result.Errors()) > 0 && strings.Contains(result.Errors()[0].String(), arrayError) {

				exampleLoaders, err := unpackArray(*exampleLoader)
				if err != nil {
					errors = append(errors, fmt.Errorf("%s: %s", example.examplePath, err))
					return
				}

				for _, exampleLoader := range exampleLoaders {
					result, valErr := gojsonschema.Validate(*exampleSchemaLoader, exampleLoader)
					if valErr != nil {
						errors = append(errors, fmt.Errorf("%s: %s", example.examplePath, valErr))
						return
					}

					for _, err := range result.Errors() {
						errors = append(errors, fmt.Errorf("%s: %s", example.examplePath, err.String()))
					}
				}
			}

			for _, err := range result.Errors() {
				if !strings.Contains(result.Errors()[0].String(), arrayError) {
					errors = append(errors, fmt.Errorf("%s: %s", example.examplePath, err.String()))
				}
			}
		})
	})
	return errors
}

// `gojsonschema` doesn't handle reference paths that point to specific fields like:
// file://aaa/bbb.json#definitions/example/something
// It will return the root object of that JSON file, completely ignoring the part after #
// As a workaround, this function loads that file, calls a JsonPath query to get the referenced object.
// Also, gojsonschema is not able to resolve schema references if a JSON object is not saved on disk.
// *Side effects*
// If arefPath contains the link to an object, it will generate a file on your file system.
// Related: https://github.com/xeipuuv/gojsonschema/issues/262
func GetReferenceLoader(refPath string) (*gojsonschema.JSONLoader, error) {

	simpleLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", refPath))
	if !strings.Contains(refPath, "#") {
		return &simpleLoader, nil
	}

	pathParts := strings.Split(refPath, "#")
	filePath := pathParts[0]
	objectPath := pathParts[1]

	if len(filePath) == 0 || len(objectPath) == 0 {
		return &simpleLoader, nil
	}

	jsonLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", filePath))
	jsonPath := TranslateReferenceToJSONPath(objectPath)
	jsonObj, refErr := jsonLoader.LoadJSON()
	if refErr != nil {
		return nil, fmt.Errorf("can't load the path: %s:%s", jsonPath, refErr)
	}

	foundObject, jsonPathErr := jsonpath.Get(jsonPath, jsonObj)
	if jsonPathErr != nil {
		return nil, fmt.Errorf("can't find the path: %s:%s %+v ", jsonPath, jsonPathErr, jsonObj)
	}

	partialLoader := gojsonschema.NewGoLoader(foundObject)
	return &partialLoader, nil
}

// unpack examples array into separate objects
func unpackArray(loader gojsonschema.JSONLoader) ([]gojsonschema.JSONLoader, error) {
	obj, err := loader.LoadJSON()
	if err != nil {
		return nil, err
	}

	objArray, ok := obj.([]interface{})
	if !ok {
		return nil, fmt.Errorf("can't unpack example array")
	}

	loaders := make([]gojsonschema.JSONLoader, 0, len(objArray))
	for _, elem := range objArray {
		loaders = append(loaders, gojsonschema.NewGoLoader(elem))
	}

	return loaders, nil
}

// Translate a OpenAPI reference path ($ref) to a JsonPath query in a dialect accepted by the `jsonpath` lib.
// e.g.
// /definitions/request -> $.definitions.request
// /definitions/request/200 -> $.definitions.request["200"]
// `jsonpath` throws an error if you try to access fields like this:
// `$.definitions.request.200`
// Related: https://github.com/PaesslerAG/jsonpath/issues/23
func TranslateReferenceToJSONPath(refPath string) string {
	var jsonPath []string
	for _, part := range strings.Split(refPath, "/") {
		if _, err := strconv.Atoi(part); err == nil {
			jsonPath = append(jsonPath, fmt.Sprintf("[\"%s\"]", part))
		} else {
			if len(part) > 0 {
				jsonPath = append(jsonPath, ".", part)
			}
		}
	}
	return fmt.Sprintf("$%s", strings.Join(jsonPath, ""))
}
