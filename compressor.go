package rerpc

import (
	"io"
)

// Compressor ...
type Compressor interface {
	Name() string // e.g., "gzip" or "snappy"

	// method set should make pooling really easy
	GetReader() io.Reader
	PutReader(io.Reader)

	GetWriter() io.Writer
	PutWriter(io.Writer)
}

// Compressors is a collection of Compressors, uniquely
// identified by a name.
type Compressors struct {
	values map[string]Compressor
}

// Get returns the Compressor associated with the given name, if any.
func (cs *Compressors) Get(name string) (Compressor, bool) {
	compressor, ok := cs.values[name]
	return compressor, ok
}

// nopCompressor returns a no-op compressor.
type nopCompressor struct{}

func (n nopCompressor) Name() string {
	return "nop"
}

func (n nopCompressor) GetReader() io.Reader {
	return nil
}

func (n nopCompressor) PutReader(io.Reader) {
	return
}

func (n nopCompressor) GetWriter() io.Writer {
	return io.Discard
}

func (n nopCompressor) PutWriter(io.Writer) {
	return
}
