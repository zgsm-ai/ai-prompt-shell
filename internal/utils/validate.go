package utils

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateVariables 根据 JSON schema 验证 variables 是否合法
func ValidateVariables(variables map[string]interface{}, schema interface{}) error {
	// 将输入 variables 转换为 JSON 文档
	variablesJSON, err := json.Marshal(variables)
	if err != nil {
		return fmt.Errorf("failed to marshal variables: %w", err)
	}

	// 将 schema 转换为 JSON 文档
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	// 加载 schema
	schemaLoader := gojsonschema.NewBytesLoader(schemaJSON)
	documentLoader := gojsonschema.NewBytesLoader(variablesJSON)

	// 执行验证
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

func ValidateArgs(args []interface{}, schema interface{}) error {
	// 将输入 variables 转换为 JSON 文档
	variablesJSON, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("failed to marshal variables: %w", err)
	}

	// 将 schema 转换为 JSON 文档
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	// 加载 schema
	schemaLoader := gojsonschema.NewBytesLoader(schemaJSON)
	documentLoader := gojsonschema.NewBytesLoader(variablesJSON)

	// 执行验证
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
