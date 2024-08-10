package gonuts

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// JSONPathExtractor extracts values from JSON data using a path-like syntax
type JSONPathExtractor struct {
	data interface{}
}

// NewJSONPathExtractor creates a new JSONPathExtractor
//
// Example:
//
//	jsonData := `{"name": "John", "age": 30, "address": {"city": "New York"}}`
//	extractor, err := NewJSONPathExtractor(jsonData)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewJSONPathExtractor(jsonData string) (*JSONPathExtractor, error) {
	var data interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &JSONPathExtractor{data: data}, nil
}

// Extract retrieves a value from the JSON data using the given path
//
// The path can include dot notation for nested objects, bracket notation for array indices,
// and wildcards (*) for matching multiple elements.
//
// Example:
//
//	value, err := extractor.Extract("address.city")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(value) // Output: New York
func (jpe *JSONPathExtractor) Extract(path string) (interface{}, error) {
	parts := strings.FieldsFunc(path, func(r rune) bool {
		return r == '.' || r == '[' || r == ']'
	})

	var current interface{} = jpe.data
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			if part == "*" {
				return jpe.handleWildcard(v, parts[len(parts)-1])
			}
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil, fmt.Errorf("key not found: %s", part)
			}
		case []interface{}:
			if part == "*" {
				return jpe.handleArrayWildcard(v, parts[len(parts)-1])
			}
			index, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid array index: %s", part)
			}
			if index < 0 || index >= len(v) {
				return nil, fmt.Errorf("array index out of bounds: %d", index)
			}
			current = v[index]
		default:
			return nil, fmt.Errorf("cannot navigate further from %T", v)
		}
	}

	return current, nil
}

// handleWildcard processes wildcard matching for objects
func (jpe *JSONPathExtractor) handleWildcard(obj map[string]interface{}, lastPart string) (interface{}, error) {
	result := make(map[string]interface{})
	for key, value := range obj {
		if lastPart == "*" {
			result[key] = value
		} else if nestedObj, ok := value.(map[string]interface{}); ok {
			if nestedValue, ok := nestedObj[lastPart]; ok {
				result[key] = nestedValue
			}
		}
	}
	return result, nil
}

// handleArrayWildcard processes wildcard matching for arrays
func (jpe *JSONPathExtractor) handleArrayWildcard(arr []interface{}, lastPart string) (interface{}, error) {
	var result []interface{}
	for _, item := range arr {
		if lastPart == "*" {
			result = append(result, item)
		} else if obj, ok := item.(map[string]interface{}); ok {
			if value, ok := obj[lastPart]; ok {
				result = append(result, value)
			}
		}
	}
	return result, nil
}

// ExtractString is a convenience method that extracts a string value
//
// Example:
//
//	name, err := extractor.ExtractString("name")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(name) // Output: John
func (jpe *JSONPathExtractor) ExtractString(path string) (string, error) {
	v, err := jpe.Extract(path)
	if err != nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("value at path %s is not a string", path)
	}
	return s, nil
}

// ExtractInt is a convenience method that extracts an int value
func (jpe *JSONPathExtractor) ExtractInt(path string) (int, error) {
	v, err := jpe.Extract(path)
	if err != nil {
		return 0, err
	}
	switch n := v.(type) {
	case float64:
		return int(n), nil
	case int:
		return n, nil
	default:
		return 0, fmt.Errorf("value at path %s is not a number", path)
	}
}

// ExtractFloat is a convenience method that extracts a float64 value
func (jpe *JSONPathExtractor) ExtractFloat(path string) (float64, error) {
	v, err := jpe.Extract(path)
	if err != nil {
		return 0, err
	}
	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("value at path %s is not a float", path)
	}
	return f, nil
}

// ExtractBool is a convenience method that extracts a bool value
func (jpe *JSONPathExtractor) ExtractBool(path string) (bool, error) {
	v, err := jpe.Extract(path)
	if err != nil {
		return false, err
	}
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("value at path %s is not a boolean", path)
	}
	return b, nil
}

// ExtractArray is a convenience method that extracts an array value
func (jpe *JSONPathExtractor) ExtractArray(path string) ([]interface{}, error) {
	v, err := jpe.Extract(path)
	if err != nil {
		return nil, err
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value at path %s is not an array", path)
	}
	return arr, nil
}

// ExtractMultiple extracts multiple values from the JSON data using the given paths
//
// Example:
//
//	values, err := extractor.ExtractMultiple("name", "age", "address.city")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(values) // Output: map[name:John age:30 address.city:New York]
func (jpe *JSONPathExtractor) ExtractMultiple(paths ...string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, path := range paths {
		value, err := jpe.Extract(path)
		if err != nil {
			return nil, fmt.Errorf("error extracting path %s: %w", path, err)
		}
		result[path] = value
	}
	return result, nil
}
