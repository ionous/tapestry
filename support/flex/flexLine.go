package flex

import (
	"io"
	"strings"
)

// return a string that doesn't include the newline
// err can be nil or eof
// due to the way section reader works; a fake eof is always an empty line.
func readLine(runes io.RuneReader) (ret string, err error) {
	var str strings.Builder
	for {
		if n, _, e := runes.ReadRune(); e != nil || n == newline {
			ret = str.String()
			err = e
			break
		} else {
			str.WriteRune(n)
		}
	}
	return
}
