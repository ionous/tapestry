package writer

import "os"

func NewStdout() *FileWriter {
	return &FileWriter{os.Stdout}
}

type FileWriter struct {
	*os.File
}

func (fp FileWriter) WriteByte(b byte) error {
	_, e := fp.File.Write([]byte{b})
	return e
}

func (fp FileWriter) WriteRune(q rune) (int, error) {
	return WriteRune(fp.File, q)
}
