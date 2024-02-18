package match

import (
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"io"
	"unicode"
	"unicode/utf8"
)

// transform a string into a customized set of hash span.
// . lowercases and trims the string ( using ToLower since grok depends on english span and phrasing anyway )
// . considers commas their own words ( otherwise commas would wind up as part of span )
// . combines double quoted text into a single Word ( errors on unmatched quote )
// fix? quote escaping?
// rationale: since trimming and separating by spaces would require string allocation (probably multiple)
// we might as well generate some hashes instead.
func MakeSpan(sentence string) (ret Span, err error) {
	if out, rem, e := makeSpan(sentence); e != nil {
		err = e
	} else if len(rem) > 0 {
		err = errors.New("expected at most a single sentence")
	} else {
		ret = out
	}
	return
}

// generates a set of spans from a string containing one or more sentences.
func MakeSpans(sentences string) (ret []Span, err error) {
	next := sentences
	for len(next) > 0 {
		if out, rem, e := makeSpan(next); e != nil {
			err = e
		} else {
			ret = append(ret, out)
			next = rem
		}
	}
	return
}

func makeSpan(s string) (out Span, rem string, err error) {
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
				rem = s[i:]
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
			flushWord(i, i+1, comma)
			wordStart = -1

		case r != '-' && unicode.IsPunct(r):
			// alt: eat like normalize does?
			err = fmt.Errorf("unexpected punctuation at %d", i)
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
		err = fmt.Errorf("unmatched quotes at %d", quoteMatching-1)
	} else if err == nil {
		flushWord(wordStart, len(s), sumReset(w))
	}
	return
}

var comma = Hash(",")

func sumReset(w hash.Hash64) uint64 {
	out := w.Sum64()
	w.Reset()
	return out
}

func Hash(s string) uint64 {
	w, rbs := fnv.New64a(), makeRuneWriter()
	for _, r := range s {
		r := unicode.ToLower(r)
		rbs.writeRune(r, w)
	}
	return w.Sum64()
}

func Hashes(strs []string) []uint64 {
	out := make([]uint64, len(strs))
	for i, s := range strs {
		out[i] = Hash(s)
	}
	return out
}

type runeWriter []byte

func makeRuneWriter() runeWriter {
	return make([]byte, utf8.UTFMax)
}

func (rbs runeWriter) writeRune(r rune, w io.Writer) {
	c := utf8.EncodeRune(rbs, r)
	_, _ = w.Write(rbs[:c])
}
