package rerpc

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

// Error wraps Go's built-in error interface and adds support for gRPC error
// codes and error details. The output of the wrapped error's Error() method is
// sent to the client as the gRPC error message; when building public APIs,
// take care not to leak sensitive information.
//
// Error codes and messages are explained in the gRPC documentation linked
// below. Unfortunately, error details were introduced before gRPC adopted a
// formal proposal process, so they're not clearly documented anywhere and
// differ slightly between implementations. Roughly, they're an optional
// mechanism for servers, middleware, and proxies to send strongly-typed errors
// to clients.
//
// Related documents:
//   gRPC HTTP/2 specification: https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
//   gRPC status codes: https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
//
// FIXME: examples!
type Error struct {
	code    Code
	err     error
	details []*any.Any
}

// Wrap annotates any error with a gRPC status code.
func Wrap(c Code, err error) *Error {
	return &Error{
		code:  c,
		err:   err,
	}
}

// Wrapf calls errors.Errorf with the supplied template and arguments, then
// wraps the resulting error.
func Wrapf(c Code, template string, args ...interface{}) *Error {
	return Wrap(c, fmt.Errorf(template, args...))
}

// AsError uses errors.As to unwrap any error and look for a reRPC Error.
func AsError(err error) (*Error, bool) {
	var re *Error
	ok := errors.As(err, &re)
	return re, ok
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.code, e.err)
}

// Unwrap implements errors.Wrapper, which allows errors.Is and errors.As
// access to the underlying error.
func (e *Error) Unwrap() error {
	return e.err
}

// Code returns the error's gRPC code.
func (e *Error) Code() Code {
	return e.code
}

// Detail returns a deep copy of the error's details.
func (e *Error) Details() []*any.Any {
	if len(e.details) == 0 {
		return nil
	}
	ds := make([]*any.Any, len(e.details))
	for i, d := range e.details {
		ds[i] = proto.Clone(d).(*any.Any)
	}
	return ds
}

// AddDetail appends a message to the error's details.
func (e *Error) AddDetail(m proto.Message) error {
	if d, ok := m.(*any.Any); ok {
		e.details = append(e.details, proto.Clone(d).(*any.Any))
		return nil
	}
	detail, err := ptypes.MarshalAny(m)
	if err != nil {
		return fmt.Errorf("can't add message to error details: %w", err)
	}
	e.details = append(e.details, detail)
	return nil
}

// SetDetails overwrites the error's details.
func (e *Error) SetDetails(details ...proto.Message) error {
	e.details = make([]*any.Any, 0, len(details))
	for _, d := range details {
		if err := e.AddDetail(d); err != nil {
			return err
		}
	}
	return nil
}
