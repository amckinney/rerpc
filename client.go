package rerpc

import (
	"bytes"
	"context"

	"github.com/golang/protobuf/proto"
)

type Client struct {
	BaseURL        string // server URL, including any prefixes
	Doer           Doer   // transport-level HTTP client
	Name           string // fully-qualified gRPC-style name, e.g., "rerpc.ping.v0.Ping/Ping"
	Implementation func(context.Context, proto.Message) (proto.Message, error)
}

func (c *Client) Call(ctx context.Context, req, res proto.Message) error {
	body := &bytes.Buffer{}
	return marshalLPM(body, req, EncodingIdentity)
}
