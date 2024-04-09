package flex

import (
	"errors"
	"io"
	"strings"
)

// consumes all comments until eof ( and eats the eof error )
// lines are stripped of comment markers and trailing newlines
func ReadComments(runes io.RuneReader) (ret []string, err error) {
	if lines, e := readComments(runes); e != nil && e != io.EOF {
		err = e
	} else {
		ret = lines
	}
	return
}

// consumes all comments until eof; returns eof
func readComments(runes io.RuneReader) (ret []string, err error) {
	for err == nil {
		var line string
		if line, err = readLine(runes); err == nil || err == io.EOF {
			// ignore fully blank lines
			if len(strings.TrimSpace(line)) != 0 {
				// check for comment lines:
				if !strings.HasPrefix(line, "# ") {
					err = errors.New("header can only contain comments")
					break
				} else {
					ret = append(ret, line[2:])
				}
			}
		}
	}
	return
}
