package printer

import (
	"io"
)

// Sep implements io.Writer, treating every Write as a new word.
type Sep struct {
	io.Writer
	last  string // separators
	queue string // last string sent to Write()
	cnt   int    // number of non-zero writes to the underlying writer.
}

// AndSeparator creates a phrase: a, b, c, and d.
// Note: spacing between words is left to printer.Spacing.
func AndSeparator(w io.Writer) io.WriteCloser {
	return &Sep{Writer: w, last: "and"}
}

// OrSeparator creates a phrase: a, b, c, or d.
// Note: spacing between words is left to printer.Spacing.
func OrSeparator(w io.Writer) io.WriteCloser {
	return &Sep{Writer: w, last: "or"}
}

// Write implements io.Writer, spacing writes with separators.
func (l *Sep) Write(p []byte) (ret int, err error) {
	const mid = ","
	if len(p) > 0 {
		s := string(p)
		err = l.flush(mid)
		l.queue = s
		ret = len(s)
	}
	return
}

// Close writes all pending lines with appropriate separators.
func (l *Sep) Close() error {
	var fini string
	if l.cnt > 1 {
		fini = ", " + l.last
	} else {
		fini = l.last
	}
	return l.flush(fini)
}

// Flush empties the queue
func (l *Sep) flush(sep string) (err error) {
	if len(l.queue) > 0 {
		if l.cnt != 0 {
			_, e := io.WriteString(l.Writer, sep)
			err = e
		}
		if err == nil {
			_, e := io.WriteString(l.Writer, l.queue)
			err = e
		}
		l.queue = ""
		l.cnt++
	}
	return
}