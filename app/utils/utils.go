package utils

import (
    "regexp"
    "strings"
)

// isNumber checks if a value is a valid JSON number.
func IsNumber(value string) bool {
    // A regex to match integers or floating-point numbers
    reNumber := regexp.MustCompile(`^-?\d+(\.\d+)?([eE][+-]?\d+)?$`)
    return reNumber.MatchString(value)
}

// isJSONStructure checks if a value is a JSON object or array (e.g., starts with '{' or '[').
func IsJSONStructure(value string) bool {
    return strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[")
}

func IsComplexObject(value string) bool {
    return strings.Contains(value, ":") || IsJSONStructure(value)
}
