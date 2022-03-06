package utils

import "net/http"

func GetKeys(m map[string]int) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func SendFlash(url string) error {
	_, err := http.Get(url)
	return err
}
