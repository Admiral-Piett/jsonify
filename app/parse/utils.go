package parse

import (
	"fmt"
	"regexp"
	"strings"
)

func startsComplexDataStructure(value rune) bool {
	return strings.ContainsRune("{[", value)
}

func endsComplexDataStructure(value rune) bool {
	return strings.ContainsRune("}]", value)
}

func isSkippableCharacter(value rune) bool {
	switch {
	case strings.ContainsRune(`"' `, value):
		return true
	//case '\n' == value:
	//	return true
	case '\t' == value:
		return true
	}
	return strings.ContainsRune(`"'`, value)
}

func IsNumber(value string) bool {
	// A regex to match integers or floating-point numbers
	reNumber := regexp.MustCompile(`^-?\d+(\.\d+)?([eE][+-]?\d+)?$`)
	return reNumber.MatchString(value)
}

func IsComplexObject(value string) bool {
	return strings.Contains(value, ":") || strings.Contains(value, ",")
}

func containsRune(runes []rune, r rune) bool {
	for _, v := range runes {
		if v == r {
			return true
		}
	}
	return false
}

func writeLine(key, value, processedData strings.Builder) strings.Builder {
	if key.Len() > 0 {
		processedData.WriteString(fmt.Sprintf("%s: %s,", key.String(), value.String()))
	} else {
		processedData.WriteString(fmt.Sprintf("%s,", value.String()))
	}
	return processedData
}

func findNextAnchorCharacter(data []rune) rune {
	for _, r := range data {
		switch r {
		case ':':
			return r
		case ',':
			return r
		case '\n':
			return r
		}
	}
	return 0
}

//func stripUnneededQuotes(input interface{}) interface{} {
//    dict, ok := input.(map[string]interface{})
//    if ok {
//        for k, v := range dict {
//
//        }
//    }
//    arr, ok := input.([]interface{})
//    if ok {
//
//    }
//    return input
//}
