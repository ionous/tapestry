package grok

import (
	"hash"
	"hash/fnv"
	"io"
	"unicode"
	"unicode/utf8"

	"github.com/ionous/errutil"
)

func panicHash(s string) []Word {
	out, e := hashWords(s)
	if e != nil {
		panic(e)
	}
	return out
}

// transform a string into a customized set of hash span.
// . lowercases and trims the string ( using ToLower since grok depends on english span and phrasing anyway )
// . considers commas their own words ( otherwise commas would wind up as part of span )
// . combines double quoted text into a single Word ( errors on unmatched quote )
// fix? quote escaping?
// rationale: since trimming and separating by spaces would require string allocation (probably multiple)
// we might as well generate some hashes instead.
func hashWords(s string) (out []Word, err error) {
	flushWord := func(start, end int, hash uint64) {
		if start >= 0 {
			if end > start {
				out = append(out, Word{
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
	_, _ = w.Write(rbs[:c])
}
