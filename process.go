package main

import (
	"fmt"
	"strings"
)

func Set(k string, v string) error {
	err := Write(k, v)
	if err != nil {
		return fmt.Errorf("process.Set(%s,%s): %w", k, v, err)
	}
	return nil
}

func Get(k string) (string, error) {
	data, err := Read()
	if err != nil {
		return "", fmt.Errorf("process.Get(): %w", err)
	}

	lines := strings.Split(data, "\n")
	table := make(map[string]string)
	for _, line := range lines {
		segments := strings.Split(line, "=")
		if len(segments) < 2 {
			continue
		}
		key, val := segments[0], segments[1]
		table[key] = val
	}

	value, ok := table[k]

	if !ok {
		return "", fmt.Errorf("process.Get(): %w", ErrNotFound)
	}

	return value, nil
}
