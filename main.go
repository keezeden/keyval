package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
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
			log.Println("[ERR] Accepting Connection:", err)
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
			err := checkArgs(segments, 3)
			if err != nil {
				connection.Write([]byte(err.Error()))
				continue
			}

			key, val := segments[1], segments[2]

			if len(segments) > 3 {
				ttl, err := strconv.Atoi(segments[3])
				if err != nil {
					connection.Write([]byte("[ERR] TTL must be an int.\n"))
				}
				store.Set(key, val, ttl)
			} else {
				store.Set(key, val, 300)
			}

			connection.Write([]byte("[OK]\n"))

			continue

		case "DEL":
			err := checkArgs(segments, 2)
			if err != nil {
				connection.Write([]byte(err.Error()))
			}

			key := segments[1]

			store.Delete(key)

			connection.Write([]byte("[OK]\n"))

			continue
		case "GET":
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

		case "KEYS":
			keys := store.ListKeys()
			connection.Write([]byte(fmt.Sprintf("[OK] %s\n", strings.Join(keys, ", "))))

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
