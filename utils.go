package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func mergeMap(a map[string]string, b map[string]string) map[string]string {
	newMap := make(map[string]string)
	for k, v := range a {
		newMap[k] = v
	}
	for k, v := range b {
		newMap[k] = v
	}
	return newMap
}

func askYesOrNo(reader io.Reader) bool {
	line, err := readline(reader)
	if err != nil {
		return false
	}
	return strings.HasPrefix(strings.ToUpper(strings.Trim(line, " ")), "Y")
}

func readline(reader io.Reader) (value string, err error) {
	var valb []byte
	var n int
	b := make([]byte, 1)
	for {
		n, err = reader.Read(b)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 || b[0] == '\n' {
			break
		}
		valb = append(valb, b[0])
	}

	return strings.TrimSuffix(string(valb), "\r"), nil
}

func sortKeys(m map[string]map[string]string) []string {
	keys := []string{}
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("Not found env key: %v", key)
	}
	return v, nil
}
