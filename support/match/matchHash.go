package match

import (
	"hash/fnv"
	"io"
	"unicode"
	"unicode/utf8"
)

func Hash(s string) uint64 {
	w, rbs := fnv.New64a(), makeRuneWriter()
	for _, r := range s {
		r := unicode.ToLower(r)
		rbs.writeRune(r, w)
	}
	return w.Sum64()
}

type runeWriter []byte

func makeRuneWriter() runeWriter {
	return make([]byte, utf8.UTFMax)
}

func (rbs runeWriter) writeRune(r rune, w io.Writer) {
	c := utf8.EncodeRune(rbs, r)
	_, _ = w.Write(rbs[:c])
}
