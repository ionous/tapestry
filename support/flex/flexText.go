package flex

import (
	"io"
)

// consumes all text until eof ( and eats the eof error )
// FIX: have to be able to pivot to tell subsections
func ReadText(runes io.RuneReader) (ret []string, err error) {
	if lines, e := readLines(runes); e != nil && e != io.EOF {
		err = e
	} else {
		ret = lines
	}
	return
}

// consumes all lines until eof; returns eof
func readLines(runes io.RuneReader) (ret []string, err error) {
	for err == nil {
		var line string
		if line, err = readLine(runes); err == nil || err == io.EOF {
			ret = append(ret, line)
		}
	}
	return
}
