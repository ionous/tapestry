package print

import (
	"bytes"

	"git.sr.ht/~ionous/iffy/rt/writer"
)

// Sep implements writer.Output, treating every Write as a new word.
type Sep struct {
	target    writer.Output
	mid, last string       // separators
	pending   bytes.Buffer // last string sent to Write()
	cnt       int          // number of non-zero writes to the underlying writer.
}

// AndSeparator creates a phrase: a, b, c, and d.
// Note: spacing between words is left to print.Spacing.
func AndSeparator(w writer.Output) writer.OutputCloser {
	sep := &Sep{target: w, mid: ",", last: "and"}
	return sep.ChunkOutput()
}

// OrSeparator creates a phrase: a, b, c, or d.
// Note: spacing between words is left to print.Spacing.
func OrSeparator(w writer.Output) writer.OutputCloser {
	sep := &Sep{target: w, mid: ",", last: "or"}
	return sep.ChunkOutput()
}

func (l *Sep) ChunkOutput() writer.ChunkOutput {
	return l.WriteChunk
}

// Write implements writer.Output, spacing writes with separators.
func (l *Sep) WriteChunk(c writer.Chunk) (ret int, err error) {
	if c.IsClosed() {
		if l.cnt > 1 {
			l.target.WriteRune(',')
		}
		err = l.flush(l.last)
	} else {
		if !c.IsEmpty() {
			if e := l.flush(l.mid); e != nil {
				err = e
			} else {
				ret, err = c.WriteTo(&l.pending)
			}
		}
	}
	return
}

// Flush writes pending text, prefixed if needed with a separator
func (l *Sep) flush(sep string) (err error) {
	// pending text pending, write it.
	if l.pending.Len() > 0 {
		// separate text already written
		if l.cnt != 0 {
			_, e := l.target.WriteString(sep)
			err = e
		}
		// write the pending text
		if err == nil {
			_, e := l.target.Write(l.pending.Bytes())
			err = e
		}
		l.pending.Reset()
		l.cnt++
	}
	return
}
