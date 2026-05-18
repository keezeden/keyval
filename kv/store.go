package kv

import (
	"sync"
	"time"
)

type Event struct {
	Type  string
	Time  time.Time
	Key   string
	Value string
	TTL   int
}

type Store struct {
	mu   sync.RWMutex
	Log  Log
	Data map[string]Event
}

func NewStore() *Store {
	return &Store{
		Log:  *NewLog("log.txt"),
		Data: make(map[string]Event),
	}
}

func (store *Store) Load() error {
	text, err := store.Log.ReadAll()
	if err != nil {
		return err
	}

	records, err := DecodeEvents(text)
	if err != nil {
		return err
	}

	store.Data = records
	return nil
}

func (store *Store) Set(key string, value string, ttl int) error {
	event := Event{
		Type:  "SET",
		Key:   key,
		Value: value,
		Time:  time.Now(),
		TTL:   ttl,
	}

	errChannel := make(chan error, 1)

	go func() {
		err := store.Log.Append(event)
		if err != nil {
			errChannel <- err
		}
		close(errChannel)
	}()

	go func() {
		err := <-errChannel
		if err != nil {
			store.mu.Lock()
			defer store.mu.Unlock()
			delete(store.Data, key)
		}
	}()

	store.mu.Lock()
	defer store.mu.Unlock()
	store.Data[key] = event
	return nil
}

func (store *Store) Delete(key string) error {
	event := Event{
		Type: "DEL",
		Key:  key,
	}

	go store.Log.Append(event)

	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.Data, key)
	return nil
}

func (store *Store) Expire(key string) error {
	event := Event{
		Type: "EXP",
		Key:  key,
	}

	go store.Log.Append(event)

	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.Data, key)
	return nil
}

func (store *Store) Get(key string) (Event, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	event, found := store.Data[key]

	if !found {
		return Event{}, ErrNotFound
	}

	if store.isExpired(event) {
		store.Expire(event.Key)
		return Event{}, ErrNotFound
	}

	return event, nil
}

func (store *Store) ListKeys() []string {
	var values []string

	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, event := range store.Data {
		values = append(values, event.Key)
	}

	return values
}

func (store *Store) ListValues() []string {
	var values []string

	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, event := range store.Data {
		values = append(values, event.Value)
	}

	return values
}

func (store *Store) isExpired(event Event) bool {
	now := time.Now()

	difference := now.Sub(event.Time)
	seconds := difference.Seconds()

	if seconds < float64(event.TTL) {
		return false
	}

	return true
}
