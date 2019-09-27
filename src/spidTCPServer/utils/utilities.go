package utils

func CheckKeys(m map[string]interface{}, keys []string) string {
	for _, key := range keys {
		if m[key] == nil {
			return key
		}
	}
	return ""
}
