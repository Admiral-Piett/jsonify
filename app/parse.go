package main

import (
    "encoding/json"
    "fmt"
    "regexp"
    "strings"
)

var invalidRunes = []rune("'\"\\',")

var digitsOnlyFilter = regexp.MustCompile(`^[0-9]+$`)

func prepText(text string) string {
    text = strings.TrimSpace(text)
    if !strings.HasPrefix(text, "{") {
        text = "{" + text
    }
    if !strings.HasSuffix(text, "}") {
        if strings.HasSuffix(text, ",") {
            text = text[:len(text)-1]
        }
        text = text + "}"
    }
    return text
}

func sanitize(text string) string {
    textRune := []rune(text)

    filteredRunes := []rune{}
    for _, r := range textRune {
        if SliceContains(invalidRunes, r) {
            continue
        }
        filteredRunes = append(filteredRunes, r)
    }
    resultString := string(filteredRunes)
    if !digitsOnlyFilter.MatchString(resultString) && resultString != "null" && resultString != "true" && resultString != "false" {
        resultString = "\"" + resultString + "\""
    }
    return resultString
}

// TAG - this won't work for [{}] or any other kind of array.
func correctInvalidValues(text string) (string, error) {
    formatedLines := []string{}
    splitText := strings.Split(text, "\n")

    lastLine := len(splitText) - 1
    for i, line := range splitText {
        splitLine := strings.SplitN(line, ":", 2)
        if 2 < len(splitLine) {
            return "", fmt.Errorf("Invalid line: %s", line)
        }

        key := strings.TrimSpace(splitLine[0])
        value := strings.TrimSpace(splitLine[1])

        key = strings.TrimLeft(key, "{")
        value = strings.TrimRight(value, ",")
        value = strings.TrimRight(value, "}")
        value = strings.Trim(value, "\"")

        key = sanitize(key)
        // This is best effort for now.  If we detect these characters at the start of these strings, we
        //  will assume it's a complex data type and do our best to sort its punctuation out before trying to push
        //  it out.
        if strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[") {
            value = strings.ReplaceAll(value, "\\\"", "\"")
            value = strings.ReplaceAll(value, "\\'", "\"")
            value = strings.ReplaceAll(value, "'", "\"")
        } else {
            value = sanitize(value)
        }

        formattedLine := key + ":" + value
        if i != lastLine {
            formattedLine = formattedLine + ","
        }
        formatedLines = append(formatedLines, formattedLine)
    }
    return strings.Join(formatedLines, "\n"), nil
}

func parse(text string) (string, error) {
    data := &map[string]interface{}{}

    text = prepText(text)

    err := json.Unmarshal([]byte(text), data)
    if err != nil {
        text, err = correctInvalidValues(text)
        // There are limitations with Go's regex, so we can't be as exact as we need to be.  We will make sure our
        //  overall document is in a good state after our refactoring.
        text = prepText(text)
        if err != nil {
            return "", err
        }
        err = json.Unmarshal([]byte(text), data)
        if err != nil {
            return "", err
        }
    }

    result := make(map[string]interface{})
    for path, value := range *data {
        if path == "" {
            continue
        }
        setNestedValue(result, path, value)
    }

    jsonData, err := json.MarshalIndent(result, "", "  ")
    if err != nil {
        return "", err
    }
    jsonStr := string(jsonData)
    return jsonStr, nil
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
