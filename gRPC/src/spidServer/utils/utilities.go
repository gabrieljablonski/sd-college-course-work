package utils

type IP struct {
	Address string
	Port    string
}

func CheckKeys(m map[string]interface{}, keys []string) string {
	for _, key := range keys {
		if m[key] == nil {
			return key
		}
	}
	return ""
}
