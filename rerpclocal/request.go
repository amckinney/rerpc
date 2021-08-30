package rerpclocal

import (
	"io"

	"github.com/rerpc/rerpc"
)

// Request represents a request sent between a local client and server.
type Request struct {
	Type        rerpc.StreamType
	Path        string
	ContentType string

	Body io.ReadCloser
}
