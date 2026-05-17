package kv

import (
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

	err := store.Log.Append(event)
	if err != nil {
		return err
	}

	store.Data[key] = event
	return nil
}

func (store *Store) Delete(key string) error {
	event := Event{
		Type: "DEL",
		Key:  key,
	}

	err := store.Log.Append(event)
	if err != nil {
		return err
	}

	delete(store.Data, key)
	return nil
}

func (store *Store) Expire(key string) error {
	event := Event{
		Type: "EXP",
		Key:  key,
	}

	err := store.Log.Append(event)
	if err != nil {
		return err
	}

	delete(store.Data, key)
	return nil
}

func (store *Store) Get(key string) (Event, error) {
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

	for _, event := range store.Data {
		values = append(values, event.Key)
	}

	return values
}

func (store *Store) ListValues() []string {
	var values []string

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
