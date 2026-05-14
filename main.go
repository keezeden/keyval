package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "[ERR] Failed to read body", http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		request := string(bodyBytes)

		segments := strings.Split(request, " ")

		if len(segments) == 0 {
			http.Error(w, "[ERR] Request must contain a command", http.StatusInternalServerError)
			return
		}

		command := segments[0]

		switch command {
		case "SET":
			// SET key val
			err := checkArgs(segments, 3)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			key, val := segments[1], segments[2]

			// fmt.Printf("Setting %s to %s\n", key, val)
			Set(key, val)

			return

		case "GET":
			// GET key
			err := checkArgs(segments, 2)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			key := segments[1]
			value, err := Get(key)

			if errors.Is(err, ErrNotFound) {
				http.Error(w, fmt.Sprintf("[ERR] %s", ErrNotFound.Error()), http.StatusNotFound)
				return
			}

			fmt.Fprintf(w, "[OK] %s", value)

			return

		case "DEL":
			// DEL key
			err := checkArgs(segments, 2)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			key := segments[1]

			fmt.Printf("Deleting %s", key)

			return

		case "EXISTS":
			//EXISTS key
			err := checkArgs(segments, 2)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			key := segments[1]

			fmt.Printf("Checking %s", key)

			return

		case "KEYS":
			// KEYS
			// no check needed

			return
		default:
			http.Error(w, "[ERR] Command not found", http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}

func checkArgs(segments []string, count int) error {
	if len(segments) < count {
		err := fmt.Sprintf("[ERR] Command expects at least %d arguments", count)
		return errors.New(err)
	}

	return nil
}
