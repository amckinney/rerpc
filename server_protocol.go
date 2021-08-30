package rerpc

import (
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
)

// HTTPServerProtocol is used by a Handler to create rerpc.Streams from incoming
// requests and format responses.
type HTTPServerProtocol interface {
	// ContentType returns the Content-Type for which this protocol is registered.
	ContentType() string

	// NewStream takes an incoming request and response writer and returns
	// a Stream that should be used for the RPC.
	NewStream(StreamType, http.ResponseWriter, *http.Request) (StreamFunc, error)
}

func defaultHTTPServerProtocols() map[string]HTTPServerProtocol {
	var (
		jsonProtobufCodec = jsonProtobufCodec{
			marshaler:   protojson.MarshalOptions{UseProtoNames: true},
			unmarshaler: protojson.UnmarshalOptions{DiscardUnknown: true},
		}
		protobufCodec = protobufCodec{}
	)
	return map[string]HTTPServerProtocol{
		"*": grpcServerProtocol{
			contentType: "application/grpc",
			codec:       protobufCodec,
		},

		"application/grpc": grpcServerProtocol{
			contentType: "application/grpc",
			codec:       protobufCodec,
		},

		"application/grpc+proto": grpcServerProtocol{
			contentType: "application/grpc+proto",
			codec:       protobufCodec,
		},

		// TODO(alex): With this abstraction, we could support
		// a separate twirpServerProtocol type that owns the
		// application/protobuf and application/json Content-Types.
		//
		// Similarly, we could have separate implementations for gRPC-Web
		// and other HTTP-oriented protocols like Websockets.
		"application/protobuf": grpcServerProtocol{
			contentType: "application/protobuf",
			codec:       protobufCodec,
		},

		"application/json": grpcServerProtocol{
			contentType: "application/json",
			codec:       jsonProtobufCodec,
		},
	}
}
