package utils

import "strings"

func TrimValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case map[string]interface{}:
		return TrimMap(v)
	case []interface{}:
		return TrimSlice(v)
	default:
		return value
	}
}

func TrimMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		result[k] = TrimValue(v)
	}
	return result
}

func TrimSlice(arr []interface{}) []interface{} {
	result := make([]interface{}, len(arr))
	for i, v := range arr {
		result[i] = TrimValue(v)
	}
	return result
}
