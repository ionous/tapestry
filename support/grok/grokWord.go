package grok

import (
	"hash"
	"hash/fnv"
	"io"
	"unicode"
	"unicode/utf8"

	"github.com/ionous/errutil"
)

type wordError struct {
	word   word
	reason string
}

func makeWordError(w word, reason string) error {
	return &wordError{w, reason}
}

func (w *wordError) Error() string {
	// i suppose if you wanted to be evil, you would unsafe pointer this string
	// back it up by start to get the actual position
	return errutil.Sprint(w.reason, ">", w.word.slice)
}

type word struct {
	hash  uint64
	slice string // go doesn't allocate a new string for a slice, it stores offset and length
	start int    // but doesn't seem to be a way of finding the offset from the string itself
}

func (w *word) isValid() bool {
	return len(w.slice) > 0
}

func (w *word) equals(other uint64) bool {
	return w.hash == other
}

func (w *word) String() string {
	return w.slice
}

type spans []span
type span []word

func makeSpans(strs []string) (out spans) {
	out = make(spans, len(strs))
	for i, str := range strs {
		out[i] = panicHash(str)
	}
	return
}

// find the index of the span which matches the passed words
func (ws spans) findPrefix(match []word) (ret int) {
	ret = -1 // provisionally
	if mcnt := len(match); mcnt > 0 {
		best := -1
		for prefixIndex, prefix := range ws {
			// every word in el has to exist in match for it to be a prefix
			// and it has to be longer than any other previous match for it to be the best match
			// ( fix? try a sort search> my first attempt failed miserably )
			if pcnt := len(prefix); pcnt <= mcnt && pcnt > best {
				var failed bool
				for i, a := range prefix {
					if a.hash != match[i].hash {
						failed = true
						break
					}
				}
				if !failed {
					ret, best = prefixIndex, pcnt
				}
			}
		}
	}
	return
}

// given this pool of known words;
// split the leading part of a phrase from its tail "[the brace of] quick foxes"
func (ws spans) cut(name []word) (ret []word, rest []word) {
	if at := ws.findPrefix(name); at < 0 {
		rest = name // no det, everything is a name
	} else {
		cnt := len(ws[at])
		ret = name[:cnt]  // the first part is the span of determiners
		rest = name[cnt:] // the second part is anything after
	}
	return
}

func panicHash(s string) []word {
	out, e := hashWords(s)
	if e != nil {
		panic(e)
	}
	return out
}

// transform a string into a customized set of hash span.
// . lowercases and trims the string ( using ToLower since grok depends on english span and phrasing anyway )
// . considers commas their own words ( otherwise commas would wind up as part of span )
// . combines double quoted text into a single word ( errors on unmatched quote )
// fix? quote escaping?
// rationale: since trimming and separating by spaces would require string allocation (probably multiple)
// we might as well generate some hashes instead.
func hashWords(s string) (out []word, err error) {
	flushWord := func(start, end int, hash uint64) {
		if start >= 0 {
			if end > start {
				out = append(out, word{
					start: start,
					hash:  hash,
					slice: s[start:end],
				})
			}
		}
	}
	var wordStart int
	var quoteMatching int // one indexed
	var terminal int
	w, rbs := fnv.New64a(), makeRuneWriter()
Loop:
	for i, r := range s {
		switch {
		case terminal > 0:
			if !unicode.IsSpace(r) {
				err = errutil.New("unexpected full stop at", terminal)
				break Loop
			}

		case r == '"':
			if quoteMatching <= 0 {
				quoteMatching = i + 1 // one indexed
			} else {
				flushWord(quoteMatching, i, sumReset(w))
				quoteMatching = 0
				wordStart = -1
			}

		case quoteMatching > 0:
			rbs.writeRune(r, w)

		case r == '.':
			flushWord(wordStart, i, sumReset(w))
			terminal = i
			wordStart = -1

		case r == ',':
			flushWord(wordStart, i, sumReset(w))
			flushWord(i, i+1, keywords.comma)
			wordStart = -1

		case unicode.IsPunct(r):
			err = errutil.New("unexpected punctuation at", i)
			break Loop

		case unicode.IsSpace(r):
			flushWord(wordStart, i, sumReset(w))
			// tricky: re: leading spaces, this blocks the final flush
			// unless something valid gets written to the hash.
			wordStart = -1

		// some normal character
		default:
			rbs.writeRune(unicode.ToLower(r), w)
			if wordStart < 0 {
				wordStart = i
			}
		}
	}
	if quoteMatching > 0 {
		err = errutil.New("unmatched quotes at:", quoteMatching-1)
	} else if err == nil {
		flushWord(wordStart, len(s), sumReset(w))
	}
	return
}

func sumReset(w hash.Hash64) uint64 {
	out := w.Sum64()
	w.Reset()
	return out
}

func plainHash(s string) uint64 {
	w, rbs := fnv.New64a(), makeRuneWriter()
	for _, r := range s {
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
	w.Write(rbs[:c])
}
