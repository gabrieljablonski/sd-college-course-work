package utils

import "fmt"

type IP struct {
	Address string
	Port    string
}

func (ip IP) String() string {
	return fmt.Sprintf("%s:%s", ip.Address, ip.Port)
}

func CheckKeys(m map[string]interface{}, keys []string) string {
	for _, key := range keys {
		if m[key] == nil {
			return key
		}
	}
	return ""
}
