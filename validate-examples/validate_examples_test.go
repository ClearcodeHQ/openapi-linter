package validate_examples

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestFindExamples(t *testing.T) {
	Assert := assert.New(t)
	obj := map[string]interface{}{}

	findExamplesHelper := func(obj map[string]interface{}) ([]Example, []error) {
		examples := []Example{}
		errors := []error{}

		FindExamples(obj, func(example Example, err error) {
			examples = append(examples, example)
			errors = append(errors, err)
		})
		return examples, errors
	}

	t.Run("The object is an array", func(t *testing.T) {
		// GIVEN
		json.Unmarshal([]byte(`[1, 2, 3]`), &obj)

		// WHEN
		examples, errors := findExamplesHelper(obj)

		// THEN
		Assert.Len(errors, 0)
		Assert.Len(examples, 0)
	})

	t.Run("No schema key or example", func(t *testing.T) {
		// GIVEN
		json.Unmarshal([]byte(`{
			"root": {
				"foo": {
					"bar": "somedoc.json"
				}
			}		
		}`), &obj)
		expectedExamples := []Example{}
		expectedErrors := []error{}

		// WHEN
		examples, errors := findExamplesHelper(obj)

		// THEN
		Assert.Equal(examples, expectedExamples)
		Assert.Equal(errors, expectedErrors)
	})

	t.Run("Example without schema", func(t *testing.T) {
		// GIVEN
		json.Unmarshal([]byte(`{
			"root": {
				"example": {
					"$ref": "somedoc.json"
				}
			}		
		}`), &obj)
		expectedExamples := []Example{
			{},
		}
		expectedErrors := []error{
			fmt.Errorf("Can't find the schema of the example."),
		}

		// WHEN
		examples, errors := findExamplesHelper(obj)

		// THEN
		Assert.Equal(examples, expectedExamples)
		Assert.Equal(errors, expectedErrors)
	})

	t.Run("Example with the inline object", func(t *testing.T) {
		// GIVEN
		json.Unmarshal([]byte(`{
			"root": {
				"example": {
					"foo": "bar",
					"baz": "buzz"
				},
				"schema": { "$ref": "smth.json" }
			}		
		}`), &obj)
		expectedExamples := []Example{
			{},
		}
		expectedErrors := []error{
			fmt.Errorf("The reference to schema/example is missing, inline objects aren't supported."),
		}

		// WHEN
		examples, errors := findExamplesHelper(obj)

		// THEN
		Assert.Equal(examples, expectedExamples)
		Assert.Equal(errors, expectedErrors)

	})
	t.Run("Schema with the inline object", func(t *testing.T) {
		// GIVEN
		json.Unmarshal([]byte(`{
			"root": {
				"example": {
					"foo": "bar",
					"baz": "buzz"
				},
				"schema": {
					"foo": "bar",
					"baz": "buzz"
				}
			}		
		}`), &obj)
		expectedExamples := []Example{
			{},
		}
		expectedErrors := []error{
			fmt.Errorf("The reference to schema/example is missing, inline objects aren't supported."),
		}

		// WHEN
		examples, errors := findExamplesHelper(obj)

		// THEN
		Assert.Equal(examples, expectedExamples)
		Assert.Equal(errors, expectedErrors)
	})

	// Create parametrized test cases for the invalid json types.
	exampleInvalidObjectTestCases := []string{
		"[]",
		"null",
		"1",
	}
	for _, invalidObjectTestCase := range exampleInvalidObjectTestCases {
		t.Run(fmt.Sprintf("Example is not an object but: %s", invalidObjectTestCase), func(t *testing.T) {
			// GIVEN
			json.Unmarshal([]byte(fmt.Sprintf(`{
				"root": {
					"example": %s,
					"schema": { "$ref": "smth.json" }
				}		
			}`, invalidObjectTestCase)), &obj)
			expectedExamples := []Example{
				{},
			}
			expectedErrors := []error{
				fmt.Errorf("Can't cast the schema/example object to map[string]."),
			}
			// WHEN
			examples, errors := findExamplesHelper(obj)

			// THEN
			Assert.Equal(examples, expectedExamples)
			Assert.Equal(errors, expectedErrors)
		})

		t.Run(fmt.Sprintf("Schema is not an object but: %s", invalidObjectTestCase), func(t *testing.T) {
			// GIVEN
			json.Unmarshal([]byte(fmt.Sprintf(`{
				"root": {
					"example": { "$ref": "smth.json"},
					"schema":  %s
				}		
			}`, invalidObjectTestCase)), &obj)
			expectedExamples := []Example{
				{},
			}
			expectedErrors := []error{
				fmt.Errorf("Can't cast the schema/example object to map[string]."),
			}

			// WHEN
			examples, errors := findExamplesHelper(obj)

			// THEN
			Assert.Equal(examples, expectedExamples)
			Assert.Equal(errors, expectedErrors)
		})
	}

	t.Run("Example with the schema", func(t *testing.T) {
		// GIVEN
		json.Unmarshal([]byte(`{
			"root": {
				"example": { "$ref": "bb.json" },
				"schema": { "$ref": "aa.json" }
			}		
		}`), &obj)
		expectedExamples := []Example{
			{
				"aa.json",
				"bb.json",
			},
		}
		expectedErrors := []error{
			nil,
		}

		// WHEN
		examples, errors := findExamplesHelper(obj)

		// THEN
		Assert.Equal(examples, expectedExamples)
		Assert.Equal(errors, expectedErrors)
	})
}

