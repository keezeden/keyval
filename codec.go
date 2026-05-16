package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func DecodeEvents(contents string) (map[string]Event, error) {
	lines := strings.Split(contents, "\n")
	data := make(map[string]Event)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var e Event
		err := json.Unmarshal([]byte(line), &e)
		if err != nil {
			return nil, fmt.Errorf("DecodeEvents(): %w", err)
		}

		switch e.Type {
		case "SET":
			data[e.Key] = e
			continue
		case "DEL":
			delete(data, e.Key)
			continue
		case "EXP":
			delete(data, e.Key)
			continue
		default:
			continue
		}
	}

	return data, nil
}
