package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	store := NewStore()
	err = store.Load()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("TCP server listening on port 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}

		go handleConnection(conn, store)
	}
}

func handleConnection(connection net.Conn, store *Store) {
	defer connection.Close()

	scanner := bufio.NewScanner(connection)

	for scanner.Scan() {
		request := string(scanner.Text())

		segments := strings.Split(request, " ")

		if len(segments) == 0 {
			connection.Write([]byte("[ERR] Request must contain a command\n"))
			continue
		}

		command := segments[0]

		switch command {
		case "SET":
			// SET key val
			err := checkArgs(segments, 3)
			if err != nil {
				connection.Write([]byte(err.Error()))
				continue
			}

			key, val := segments[1], segments[2]

			store.Set(key, val)

			connection.Write([]byte("[OK]\n"))

			continue

		case "DEL":
			// DEL key
			err := checkArgs(segments, 2)
			if err != nil {
				connection.Write([]byte(fmt.Sprintf("[ERR] %s\n", ErrNotFound.Error())))
			}

			key := segments[1]

			fmt.Printf("Deleting %s", key)

			continue
		case "GET":
			// GET key
			err := checkArgs(segments, 2)
			if err != nil {
				connection.Write([]byte(err.Error()))
				continue
			}

			key := segments[1]
			entry, err := store.Get(key)

			if errors.Is(err, ErrNotFound) {
				connection.Write([]byte(fmt.Sprintf("[ERR] %s\n", ErrNotFound.Error())))
				continue
			}

			connection.Write([]byte(fmt.Sprintf("[OK] %s\n", entry)))

			continue

		case "EXISTS":
			//EXISTS key
			err := checkArgs(segments, 2)
			if err != nil {
				connection.Write([]byte(fmt.Sprintf("[ERR] %s\n", ErrNotFound.Error())))
			}

			key := segments[1]

			fmt.Printf("Checking %s\n", key)

			continue

		case "KEYS":
			// KEYS
			// no check needed

			continue
		default:
			connection.Write([]byte(fmt.Sprintf("[ERR] Command %s\n", ErrNotFound.Error())))
			continue
		}
	}
}

func checkArgs(segments []string, count int) error {
	if len(segments) < count {
		err := fmt.Sprintf("[ERR] Command expects at least %d arguments", count)
		return errors.New(err)
	}

	return nil
}
