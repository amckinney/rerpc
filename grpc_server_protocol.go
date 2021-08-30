package rerpc

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
)

// GPRCServerProtocolOption is an option specific to the GRPC protocol.
type GPRCServerProtocolOption func(*grpcServerProtocolOptions)

// GRPCCodec configures the gRPC HTTPServerProtocol with the
// given Compressors.
func GRPCCodec(codec Codec) GPRCServerProtocolOption {
	return func(opts *grpcServerProtocolOptions) {
		opts.codec = codec
	}
}

// GRPCCompressors configures the gRPC HTTPServerProtocol with the
// given Compressors.
func GRPCCompressors(compressors *Compressors) GPRCServerProtocolOption {
	return func(opts *grpcServerProtocolOptions) {
		opts.compressors = compressors
	}
}

// GRPCDisableTwirp disables Twirp support.
func GRPCDisableTwirp() GPRCServerProtocolOption {
	return func(opts *grpcServerProtocolOptions) {
		opts.disableTwirp = true
	}
}

// GRPCMaxRequestSize limits GRPC requests to the given size.
// If 0 is provided, the default max is used.
func GRPCMaxRequestSize(size uint64) GPRCServerProtocolOption {
	return func(opts *grpcServerProtocolOptions) {
		opts.maxRequestSize = int64(size)
	}
}

// NewGPRCServerProtocols is a convenience function for configuring the default gRPC HTTPServerProtocols
// (i.e. Content-Type == "application/grpc[-proto]" with the Protobuf Codec).
func NewGPRCServerProtocols(opts ...GPRCServerProtocolOption) []HTTPServerProtocol {
	return []HTTPServerProtocol{
		NewGPRCServerProtocol("application/grpc", opts...),
		NewGPRCServerProtocol("application/grpc+proto", opts...),
		NewGPRCServerProtocol("application/protobuf", opts...),

		// TODO(alex): This is a huge hack for making sure that the jsonProtobufCodec is used for the application/json Content-Type by default.
		NewGPRCServerProtocol("application/json", append([]GPRCServerProtocolOption{GRPCCodec(jsonProtobufCodec{})}, opts...)...),
	}
}

// NewGPRCServerProtocol returns the GRPC HTTPServerProtocol.
func NewGPRCServerProtocol(contentType string, opts ...GPRCServerProtocolOption) HTTPServerProtocol {
	options := &grpcServerProtocolOptions{}
	for _, opt := range opts {
		opt(options)
	}
	var maxRequestSize int64
	if options.maxRequestSize != 0 {
		maxRequestSize = options.maxRequestSize
	}
	var codec Codec = &protobufCodec{}
	if options.codec != nil {
		codec = options.codec
	}
	return &grpcServerProtocol{
		contentType:    contentType,
		codec:          codec,
		compressors:    options.compressors,
		maxRequestSize: maxRequestSize,
		disableTwirp:   options.disableTwirp,
	}
}

type grpcServerProtocol struct {
	contentType string
	codec       Codec

	compressors    *Compressors
	maxRequestSize int64
	disableTwirp   bool
}

// ContentType returns the Content-Type associated with this HTTPServerProtocol.
func (p grpcServerProtocol) ContentType() string {
	return p.contentType
}

