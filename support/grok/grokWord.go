package grok

import (
	"github.com/ionous/errutil"
)

type wordError struct {
	word   Word
	reason string
}

func makeWordError(w Word, reason string) error {
	return &wordError{w, reason}
}

func (w *wordError) Error() string {
	// i suppose if you wanted to be evil, you would unsafe pointer this string
	// back it up by start to get the actual position
	return errutil.Sprint(w.reason, ">", w.word.slice)
}

type Word struct {
	hash  uint64
	slice string // go doesn't allocate a new string for a slice, it stores offset and length
	start int    // but doesn't seem to be a way of finding the offset from the string itself
}

func (w *Word) isValid() bool {
	return len(w.slice) > 0
}

func (w *Word) equals(other uint64) bool {
	return w.hash == other
}

func (w *Word) String() string {
	return w.slice
}

type spans [][]Word

func panicSpans(strs []string) (out spans) {
	out = make(spans, len(strs))
	for i, str := range strs {
		out[i] = panicHash(str)
	}
	return
}

// find the index and length of a prefix matching the passed words
func (ws spans) findPrefix(words []Word) (retWhich int, retLen int) {
	if wordCount := len(words); wordCount > 0 {
		for prefixIndex, prefix := range ws {
			// every Word in el has to exist in words for it to be a prefix
			// and it has to be longer than any other previous match for it to be the best match
			// ( fix? try a sort search> my first attempt failed miserably )
			if prefixLen := len(prefix); prefixLen <= wordCount && prefixLen > retLen {
				var failed bool
				for i, a := range prefix {
					if a.hash != words[i].hash {
						failed = true
						break
					}
				}
				if !failed {
					retWhich, retLen = prefixIndex, prefixLen
				}
			}
		}
	}
	return
}
