package print

import (
	"bytes"

	"git.sr.ht/~ionous/tapestry/rt/writer"
)

// Lines implements io.Writer, buffering every Write as a new line.
// use MakeChunks to construct a valid line writer.
type Lines struct {
	writer.ChunkOutput
	lines []string
}

func NewLines() *Lines {
	ls := new(Lines)
	ls.ChunkOutput = writer.ChunkOutput(ls.WriteChunk)
	return ls
}

// Lines returns all current lines.
// There is no flush. A new line writer can be constructed instead.
func (ls *Lines) Lines() []string {
	return ls.lines
}

// Write implements writer.Output, spacing writes with separators.
func (ls *Lines) WriteChunk(c writer.Chunk) (int, error) {
	var buf bytes.Buffer
	n, e := c.WriteTo(&buf)
	ls.lines = append(ls.lines, buf.String())
	return n, e
}
