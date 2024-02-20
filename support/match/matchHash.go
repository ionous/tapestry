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

var Keywords = struct {
	Comma, QuotedText uint64
}{Hash(`,`), Hash(`"`)}

// transform a string into a slice of hashes.
// . lowercases and trims the string ( using ToLower since matching depends on english span and phrasing anyway )
// . considers commas their own words ( otherwise commas would wind up as part of span )
// . combines any double quoted or backtick'd text into two Words.
//   - a leading double quote,
//   - a word consisting of all the text until the end of the quote.
//   - ( errors on unmatched quote or backtick )
//
// rationale: since trimming and separating by spaces would require string allocation;
// might as well generate some hashes and slightly speed up comparisons.
// .
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
	var tickMatching int
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

		// back ticks
		case r == '`' && quoteMatching == 0:
			if tickMatching <= 0 {
				tickMatching = i + 1 // one indexed
				flushWord(i, i+1, Keywords.QuotedText)
			} else {
				flushWord(tickMatching, i, sumReset(w))
				tickMatching = 0
				wordStart = -1
			}

		case tickMatching > 0:
			rbs.writeRune(r, w)

		// double quotes
		case r == '"' && tickMatching == 0:
			if quoteMatching <= 0 {
				quoteMatching = i + 1 // one indexed
				flushWord(i, i+1, Keywords.QuotedText)
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
			flushWord(i, i+1, Keywords.Comma)
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
		err = fmt.Errorf("unmatched quote at %d", quoteMatching-1)
	} else if tickMatching > 0 {
		err = fmt.Errorf("unmatched tick at %d", tickMatching-1)
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
