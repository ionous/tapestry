package writer

// Output defines a sink for go friendly text output.
// It's a subset of strings.Builder.
type Output interface {
	Write(p []byte) (int, error)
	WriteByte(c byte) error
	WriteRune(r rune) (int, error)
	WriteString(s string) (int, error)
}

// OutputCloser - callers should attempt to Close() when they are done with a writer.
// ( this conforms to the io.Closer interface )
type OutputCloser interface {
	Output
	Close() error
}