// ScanForExamples is a function that's heavily IO based and it doesn't implement the dependency injection pattern.
func TestScanForExamples(t *testing.T) {
	Assert := assert.New(t)
	_, testsFile, _, _ := runtime.Caller(0)
	testsRoot := filepath.Dir(testsFile)

	getFixturesPath := func(fixtureName string) string {
		fixturePath := filepath.Join(testsRoot, "..", "tests", "validate_examples", "scan_for_examples", fixtureName)
		_, err := os.Stat(fixturePath)

		Assert.Nilf(err, fmt.Sprintf("Invalid fixture name or can't find the directory: %s", fixturePath))
		return fixturePath
	}

	t.Run("No examples, no errors", func(t *testing.T) {
		// WHEN
		errors := ScanForExamples(getFixturesPath("no_examples"))

		// THEN
		Assert.Len(errors, 0)

	})
	t.Run("Invalid example", func(t *testing.T) {
		// GIVEN
		expectedErrors := []error{
			fmt.Errorf("example.json#rootobject1/aaa/bbb: (root): exampleField is required"),
			fmt.Errorf("example.json#rootobject1/aaa/bbb: (root): Additional property xxx is not allowed"),
		}
		// WHEN
		errors := ScanForExamples(getFixturesPath("invalid_example"))

		// THEN
		Assert.Equal(expectedErrors, errors)
	})
	t.Run("Example without errors", func(t *testing.T) {
		// WHEN
		errors := ScanForExamples(getFixturesPath("without_errors"))

		// THEN
		Assert.Len(errors, 0)
	})
}

func TestGetReferenceLoader(t *testing.T) {
	Assert := assert.New(t)
	_, testsFile, _, _ := runtime.Caller(0)
	testsRoot := filepath.Dir(testsFile)

	getFixturesPath := func(fixtureName string) string {
		fixturePath := filepath.Join(testsRoot, "..", "tests", "validate_examples", "reference_loader", fixtureName)
		fixtureFileName := strings.Split(fixturePath, "#")
		_, err := os.Stat(fixtureFileName[0])

		Assert.Nilf(err, fmt.Sprintf("Invalid fixture name or can't find the directory: %s", fixturePath))
		return fixturePath
	}

	getReferenceLoaderHelper := func(fixtureName string) (map[string]interface{}, gojsonschema.JSONLoader) {
		fixturePath := getFixturesPath(fixtureName)
		loader, _ := GetReferenceLoader(fixturePath)
		obj, _ := loader.LoadJSON()
		objMap, _ := obj.(map[string]interface{})

		return objMap, loader
	}

	t.Run("Load file without the path reference", func(t *testing.T) {
		// GIVEN
		expectedObject := map[string]interface{}{
			"aaa": "bbb",
		}
		// WHEN
		extractedObject, _ := getReferenceLoaderHelper("basic.json")

		// THEN
		Assert.Equal(extractedObject, expectedObject)
	})

	t.Run("Load File with the object query", func(t *testing.T) {
		// GIVEN
		expectedObject := map[string]interface{}{
			"eee": "fff",
		}
		// WHEN
		extractedObject, partialPath := getReferenceLoaderHelper("nested.json#/aaa/ccc")

		// THEN
		Assert.Equal(extractedObject, expectedObject)

		RemovePartialLoader(partialPath)
	})
}

func TestTranslationOfReferencesToJsonPath(t *testing.T) {
	Assert := assert.New(t)

	t.Run("Empty reference", func(t *testing.T) {
		Assert.Equal("$", TranslateReferenceToJSONPath(""))
	})

	t.Run("Nested reference", func(t *testing.T) {
		Assert.Equal("$.documents.request", TranslateReferenceToJSONPath("/documents/request"))
	})

	t.Run("Reference contains a key expressed as a number", func(t *testing.T) {
		Assert.Equal("$.documents.request[\"200\"]", TranslateReferenceToJSONPath("/documents/request/200"))
		Assert.Equal("$.documents.request[\"200\"].headers", TranslateReferenceToJSONPath("/documents/request/200/headers"))
	})
}
