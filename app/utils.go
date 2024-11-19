package main

import "encoding/json"

func SliceContains(slice []rune, target rune) bool {
    for _, item := range slice {
        if target == item {
            return true
        }
    }
    return false
}

// RecursiveUnmarshal unmarshals JSON data into a given interface{} recursively.
func RecursiveUnmarshal(data []byte, v interface{}) error {
    err := json.Unmarshal(data, v)
    if err != nil {
        return err
    }

    switch vt := v.(type) {
    case map[string]interface{}:
        for key, value := range vt {
            if raw, ok := value.(json.RawMessage); ok {
                var inner interface{}
                err := RecursiveUnmarshal(raw, &inner)
                if err != nil {
                    return err
                }
                vt[key] = inner
            }
        }
    case []interface{}:
        for i, value := range vt {
            if raw, ok := value.(json.RawMessage); ok {
                var inner interface{}
                err := RecursiveUnmarshal(raw, &inner)
                if err != nil {
                    return err
                }
                vt[i] = inner
            }
        }
    }

    return nil
}
