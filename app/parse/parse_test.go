package parse

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse_success_valid_json_dict(t *testing.T) {
	input := `{"Iam": "validJSON"}`
	expected := `{
    "Iam": "validJSON"
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_valid_json_array(t *testing.T) {
	input := `[1, "test", true, null]`
	expected := `[
    1,
    "test",
    true,
    null
]`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_non_json_string(t *testing.T) {
	input := `"test"`
	expected := `"test"`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_non_json_int(t *testing.T) {
	input := `1`
	expected := `1`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_non_json_boolean(t *testing.T) {
	input := `true`
	expected := `true`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_unpacks_dot_notation_into_sub_docs(t *testing.T) {
	input := `{"nested.field1": "value1", "nested.field2": 2, "nested.field3": true, "nested.field4": null}`
	expected := `{
    "nested": {
        "field1": "value1",
        "field2": 2,
        "field3": true,
        "field4": null
    }
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_unpacks_json_strings(t *testing.T) {
	input := `"{\"Iam\": \"validJSON\"}"`
	expected := `{
    "Iam": "validJSON"
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_unpacks_nested_string_json_objects(t *testing.T) {
	input := `{"field1": "[\"test\", 1, false, null]", "nested.field2": "{\"nestedKey1\": 1, \"nestedKey2\": \"test\", \"nestedKey3\": true, \"nestedKey4\": null}"}`
	expected := `{
    "field1": [
        "test",
        1,
        false,
        null
    ],
    "nested": {
        "field2": {
            "nestedKey1": 1,
            "nestedKey2": "test",
            "nestedKey3": true,
            "nestedKey4": null
        }
    }
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_handles_missing_commas(t *testing.T) {
	input := `
response_body.line1: BAD_REQUEST
response_body.line2: "test"
response_body.line3: null
response_body.line4: true
response_body.line5: 1
`
	expected := `{
    "response_body": {
        "line1": "BAD_REQUEST",
        "line2": "test",
        "line3": null,
        "line4": true,
        "line5": 1
    }
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_handles_missing_commas_with_mid_value_new_lines(t *testing.T) {
	input := `
response_body.line1: BAD_
REQUEST
response_body.line2: "te
st"
response_body.line3: nul
l
response_body.line4: 
true
response_body.line5: 1
`
	expected := `{
    "response_body": {
        "line1": "BAD_REQUEST",
        "line2": "test",
        "line3": null,
        "line4": true,
        "line5": 1
    }
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_strips_leading_trailing_spaces(t *testing.T) {
	input := `    "  test "      `
	expected := `"  test "`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_strips_leading_trailing_quotes(t *testing.T) {
	input := `"{"IWantToBe": "validJSON"}"`
	expected := `{
    "IWantToBe": "validJSON"
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_replaces_multi_escaped_quotes(t *testing.T) {
	input := `"{\\\\\\"IWantToBe\": \\"validJSON\\\\\\\\\\\\\"}"`
	expected := `{
    "IWantToBe": "validJSON"
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_replaces_single_quotes(t *testing.T) {
	input := `{'IWantToBe':'validJSON'}`
	expected := `{
    "IWantToBe": "validJSON"
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_adds_quotes_for_unquoted_keys(t *testing.T) {
	input := `{IWantToBe:'validJSON'}`
	expected := `{
    "IWantToBe": "validJSON"
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_adds_quotes_for_keys_and_values(t *testing.T) {
	input := `{IWantToBe: validJSON}`
	expected := `{
    "IWantToBe": "validJSON"
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_simple_nested_structures(t *testing.T) {
	input := `{
"IWantToBe": "validJSON", 
"array": [1, [2, "2test"], 5],
"dict": {dict1: {"2key1": true }}
}
`
	expected := `{
    "IWantToBe": "validJSON",
    "array": [
        1,
        [
            2,
            "2test"
        ],
        5
    ],
    "dict": {
        "dict1": {
            "2key1": true
        }
    }
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_complex_nested_structures(t *testing.T) {
	input := `{
    "IWantToBe": "validJSON", 
    "array": [1, "test", true, false, null, [2, "2test", false, null, [3], 4], 5],
    "dict": {dict1: {"2key1": true, "2key2": false, "2key3": test, 2key4: null, dict2: {"3key1": false, "3key2": "2test", 3key3: null}, array2: [1, 2, 3]}, key1: true, key2: null, key3: "test", key4: 1}
    }
`

	expected := `{
    "IWantToBe": "validJSON",
    "array": [
        1,
        "test",
        true,
        false,
        null,
        [
            2,
            "2test",
            false,
            null,
            [
                3
            ],
            4
        ],
        5
    ],
    "dict": {
        "dict1": {
            "2key1": true,
            "2key2": false,
            "2key3": "test",
            "2key4": null,
            "array2": [
                1,
                2,
                3
            ],
            "dict2": {
                "3key1": false,
                "3key2": "2test",
                "3key3": null
            }
        },
        "key1": true,
        "key2": null,
        "key3": "test",
        "key4": 1
    }
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_adds_quotes_for_unquoted_and_misquotes_string_values(t *testing.T) {
	input := `{
"IWantToBe": validJSON,
"NonStringField1": 1,
"NonStringField2": true,
"NonStringField3": null,
"NonStringField4": [1, tmp", false, null],
"NonStringField5": {"key1": 1, "key2": tmp, "key3: false, "key4": null],
}`
	expected := `{
    "IWantToBe": "validJSON",
    "NonStringField1": 1,
    "NonStringField2": true,
    "NonStringField3": null,
    "NonStringField4": [
        1,
        "tmp",
        false,
        null
    ],
    "NonStringField5": {
        "key1": 1,
        "key2": "tmp",
        "key3": false,
        "key4": null
    }
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_fixes_bad_data_structure_suffixes(t *testing.T) {
	input := `{
"array": [1, tmp", false, null},
"dict": {"key1": 1, "key2": tmp, "key3: false, "key4": null],
}`
	expected := `{
    "array": [
        1,
        "tmp",
        false,
        null
    ],
    "dict": {
        "key1": 1,
        "key2": "tmp",
        "key3": false,
        "key4": null
    }
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_adds_leading_trailing_curly_brackets_if_not_present(t *testing.T) {
	input := `IWantToBe:'validJSON'`
	expected := `{
    "IWantToBe": "validJSON"
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_strips_trailing_commas(t *testing.T) {
	// Root and nested
	input := `{
"IWantToBe": "validJSON",
"array": [1, [2, "2test"],],
"dict": {dict1: {"2key1": true, },},
}`
	expected := `{
    "IWantToBe": "validJSON",
    "array": [
        1,
        [
            2,
            "2test"
        ]
    ],
    "dict": {
        "dict1": {
            "2key1": true
        }
    }
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParse_success_strips_extra_white_spaces(t *testing.T) {
	// Root and nested
	input := "{\"test\": \"1\",\n\n\n\n\"test2\"\t\t\n: 2,\t\t\t\t\"test3\"\n\n\t: true}"
	expected := `{
    "test": 1,
    "test2": 2,
    "test3": true
}`

	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
