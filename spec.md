Key-Value Store Spec

What it is
A persistent, networked key-value store that clients can talk to over a raw TCP connection using a simple text protocol.

Functional Requirements
Core operations

SET key value — store a value under a key, overwrite if exists
GET key — retrieve a value by key
DEL key — delete a key
EXISTS key — return whether a key exists
KEYS — return all current keys

Expiry

SET key value EX seconds — key automatically disappears after N seconds
TTL key — return how many seconds remain before a key expires, -1 if no expiry, -2 if key doesn't exist

Persistence

All writes are recorded to an append-only log on disk before being acknowledged
On startup, the store rebuilds its state by replaying the log from beginning to end
The log should be human-readable

Non-Functional Requirements

Multiple clients can connect and issue commands simultaneously
A slow or stuck client must not block other clients
Data written to disk must survive a process crash or restart
Response time for GET/SET should be under 1ms on localhost

Protocol

Plain text over TCP, one command per line
Server responds with a single line per command
Success and error responses are distinct and consistent
Unknown commands return an error, never silence

Failure Behaviour

GET on a missing or expired key returns a clear not-found response (not an error)
DEL on a non-existent key is a no-op, not an error
A malformed command returns an error without crashing the server
If the log file is corrupted on startup, the server should load as much state as it can and report what it skipped

Out of Scope

Authentication
Replication or clustering
Binary values
More than one database / namespace

Done When

A client can SET a key, kill the server process, restart it, and GET the same key back
10 simultaneous clients can read and write without corrupting each other's data
An expired key is never returned, even if it still exists in memory
