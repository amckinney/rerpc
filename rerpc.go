package rerpc

import (
	"fmt"
	"math"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// Version is the semantic version of the reRPC module.
const Version = "0.0.1"

// EightKiB is 8KiB in bytes, gRPC's recommended maximum header size.
const EightKiB = 1024 * 8

// UserAgent describes reRPC to servers, following the convention outlined in
// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#user-agents.
var UserAgent = fmt.Sprintf("grpc-go-rerpc/%s (%s)", Version, runtime.Version())

// ReRPC's supported HTTP content types. The gRPC variants follow gRPC's
// HTTP/2 protocol, while the JSON variant uses a closely-related protocol
// outlined in reRPC's PROTOCOL.md.
const (
	TypeDefaultGRPC = "application/grpc"
	TypeProtoGRPC   = "application/grpc+proto"
	TypeJSON        = "application/json"
)

// ReRPC's supported compression methods.
const (
	// FIXME: rename to compression
	EncodingIdentity = "identity"
	EncodingGzip     = "gzip"
	EncodingSnappy   = "snappy" // FIXME, implement
)

var acceptEncodingValue = strings.Join([]string{EncodingGzip, EncodingIdentity}, ", ")

// maxHours is the largest number of hours that can be expressed in a
// time.Duration without overflowing.
var maxHours = math.MaxInt64 / int64(time.Hour)

// Doer is the transport-level interface reRPC expects HTTP clients to
// implement. The standard library's http.Client implements Doer.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}
