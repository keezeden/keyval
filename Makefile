.PHONY: test test-kv test-concurrent bench bench-average

test:
	go test ./...

test-kv:
	go test ./kv

test-concurrent:
	go test ./kv -run TestStoreSetConcurrentWrites

bench:
	go test ./kv -run ^$$ -bench .

bench-average:
	go test ./kv -run ^$$ -bench Average