// NewStream creates a new gRPC server stream.
//
// Note that this function acts very similar to Handler.ServeHTTP in its current state.
func (p grpcServerProtocol) NewStream(stype StreamType, w http.ResponseWriter, r *http.Request) (StreamFunc, error) {
	isBidi := (stype & StreamTypeBidirectional) == StreamTypeBidirectional
	if isBidi && r.ProtoMajor < 2 {
		w.WriteHeader(http.StatusHTTPVersionNotSupported)
		io.WriteString(w, "bidirectional streaming requires HTTP/2")
		return nil, errors.New("TODO" /* TODO(alex): Return a typed error here so that the Handler can exit early. */)
	}
	if r.Method != http.MethodPost {
		// grpc-go returns a 500 here, but interoperability with non-gRPC HTTP
		// clients is better if we return a 405.
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil, errors.New("TODO" /* TODO(alex): Return a typed error here so that the Handler can exit early. */)
	}
	ctype := p.contentType
	if (ctype == TypeJSON || ctype == TypeProtoTwirp) && (p.disableTwirp || stype != StreamTypeUnary) {
		w.Header().Set("Accept-Post", acceptPostValueWithoutTwirp)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return nil, errors.New("TODO" /* TODO(alex): Return a typed error here so that the Handler can exit early. */)
	}
	if ctype != TypeDefaultGRPC && ctype != TypeProtoGRPC && ctype != TypeProtoTwirp && ctype != TypeJSON {
		// grpc-go returns 500, but the spec recommends 415.
		// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests
		w.Header().Set("Accept-Post", acceptPostValueDefault)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return nil, errors.New("TODO" /* TODO(alex): Return a typed error here so that the Handler can exit early. */)
	}

	// We need to parse metadata before entering the interceptor stack, but we'd
	// like to report errors to the client in a format they understand (if
	// possible). We'll collect any such errors here and use them to
	// short-circuit early later on.
	//
	// NB, future refactorings will need to take care to avoid typed nils here.
	var failed *Error

	timeout, err := parseTimeout(r.Header.Get("Grpc-Timeout"))
	if err != nil && err != errNoTimeout {
		// Errors here indicate that the client sent an invalid timeout header, so
		// the error text is safe to send back.
		failed = Wrap(CodeInvalidArgument, err).(*Error)
	} else if err == nil {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		r = r.WithContext(ctx)
	}

	var (
		requestCompression  = CompressionIdentity
		responseCompression = CompressionIdentity
	)
	if ctype == TypeJSON || ctype == TypeProtoTwirp {
		if r.Header.Get("Content-Encoding") == "gzip" {
			requestCompression = CompressionGzip
		}
		// TODO: Actually parse Accept-Encoding instead of this hackery.
		//
		// TODO(alex): This isn't right. Use the compressor implementation
		// rather than just relying on the string identifier. We actually
		// need to implement the Compressors type though.
		//
		if p.compressors == nil && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			responseCompression = CompressionGzip
		}
	} else {
		if me := r.Header.Get("Grpc-Encoding"); me != "" {
			switch me {
			case CompressionIdentity:
				requestCompression = CompressionIdentity
			case CompressionGzip:
				requestCompression = CompressionGzip
			default:
				// Per https://github.com/grpc/grpc/blob/master/doc/compression.md, we
				// should return CodeUnimplemented and specify acceptable compression(s)
				// (in addition to setting the Grpc-Accept-Encoding header).
				if failed == nil {
					failed = Errorf(
						CodeUnimplemented,
						"unknown compression %q: accepted grpc-encoding values are %v",
						me, acceptEncodingValue,
					).(*Error)
				}
			}
		}
		// Follow https://github.com/grpc/grpc/blob/master/doc/compression.md.
		// (The grpc-go implementation doesn't read the "grpc-accept-encoding" header
		// and doesn't support compression method asymmetry.)
		responseCompression = requestCompression
		// TODO(alex): Similar to above, this isn't right. Now that compressors are pluggable we
		// should be able to resolve the compressor based on these headers.
		if p.compressors == nil {
			responseCompression = CompressionIdentity
		} else if mae := r.Header.Get("Grpc-Accept-Encoding"); mae != "" {
			for _, enc := range strings.FieldsFunc(mae, splitOnCommasAndSpaces) {
				switch enc {
				case CompressionGzip:
					responseCompression = CompressionGzip
					// prefer gzip, so no continue
				case CompressionIdentity:
					responseCompression = CompressionIdentity
					continue
				default:
					continue
				}
				break
			}
		}
	}

	// We should write any remaining headers here, since: (a) the implementation
	// may write to the body, thereby sending the headers, and (b) interceptors
	// should be able to see this data.
	//
	// Since we know that these header keys are already in canonical form, we can
	// skip the normalization in Header.Set. To avoid allocating re-allocating
	// the same slices over and over, we use pre-allocated globals for the header
	// values.
	w.Header()["Content-Type"] = typeToSlice(ctype)
	if ctype != TypeJSON && ctype != TypeProtoTwirp {
		w.Header()["Grpc-Accept-Encoding"] = acceptEncodingValueSlice
		w.Header()["Grpc-Encoding"] = compressionToSlice(responseCompression)
		// Every gRPC response will have these trailers.
		w.Header()["Trailer"] = grpcStatusTrailers
	}

	// Unlike gRPC, Twirp manages compression using the standard HTTP mechanisms.
	// Since they apply to the whole stream, it's easiest to handle it here.
	var requestBody io.Reader = r.Body
	if ctype == TypeJSON || ctype == TypeProtoTwirp {
		if requestCompression == CompressionGzip {
			gr, err := getGzipReader(requestBody)
			if err != nil && failed == nil {
				failed = Errorf(CodeInvalidArgument, "can't read gzipped body: %w", err).(*Error)
			} else if err == nil {
				defer putGzipReader(gr)
				defer gr.Close()
				requestBody = gr
			}
		}
		// Checking Content-Encoding ensures that some other user-supplied
		// middleware isn't already compressing the response.
		if responseCompression == CompressionGzip && w.Header().Get("Content-Encoding") == "" {
			w.Header().Set("Content-Encoding", "gzip")
			gw := getGzipWriter(w)
			defer putGzipWriter(gw)
			w = &gzipResponseWriter{ResponseWriter: w, gw: gw}
		}
	}

	sf := StreamFunc(func(ctx context.Context) Stream {
		return newServerStream(
			ctx,
			w,
			&readCloser{Reader: requestBody, Closer: r.Body},
			ctype,
			p.maxRequestSize,
			responseCompression == CompressionGzip,
		)
	})

	if failed != nil {
		stream := sf(r.Context() /* This context will already have been wrapped by the Handler layer */)
		_ = stream.CloseReceive()
		_ = stream.CloseSend(failed)
		return nil, errors.New("TODO" /* TODO(alex): Return a typed error here so that the Handler can exit early. */)
	}

	return sf, nil
}

type grpcServerProtocolOptions struct {
	codec          Codec
	compressors    *Compressors
	maxRequestSize int64
	disableTwirp   bool // Separate protocol implementation.
}
