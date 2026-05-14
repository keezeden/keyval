package main

import (
	"fmt"
	"os"
	"strings"
)

const STORAGE_FILE = "data.txt"

func Write(k string, v string) error {
	delimiter := "="
	cleanKey := strings.ReplaceAll(k, "=", "\\=")
	cleanValue := strings.ReplaceAll(v, "=", "\\=")
	formatted := fmt.Sprintf("%s%s%s\n", cleanKey, delimiter, cleanValue)
	// 0644 permissions (read/write for owner, read for others)
	f, err := os.OpenFile(STORAGE_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("storage.Write(%s, %s): %w", k, v, err)
	}

	defer f.Close()

	_, err = f.WriteString(formatted)

	if err != nil {
		return fmt.Errorf("storage.Write(%s, %s): %w", k, v, err)
	}

	return nil
}

func Read() (string, error) {
	data, err := os.ReadFile(STORAGE_FILE)
	if err != nil {
		return "", fmt.Errorf("storage.Read(): %w", err)
	}

	return string(data), nil
}
