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
			break
		} else {
			ret = append(ret, out)
			next = rem
		}
	}
	return
}

// fix: replace this with flex plainText
// ( in a way that keeps this independent of flex;
// | a default implementation that doesnt generate commands;
// | ignores or errors on comments; or that accumulates a different kind of output
// )
func makeSpan(s string) (out Span, rem string, err error) {
	flushWord := func(start, end int, hash uint64) {
		if start >= 0 {
			if end > start {
				w := MakeWord(hash, s[start:end])
				out = append(out, w)
			}
		}
	}
	var wordStart int
	var quoteStart int // one indexed
	var tickStart int
	var quoteTerminal bool // watches for fullstops at the end of a quote.
	var terminal int       // flag for end of sentence; uses a number for debugging.
	// the rune writer writes into w to accumulate the hash.
	w, rbs := fnv.New64a(), makeRuneWriter()
Loop:
	for i, r := range s {
		switch {
		case terminal > 0:
			// eat spaces after terminal; anything else is a new sentence.
			if !unicode.IsSpace(r) {
				rem = s[i:]
				break Loop
			}

		// backticked text is written with two "words":
		// 1. a word containing a single doubleQuote ( ie. the leading backtick is replaced )
		// 2. a word containing the entire contents of the quoted text
		case r == '`' && quoteStart == 0:
			// start reading quoted text:
			if tickStart == 0 {
				tickStart = i + 1 // one indexed
				flushWord(i, i+1, Keywords.QuotedText)
			} else {
				// done reading quoted text
				flushWord(tickStart, i, sumReset(w))
				if quoteTerminal {
					terminal, quoteTerminal = i, false
				}
				tickStart = 0
				wordStart = -1
			}

		// double quoted text is written with two "words":
		// 1. a word containing a single doubleQuote
		// 2. a word containing the entire contents of the quoted text
		case r == '"' && tickStart == 0:
			// start reading quoted text:
			if quoteStart == 0 {
				quoteStart = i + 1 // one indexed
				flushWord(i, i+1, Keywords.QuotedText)
			} else {
				// done reading quoted text
				if quoteTerminal {
					terminal, quoteTerminal = i, false
				}
				flushWord(quoteStart, i, sumReset(w))
				quoteStart = 0
				wordStart = -1
			}

		case quoteStart > 0 || tickStart > 0:
			quoteTerminal = unicode.Is(unicode.Sentence_Terminal, r)
			rbs.writeRune(r, w)

		case r == '.':
			flushWord(wordStart, i, sumReset(w))
			terminal = i
			wordStart = -1

		// commas are written as their own word.
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
	if quoteStart > 0 {
		err = fmt.Errorf("unmatched quote at %d", quoteStart-1)
	} else if tickStart > 0 {
		err = fmt.Errorf("unmatched tick at %d", tickStart-1)
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
