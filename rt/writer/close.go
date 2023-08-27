package writer

// Close(es) the writer if it implements io.Closer
func xClose(w interface{}) (err error) {
	type closer interface{ Close() error }
	if c, ok := w.(closer); ok {
		err = c.Close()
	}
	return
}
