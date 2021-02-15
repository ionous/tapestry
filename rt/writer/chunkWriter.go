package writer

import (
	"github.com/ionous/errutil"
)

// ChunkWriter adapts a single WriteChunk() method into writer.Output friendly interface.
type ChunkWriter interface {
	WriteChunk(Chunk) (int, error)
}

// ChunkOutput - implements go standard output for the specific method
type ChunkOutput func(Chunk) (int, error)

// Write redirects the call to WriteChunk
func (n ChunkOutput) Write(p []byte) (int, error) {
	return n(Chunk{p})
}

// WriteByte redirects the call to WriteChunk
func (n ChunkOutput) WriteByte(c byte) error {
	_, e := n(Chunk{c})
	return e
}

// WriteRune redirects the call to WriteChunk
func (n ChunkOutput) WriteRune(r rune) (int, error) {
	return n(Chunk{r})
}

// WriteString redirects the call to WriteChunk
func (n ChunkOutput) WriteString(s string) (int, error) {
	return n(Chunk{s})
}

// Close redirects the call to WriteChunk with a Closed error
func (n ChunkOutput) Close() error {
	_, e := n(Chunk{Closed})
	return e
}

// Closed is used by ChunkOutput to indicate a request to close the writer.
// ( why not io.EOF?  ChunkWriter is custom so the error can be custom too. )
const Closed = errutil.Error("Closed")
