package parse

import (
    "encoding/json"
    "github.com/dlclark/regexp2"
    "regexp"
    "strings"
)

// Finds all quotes, with any amount of escaped slashes
var targetQuotes = regexp.MustCompile(`\\*"`)
var targetSingleQuotes = regexp.MustCompile(`\\*'`)

// This targets keys (fully quoted), that are missing a comma followed by white space before them.
//FIXME - I can't figure out how to accurately target the whitespace as optional.
var targetMissingCommasBeforeKeysWithWhiteSpace = regexp2.MustCompile(`(?m)(?<!,)\s(^|\")([a-zA-Z0-9_\.]+)\"\s*:`, 0)

var targetUnquotedKeys = regexp.MustCompile(`(?m)(^|[{\s,])([a-zA-Z0-9_\.]+)\s*:`)
var targetUnquotedValues = regexp.MustCompile(`:\s*([a-zA-Z0-9_\-:.TZ\s]+)(\s*[,}\]])`)

// QUESTION - with the inital destruction of the escaping, what about nested escaped strings?
//  - We definitely want to consider them part of the JSON doc (too bad that choice for now, maybe a switch later?),
//  but would that work if I just blindly replace all the escaped strings?
func Parse(input string) (string, error) {
    // Strip out escaping
    filtered := targetSingleQuotes.ReplaceAllString(input, "\"")

    // Step 1: Fix unquoted keys using a regex
    // Matches keys that are unquoted (e.g., key: "value") and ensures they are quoted.
    filtered = targetUnquotedKeys.ReplaceAllString(filtered, `$1"$2":`)

    // Step 2: Add commas after every key-value pair, if missing
    // Matches key-value pairs that are not followed by a comma or a closing brace/bracket.
    // Start at index one to make sure we don't hit the very first key.
    // FIXME - there's a bug here where you have stuff like `"key": [` or `"key": {` on a line on its own.
    //  It'll screw up and put commas at the start, like `, "key": [`
    filtered, _ = targetMissingCommasBeforeKeysWithWhiteSpace.ReplaceFunc(filtered, addCommas, 1, -1)

    // Step 2: Fix unquoted string values using a regex
    // Matches unquoted values that are not numbers, booleans, null, or JSON objects/arrays.
    filtered = targetUnquotedValues.ReplaceAllStringFunc(filtered, func(match string) string {
        // Extract the value and ensure it requires quoting
        colonIndex := strings.Index(match, ":")
        value := strings.TrimSpace(match[colonIndex+1:])
        value = strings.Trim(value, ",")
        if value == "true" || value == "false" || value == "null" || isNumber(value) || isJSONStructure(value) {
            return match
        }
        return match[:colonIndex+1] + ` "` + value + `"` + match[len(match)-1:]
    })

    // If we're an array, you're going to have to have said that with the square brackets in the input,
    //otherwise if we just have key: values, we need to make sure we are wrapped in curly braces.
    if !strings.HasPrefix(filtered, "[") && !strings.HasPrefix(filtered, "{") {
        filtered = "{" + filtered
    }
    if !strings.HasSuffix(filtered, "[") && !strings.HasSuffix(filtered, "{") {
        filtered = filtered + "}"
    }

    // Step 4: Remove the trailing comma from the final key-value pair in each block
    // Matches a trailing comma before a closing brace or bracket.
    reTrailingComma := regexp.MustCompile(`,(\s*[}\]])`)
    filtered = reTrailingComma.ReplaceAllString(filtered, `$1`)

    data, err := RecursiveUnmarshal(filtered)
    if err != nil {
        return filtered, nil
    }

    result, err := json.MarshalIndent(data, "", "    ")
    if err != nil {
        return filtered, nil
    }
    return string(result), nil
}

// isNumber checks if a value is a valid JSON number.
func isNumber(value string) bool {
    // A regex to match integers or floating-point numbers
    reNumber := regexp.MustCompile(`^-?\d+(\.\d+)?([eE][+-]?\d+)?$`)
    return reNumber.MatchString(value)
}

// isJSONStructure checks if a value is a JSON object or array (e.g., starts with '{' or '[').
func isJSONStructure(value string) bool {
    return strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[")
}

func addCommas(match regexp2.Match) string {
    return "," + match.String()
}

// RecursiveUnmarshal takes a JSON byte array and recursively unmarshals it into a nested map or slice.
func RecursiveUnmarshal(data string) (interface{}, error) {
    var result interface{}

    // Attempt to unmarshal the JSON into an empty interface
    err := json.Unmarshal([]byte(data), &result)
    if err != nil {
        return nil, err
    }

    // Process the unmarshaled data recursively
    result = processRecursively(result)

    return result, nil
}

// processRecursively handles maps and slices recursively to ensure all values are processed.
func processRecursively(input interface{}) interface{} {
    switch value := input.(type) {
    case map[string]interface{}: // Process a JSON object
        value = handleDotNotation(value)
        for key, val := range value {
            mapVal, ok := val.(map[string]interface{})
            if ok {
                value[key] = processRecursively(mapVal)
            }
            arrVal, ok := val.([]interface{})
            if ok {
                value[key] = processRecursively(arrVal)
            }
            strVal, ok := val.(string)
            if ok {
                if strings.HasPrefix(strVal, "{") {
                    var tmp map[string]interface{}
                    err := json.Unmarshal([]byte(strVal), &tmp)
                    if err == nil {
                        result := processRecursively(tmp)
                        value[key] = result
                    }
                }
                if strings.HasPrefix(strVal, "[") {
                    var tmp []interface{}
                    err := json.Unmarshal([]byte(strVal), &tmp)
                    if err == nil {
                        result := processRecursively(tmp)
                        value[key] = result
                    }
                }
            }
        }
        return value
    case []interface{}: // Process a JSON array
        for i, val := range value {
            mapVal, ok := val.(map[string]interface{})
            if ok {
                value[i] = processRecursively(mapVal)
            }
            arrVal, ok := val.([]interface{})
            if ok {
                value[i] = processRecursively(arrVal)
            }
            strVal, ok := val.(string)
            if ok {
                if strings.HasPrefix(strVal, "{") {
                    var tmp map[string]interface{}
                    err := json.Unmarshal([]byte(strVal), &tmp)
                    if err == nil {
                        result := processRecursively(tmp)
                        value[i] = result
                    }
                }
                if strings.HasPrefix(strVal, "[") {
                    var tmp []interface{}
                    err := json.Unmarshal([]byte(strVal), &tmp)
                    if err == nil {
                        result := processRecursively(tmp)
                        value[i] = result
                    }
                }
            }
        }
        return value
    default: // Base case: return the value as is for primitive types
        return value
    }
}

func handleDotNotation(data map[string]interface{}) map[string]interface{} {
    result := map[string]interface{}{}
    for key, value := range data {
        if key == "" {
            continue
        }
        setNestedValue(result, key, value)
    }
    return result
}

func setNestedValue(data map[string]interface{}, path string, value interface{}) {
    keys := strings.Split(path, ".")
    current := data

    for _, key := range keys[:len(keys)-1] {
        if _, ok := current[key]; !ok {
            current[key] = make(map[string]interface{})
        }
        current = current[key].(map[string]interface{})
    }

    current[keys[len(keys)-1]] = value
}
