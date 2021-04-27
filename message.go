package rerpc

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

// MarshalLPM FIXME unexport
var (
	MarshalLPM   = marshalLPM
	UnmarshalLPM = unmarshalLPM
)

func marshal(w io.Writer, msg proto.Message, ctype, encoding string) error {
	switch ctype {
	case TypeJSON:
		return marshalJSON(w, msg, encoding)
	case TypeDefaultGRPC, TypeProtoGRPC:
		return marshalLPM(w, msg, encoding)
	default:
		return Wrapf(CodeInternal, "unsupported content-type %q", ctype)
	}
}

func unmarshal(r io.Reader, msg proto.Message, ctype, encoding string) error {
	switch ctype {
	case TypeJSON:
		return unmarshalJSON(r, msg, encoding)
	case TypeDefaultGRPC, TypeProtoGRPC:
		return unmarshalLPM(r, msg, encoding)
	default:
		return Wrapf(CodeInternal, "unsupported content-type %q", ctype)
	}
}

func marshalJSON(w io.Writer, msg proto.Message, encoding string) error {
	var out io.Writer = w
	switch encoding {
	case EncodingIdentity:
	case EncodingGzip:
		out = gzip.NewWriter(out)
	default:
		return Wrapf(CodeInternal, "unsupported encoding %q", encoding)
	}
	if err := jsonpbMarshaler.Marshal(out, msg); err != nil {
		return Wrapf(CodeInternal, "couldn't marshal protobuf message with encoding %q: %w", encoding, err)
	}
	return nil
}

func unmarshalJSON(r io.Reader, msg proto.Message, encoding string) error {
	var in io.Reader = r
	switch encoding {
	case EncodingIdentity:
	case EncodingGzip:
		// TODO: pool these (or don't, for simplicitly?)
		gz, err := gzip.NewReader(in)
		if err != nil {
			return Wrapf(CodeInternal, "can't decompress gzipped data: %w", err)
		}
		defer gz.Close()
		in = gz
	default:
		return Wrapf(CodeInternal, "unsupported encoding %q", encoding)
	}

	if err := jsonpbUnmarshaler.Unmarshal(in, msg); err != nil {
		return Wrapf(CodeInvalidArgument, "can't unmarshal data into type %T: %w", msg, err)
	}
	return nil
}

func marshalLPM(w io.Writer, msg proto.Message, encoding string) error {
	raw, err := proto.Marshal(msg)
	if err != nil {
		return Wrapf(CodeInternal, "couldn't marshal protobuf message: %v", err)
	}
	data := &bytes.Buffer{}
	var dataW io.Writer = data
	switch encoding {
	case EncodingIdentity:
	case EncodingGzip:
		dataW = gzip.NewWriter(data)
	default:
		return Wrapf(CodeInternal, "unsupported encoding %q", encoding)
	}
	_, err = dataW.Write(raw) // returns uncompressed size, which isn't useful
	if err != nil {
		return Wrapf(CodeInternal, "couldn't compress with %q: %v", encoding, err)
	}
	if c, ok := dataW.(io.Closer); ok {
		if err := c.Close(); err != nil {
			return Wrapf(CodeInternal, "couldn't compress with %q: %v", encoding, err)
		}
	}

	size := data.Len()
	prefixes := [5]byte{}
	if encoding == EncodingIdentity {
		prefixes[0] = 0
	} else {
		prefixes[0] = 1
	}
	binary.BigEndian.PutUint32(prefixes[1:5], uint32(size))

	if _, err := w.Write(prefixes[:]); err != nil {
		return Wrapf(CodeInternal, "couldn't write prefix of length-prefixed message: %v", err)
	}
	if _, err := io.Copy(w, data); err != nil {
		return Wrapf(CodeInternal, "couldn't write data portion of length-prefixed message: %v", err)
	}
	return nil
}

func unmarshalLPM(r io.Reader, msg proto.Message, encoding string) error {
	// Each length-prefixed message starts with 5 bytes of metadata: a one-byte
	// unsigned integer indicating whether the payload is compressed, and a
	// four-byte unsigned integer indicating the message length.
	prefixes := make([]byte, 5)
	n, err := r.Read(prefixes)
	if err != nil || n < 5 {
		// Even an EOF is unacceptable here, since we always need a message for
		// unary RPC.
		return Wrapf(CodeInvalidArgument, "missing length-prefixed message metadata: %v", err)
	}

	var compressed bool
	switch prefixes[0] {
	case 0:
		compressed = false
		if encoding != EncodingIdentity {
			return Wrapf(
				CodeInvalidArgument,
				"length-prefixed message is uncompressed but message encoding is %q", encoding,
			)
		}
	case 1:
		compressed = true
		if encoding == EncodingIdentity {
			return Wrapf(
				CodeInvalidArgument,
				"length-prefixed message is compressed but message encoding is %q", EncodingIdentity,
			)
		}
	default:
		return Wrapf(
			CodeInvalidArgument,
			"length-prefixed message has compressed flag %v", prefixes[0],
		)
	}

	size := int(binary.BigEndian.Uint32(prefixes[1:5]))
	if size < 0 {
		return Wrapf(CodeInvalidArgument, "message size %d overflows uint32", size)
	}

	raw := make([]byte, size)
	if size > 0 {
		n, err = r.Read(raw)
		if err != nil && err != io.EOF {
			return Wrapf(CodeInternal, "error reading length-prefixed message data: %w", err)
		}
		if n < size {
			return Wrapf(
				CodeInvalidArgument,
				"promised %d bytes in length-prefixed message, got %d bytes", size, n,
			)
		}
	}

	if compressed && encoding == EncodingGzip {
		// TODO: pool, convert to switch
		gr, err := gzip.NewReader(bytes.NewReader(raw))
		if err != nil {
			return Wrapf(CodeInternal, "can't decompress gzipped data: %w", err)
		}
		defer gr.Close()
		decompressed, err := ioutil.ReadAll(gr)
		if err != nil {
			return Wrapf(CodeInvalidArgument, "can't decompress gzipped data: %w", err)
		}
		raw = decompressed
	}

	if err := proto.Unmarshal(raw, msg); err != nil {
		return Wrapf(CodeInvalidArgument, "can't unmarshal data into type %T: %w", msg, err)
	}

	return nil
}
