package utils

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

/**
 * Validate variables against JSON schema
 * @param variables Map of variables to validate
 * @param schema JSON schema definition
 * @return Error if validation fails
 */
func ValidateVariables(variables map[string]interface{}, schema interface{}) error {
	// Convert input variables to JSON document
	variablesJSON, err := json.Marshal(variables)
	if err != nil {
		return fmt.Errorf("failed to marshal variables: %w", err)
	}

	// Convert schema to JSON document
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	// Load schema
	schemaLoader := gojsonschema.NewBytesLoader(schemaJSON)
	documentLoader := gojsonschema.NewBytesLoader(variablesJSON)

	// Perform validation
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if !result.Valid() {
		var errs []string
		for _, desc := range result.Errors() {
			errs = append(errs, desc.String())
		}
		return fmt.Errorf("invalid variables: %v", errs)
	}

	return nil
}

/**
 * Validate arguments array against JSON schema
 * @param args Arguments array to validate
 * @param schema JSON schema definition
 * @return Error if validation fails
 */
func ValidateArgs(args []interface{}, schema interface{}) error {
	// Convert input args to JSON document
	variablesJSON, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("failed to marshal args: %w", err)
	}

	// Convert schema to JSON document
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	// Load schema
	schemaLoader := gojsonschema.NewBytesLoader(schemaJSON)
	documentLoader := gojsonschema.NewBytesLoader(variablesJSON)

	// Perform validation
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if !result.Valid() {
		var errs []string
		for _, desc := range result.Errors() {
			errs = append(errs, desc.String())
		}
		return fmt.Errorf("invalid args: %v", errs)
	}

	return nil
}
