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

func ParseDevice(userAgent string) (deviceName string, deviceType string) {

	ua := strings.ToLower(userAgent)

	// detect device type
	if strings.Contains(ua, "mobile") {
		deviceType = "mobile"
	} else if strings.Contains(ua, "tablet") {
		deviceType = "tablet"
	} else {
		deviceType = "desktop"
	}

	// detect OS / device name
	switch {
	case strings.Contains(ua, "iphone"):
		deviceName = "iPhone"
	case strings.Contains(ua, "android"):
		deviceName = "Android"
	case strings.Contains(ua, "windows"):
		deviceName = "Windows PC"
	case strings.Contains(ua, "mac os"):
		deviceName = "Mac"
	case strings.Contains(ua, "linux"):
		deviceName = "Linux"
	default:
		deviceName = "Unknown"
	}

	return
}
