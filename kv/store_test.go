package kv

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

const benchOpsPerIter = 128

func newTestStore(logPath string) *Store {
	return &Store{
		Log:  *NewLog(logPath),
		Data: make(map[string]Event),
	}
}

func newBenchmarkStore() *Store {
	return newTestStore(os.DevNull)
}

func reportAverageLatencyAndThroughput(b *testing.B, totalOps int, latencyUnit string, throughputUnit string) {
	b.Helper()

	totalDuration := b.Elapsed()
	avgLatency := totalDuration / time.Duration(totalOps)

	b.ReportMetric(float64(avgLatency.Nanoseconds()), latencyUnit)
	b.ReportMetric(float64(totalOps)/totalDuration.Seconds(), throughputUnit)
}

func TestNewStoreInitializesData(t *testing.T) {
	t.Parallel()

	store := NewStore()

	if store == nil {
		t.Fatal("NewStore() returned nil")
	}

	if store.Log.Path != "log.txt" {
		t.Fatalf("unexpected log path: got %q want %q", store.Log.Path, "log.txt")
	}

	if store.Data == nil {
		t.Fatal("NewStore() did not initialize Data")
	}
}

func TestStoreSetAndGet(t *testing.T) {
	t.Parallel()

	store := newTestStore(os.DevNull)

	if err := store.Set("name", "klefki", 60); err != nil {
		t.Fatalf("Set(): %v", err)
	}

	event, err := store.Get("name")
	if err != nil {
		t.Fatalf("Get(): %v", err)
	}

	if event.Key != "name" {
		t.Fatalf("unexpected key: got %q want %q", event.Key, "name")
	}

	if event.Value != "klefki" {
		t.Fatalf("unexpected value: got %q want %q", event.Value, "klefki")
	}

	if event.Type != "SET" {
		t.Fatalf("unexpected type: got %q want %q", event.Type, "SET")
	}

	if event.TTL != 60 {
		t.Fatalf("unexpected ttl: got %d want %d", event.TTL, 60)
	}

	if event.Time.IsZero() {
		t.Fatal("expected Set() to populate event time")
	}
}

func TestStoreGetMissingKey(t *testing.T) {
	t.Parallel()

	store := newTestStore(os.DevNull)

	_, err := store.Get("missing")
	if err != ErrNotFound {
		t.Fatalf("unexpected error: got %v want %v", err, ErrNotFound)
	}
}

func TestStoreDeleteRemovesKey(t *testing.T) {
	t.Parallel()

	store := newTestStore(os.DevNull)

	store.Data["name"] = Event{
		Type:  "SET",
		Time:  time.Now(),
		Key:   "name",
		Value: "klefki",
		TTL:   60,
	}

	if err := store.Delete("name"); err != nil {
		t.Fatalf("Delete(): %v", err)
	}

	if _, err := store.Get("name"); err != ErrNotFound {
		t.Fatalf("unexpected error after delete: got %v want %v", err, ErrNotFound)
	}
}

func TestStoreListKeysAndValues(t *testing.T) {
	t.Parallel()

	store := newTestStore(os.DevNull)
	now := time.Now()

	store.Data["alpha"] = Event{Type: "SET", Time: now, Key: "alpha", Value: "one", TTL: 60}
	store.Data["beta"] = Event{Type: "SET", Time: now, Key: "beta", Value: "two", TTL: 60}

	keys := store.ListKeys()
	values := store.ListValues()

	slices.Sort(keys)
	slices.Sort(values)

	if want := []string{"alpha", "beta"}; !slices.Equal(keys, want) {
		t.Fatalf("unexpected keys: got %v want %v", keys, want)
	}

	if want := []string{"one", "two"}; !slices.Equal(values, want) {
		t.Fatalf("unexpected values: got %v want %v", values, want)
	}
}

func TestStoreLoadRebuildsStateFromLog(t *testing.T) {
	t.Parallel()

	logPath := filepath.Join(t.TempDir(), "store.log")
	log := NewLog(logPath)
	now := time.Now()

	events := []Event{
		{Type: "SET", Time: now, Key: "alpha", Value: "one", TTL: 60},
		{Type: "SET", Time: now, Key: "beta", Value: "two", TTL: 60},
		{Type: "DEL", Key: "alpha"},
	}

	for _, event := range events {
		if err := log.Append(event); err != nil {
			t.Fatalf("Append(%q): %v", event.Key, err)
		}
	}

	store := newTestStore(logPath)
	if err := store.Load(); err != nil {
		t.Fatalf("Load(): %v", err)
	}

	if _, err := store.Get("alpha"); err != ErrNotFound {
		t.Fatalf("unexpected error for deleted key: got %v want %v", err, ErrNotFound)
	}

	event, err := store.Get("beta")
	if err != nil {
		t.Fatalf("Get(beta): %v", err)
	}

	if event.Value != "two" {
		t.Fatalf("unexpected loaded value: got %q want %q", event.Value, "two")
	}
}

func BenchmarkStoreAverageWriteSpeed(b *testing.B) {
	store := newBenchmarkStore()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		base := i * benchOpsPerIter

		for j := 0; j < benchOpsPerIter; j++ {
			key := fmt.Sprintf("key-%d", base+j)
			value := fmt.Sprintf("value-%d", base+j)

			if err := store.Set(key, value, 60); err != nil {
				b.Fatalf("set %q: %v", key, err)
			}
		}
	}

	reportAverageLatencyAndThroughput(b, b.N*benchOpsPerIter, "ns/write", "writes/sec")
}

func BenchmarkStoreAverageReadSpeed(b *testing.B) {
	store := newBenchmarkStore()
	now := time.Now()

	keys := make([]string, benchOpsPerIter)

	for i := 0; i < benchOpsPerIter; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)

		store.Data[key] = Event{
			Type:  "SET",
			Time:  now,
			Key:   key,
			Value: value,
			TTL:   60,
		}
		keys[i] = key
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			if _, err := store.Get(key); err != nil {
				b.Fatalf("get %q: %v", key, err)
			}
		}
	}

	reportAverageLatencyAndThroughput(b, b.N*benchOpsPerIter, "ns/read", "reads/sec")
}
