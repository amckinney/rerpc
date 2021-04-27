package internal

// Fetch the latest status definition and generate stubs.
//go:generate curl -o status.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/rpc/status.proto
//go:generate protoc -I . status.proto --go_out=.
