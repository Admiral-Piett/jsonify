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

//func TestParse_success_unpacks_nested_string_json_objects(t *testing.T) {
//    input := `{"field1": "[\"test\", 1, false, null]", "nested.field2": "{\"nestedKey1\": 1, \"nestedKey2\": \"test\", \"nestedKey3\": true, \"nestedKey4\": null}"}`
//    expected := `{
//    "field1": [
//        "test",
//        1,
//        false,
//        null
//    ],
//    "nested.field2": {
//        "nestedKey1": 1,
//        "nestedKey2": "test",
//        "nestedKey3": true,
//        "nestedKey4": null
//    }
//}`
//
//    result, err := Parse(input)
//
//    assert.Nil(t, err)
//    assert.Equal(t, expected, result)
//}

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

//func TestParse_success_adds_quotes_for_unquoted_string_values(t *testing.T) {
//    input := `{
//"IWantToBe": validJSON,
//"NonStringField1": 1,
//"NonStringField2": true,
//"NonStringField3": null,
//"NonStringField4": [1, "tmp", false, null],
//"NonStringField5": {"key1": 1, "key2": tmp, "key3": false, "key4": null],
//}`
//    expected := `{
//"IWantToBe": "validJSON",
//"NonStringField1": 1,
//"NonStringField2": true,
//"NonStringField3": null,
//"NonStringField4": [1, "tmp", false, null],
//"NonStringField5": {"key1": 1, "key2": :"tmp", "key3": false, "key4": null],
//}
//`
//
//    result, err := Parse(input)
//
//    assert.Nil(t, err)
//    assert.Equal(t, expected, result)
//}

func TestParse_success_adds_leading_trailing_curly_brackets_if_not_present(t *testing.T) {
    input := `IWantToBe:'validJSON'`
    expected := `{
    "IWantToBe": "validJSON"
}`

    result, err := Parse(input)

    assert.Nil(t, err)
    assert.Equal(t, expected, result)
}
