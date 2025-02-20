package parse

import (
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
	case '\n' == value:
		return true
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
