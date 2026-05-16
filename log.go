package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Log struct {
	Path string
}

func NewLog(path string) *Log {
	return &Log{
		Path: path,
	}
}

func (log *Log) ReadAll() (string, error) {
	bytes, err := os.ReadFile(log.Path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (log *Log) Append(event Event) error {
	file, err := os.OpenFile(log.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer file.Close()
	fmt.Printf("Key: %s, Value: %s\n", event.Key, event.Value)
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if _, err := file.Write(append(bytes, '\n')); err != nil {
		return err
	}

	return file.Sync()
}
