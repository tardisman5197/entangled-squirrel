package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

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

func GetSquirrels(path string) (map[string]int, error) {
	squirrels := make(map[string]int, 0)
	data, err := os.ReadFile(path)
	if err != nil {
		return squirrels, fmt.Errorf("could not read in squirrels, got %v", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			url := strings.TrimSpace(line)
			if url[0] != '#' {
				squirrels[url] = 0
			}
		}
	}

	return squirrels, nil
}

func WriteSquirrels(path string, squirrels map[string]int) error {
	data := ""
	for squirrel := range squirrels {
		data += squirrel + "\n"
	}
	err := os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("could not write squirrels, got %v", err)
	}
	return nil
}
