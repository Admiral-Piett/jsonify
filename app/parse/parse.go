package parse

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Also strip just lingering `\` characters at the start and end of the line.
// We may have skipped the `"` characters during processing.
var targetStartEndOfLineQuotes = regexp.MustCompile(`^\\*"|\\*"$|^\\*'|\\*'$`)
var targetStartEndResidualSlashes = regexp.MustCompile(`^\\+|\\+$`)

func Parse(input string) (string, error) {
	result := strings.TrimSpace(input)

	// NOTE: We could try to just marshal/unmarshall blind first.  If that doesn't work, we could correct.
	//  - We'd need - `stripUnneededQuotes` in `utils.go`

	if IsComplexObject(result) {
		// Trim these off the top since it's just going to throw us off later.
		result = targetStartEndOfLineQuotes.ReplaceAllString(result, "")
		result = strings.Trim(result, "{")
		result = strings.Trim(result, "}")
		result = strings.Trim(result, "[")
		result = strings.Trim(result, "]")

		runeSlice := []rune(result)
		isObj, err := determineObjectType(runeSlice)
		if err != nil {
			return "", err
		}

		filtered, _, err := correctInvalidFormatting(isObj, runeSlice)
		if err != nil {
			return "", err
		}
		result = filtered
	}

	data, err := RecursiveUnmarshal(result)
	if err != nil {
		return "", err
	}

	resultBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(resultBytes), nil
}

func determineObjectType(input []rune) (isObject bool, err error) {
	isObject = false
	err = nil

	firstChar := input[0]
	if firstChar == '{' {
		isObject = true
	} else if firstChar == '[' {
		isObject = false
	} else if containsRune(input, ':') {
		isObject = true
	} else if containsRune(input, ',') {
		isObject = false
	} else {
		err = fmt.Errorf("invalid complex data type found: %s", string(input))
	}

	return
}

func format(input string) string {
	result := strings.TrimSpace(input)
	result = targetStartEndOfLineQuotes.ReplaceAllString(result, "")
	result = targetStartEndResidualSlashes.ReplaceAllString(result, "")

	if result == "true" || result == "false" || result == "null" || IsNumber(result) || IsComplexObject(input) {
		return result
	}
	return fmt.Sprintf(`"%s"`, result)
}

func correctInvalidFormatting(isObj bool, input []rune) (result string, processedIndex int, err error) {
	result = ""
	// If we're just starting, we want to make sure we start processing normally,
	//  and don't skip in the index compare below.
	processedIndex = -1
	err = nil

	key := strings.Builder{}
	value := strings.Builder{}
	current := strings.Builder{}
	processedData := strings.Builder{}
	seenSemiColon := false

	for i, r := range input {
		// We may have processed recursively ahead, we need to skip that piece before we can start again.
		if i <= processedIndex || isSkippableCharacter(r) {
			continue
		}
		processedIndex = i

		if startsComplexDataStructure(r) {
			nestedInput := input[i+1:]
			nestedIsObj, e := determineObjectType([]rune{r})
			if e != nil {
				err = e
				return
			}
			nestedData, nestedIndex, e := correctInvalidFormatting(nestedIsObj, nestedInput)
			if e != nil {
				err = e
				return
			}
			current.WriteString(nestedData)
			// Start processing on the index after this.
			processedIndex += nestedIndex + 1
			continue
		} else if endsComplexDataStructure(r) {
			// We will ditch the opening tags by way of the above if condition.
			// Here, we'll skip the closing tags and stop the loop, so we can save and set the tags below.
			// We've already determined which ones they should be with `isObj`.
			break
		}
		if r == ':' {
			key.WriteString(format(current.String()))
			current = strings.Builder{}
			seenSemiColon = true
			continue
		}
		if r == ',' {
			value.WriteString(format(current.String()))
			processedData = writeLine(key, value, processedData)

			key = strings.Builder{}
			value = strings.Builder{}
			current = strings.Builder{}
			seenSemiColon = false
			continue
		}
		if r == '\n' {
			// If we haven't seen the middle marker of this line yet, then we can consider this cruft,
			// throw it out, and continue.
			if seenSemiColon == false {
				continue
			}

			// If the next "anchor" we find, isn't a `:` we will assume we're not yet at the end of the current
			// line's value.  That means we can't interpret this `\n` as a key/value pair break, and we need to
			// throw out this newline character and continue processing.
			// `i+1` is required to start searching after the current character which is already a `\n`
			nextAnchor := findNextAnchorCharacter(input[i+1:])
			if nextAnchor != ':' {
				continue
			}

			value.WriteString(format(current.String()))
			processedData = writeLine(key, value, processedData)

			key = strings.Builder{}
			value = strings.Builder{}
			current = strings.Builder{}
			seenSemiColon = false
			continue
		}
		current.WriteRune(r)
	}

	// Catch the end of the processing
	// NOTE: this is the last one, so we will take out the commas (that's why this isn't shared with the above saving)
	if current.Len() > 0 {
		value.WriteString(format(current.String()))
		processedData = writeLine(key, value, processedData)
	}

	result = processedData.String()
	// It's potentially possible to get through here with a residual `,` if one came in with the input.
	// If we have that, strip it before we wrap the data in tags.
	result = strings.TrimSuffix(result, ",")
	if isObj {
		result = "{" + result + "}"
	} else {
		result = "[" + result + "]"
	}
	return
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
