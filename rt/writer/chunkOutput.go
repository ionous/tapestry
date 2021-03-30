package writer

import (
	"github.com/ionous/errutil"
)

// ChunkOutput - adapts go standard output for various methods:
// when go writes bytes, runes, or strings; this turns them all into chunks.
// this is a convenience so that filters, etc. only have to implement one write method
// ( WriteChunk ) not the full set of go's five methods.
// It should return the number of bytes of the *chunk* that were consumed,
// and any error encountered along the way.
// Noting that the returned value might not match the number of bytes written
// if the output is padded or reduced in someway.
type ChunkOutput func(Chunk) (int, error)

// Write redirects the call to WriteChunk
func (n ChunkOutput) Write(p []byte) (int, error) {
	return n(Chunk{p})
}

// WriteByte redirects the call to WriteChunk
func (n ChunkOutput) WriteByte(b byte) error {
	_, e := n(Chunk{b})
	return e
}

// WriteRune redirects the call to WriteChunk
func (n ChunkOutput) WriteRune(q rune) (int, error) {
	return n(Chunk{q})
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
