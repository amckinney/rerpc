package rerpclocal

import (
	"io"
)

// Response represents a response sent between a local client and server.
type Response struct {
	Body io.WriteCloser
}
