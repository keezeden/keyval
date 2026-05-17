# Klefki

<img src="./.github/klefki.png" width="200" height="200" align="center">
Klefki is a small append-only key-value store that communicates over TCP/IP.

## Why

This project was built mainly to:

- Get reps in with Golang
- Deepen knowledge of data storage internals

## Overview

Klefki uses a simple text-based command format:

```sh
SET name Klefki
> [OK]

GET name
> [OK] Klefki

KEYS
> [OK] name

DEL name
> [OK]

GET name
> [ERR] Key not found

KEYS
> [OK]
```

## Commands

```text
SET <key> <val> <ttl?=300>
```

Stores a value under the given key. An optional TTL can be provided; if omitted, it defaults to `300` seconds.

```text
GET <key>
```

Returns the value stored for the given key.

```text
DEL <key>
```

Deletes the key and its value from the store.

```text
KEYS
```

Lists all keys currently stored.
