# Usage

This is the Go standard library implementation of a linked list 
(https://golang.org/src/container/list/list.go), with the following modifications:
* it uses Go generics
* it allows passing in a `sync.Pool` (via the `NewWithPool` constructor) to reduce allocations of `Element` structs
