package flex

import (
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"github.com/ionous/tell/charm"
)

// consumes all text until eof ( and eats the eof error )
// fix: allow line number offset
func ReadText(runes io.RuneReader) (ret []story.StoryStatement, err error) {
	var pt PlainText
	run := NewTokenizer(&pt)
	if e := charm.Read(runes, run); e != nil {
		err = e
	} else {
		ret, err = pt.Finalize()
	}
	return
}

// translate a plain text section to a series of
// comments and jess declarations
type PlainText struct {
	// declare statements, or comments
	out []story.StoryStatement
	// accumulator for declare, and comments
	phrases []match.Span
	comment []string
	// accumulator of a phrases
	str   strings.Builder
	words match.Span
}

func (pt *PlainText) Finalize() (ret []story.StoryStatement, err error) {
	if e := pt.flushAll(); e != nil {
		err = e
	} else {
		ret, pt.out = pt.out, nil
	}
	return
}

func (pt *PlainText) Decoded(p Pos, t Type, v any) (err error) {
	switch t {
	default:
		log.Println("unhandled token", t, "at", p)
		panic("unhandled token")

	case Comma:
		pt.flushComment()
		pt.writeHash(",", match.Keywords.Comma)

	case Comment:
		// fix? the tokenizer eats newlines;
		// to accurately reconstruct the block we'd want them.
		// and inline comments would be added to the declare
		// not jumped into a new comment op.
		str := v.(string)
		if e := pt.flushPhrases(nil); e != nil {
			err = e
		} else {
			pt.comment = append(pt.comment, str)
		}

	case Parenthetical:
		// fix: jess without flex needs to be able to match parens
		// the "easiest" way would be to merge / move token parsing into match
		str := v.(string)
		pt.flushComment()
		pt.writeHash(str, match.Hash("()"))

	case Quoted: // quoted string
		// fix: to preserve the phrases we need to know what kind of string it was
		// so we can reconstruct the quote markers...
		// maybe the token should be a "string literal" struct containing the original string
		// so there's no need to reconstruct
		str := v.(string)

		var terminal bool // fix? it'd make more sense to flag this in the reader
		if w, at := utf8.DecodeLastRuneInString(str); at > 0 {
			terminal = unicode.Is(unicode.Sentence_Terminal, w)
		}

		// fix? match could write/jess could read these as one element
		pt.flushComment()
		pt.writeHash("", match.Keywords.QuotedText)
		pt.writeHash(str, 0)
		if terminal {
			pt.endSentence()
		}

	case Stop:
		w := v.(rune)
		pt.flushComment()
		pt.str.WriteRune(w)
		pt.endSentence()

	case String:
		str := v.(string)
		pt.flushComment()
		pt.writeStr(str)

	case Tell:
		switch msg := v.(type) {
		case []any:
			if exe, e := decodeExecute(msg); e != nil {
				err = e
			} else {
				err = pt.flushPhrases(&assign.FromExe{
					Exe: exe,
				})
			}
		case map[string]any:
			// for now requires, From*:
			if a, e := decodeAssignment(msg); e != nil {
				err = e
			} else {
				err = pt.flushPhrases(a)
			}
		default:
			err = fmt.Errorf("can only handle tell sequences and maps, not %T", v)
		}
	}
	return
}

func (pt *PlainText) flushAll() (err error) {
	pt.flushComment()
	return pt.flushPhrases(nil)
}

func (pt *PlainText) endSentence() {
	pt.phrases = append(pt.phrases, pt.words)
	pt.words = nil
}

func (pt *PlainText) writeStr(str string) {
	pt.writeHash(str, match.Hash(str))
}

func (pt *PlainText) writeHash(str string, hash uint64) {
	pt.flushComment()
	// because we write words ( and other such things )
	// new text should have a space before;
	// terminals wouldn't but they dont come round here no more.
	if pt.str.Len() > 0 {
		pt.str.WriteRune(' ')
	}
	pt.str.WriteString(str)
	pt.words = append(pt.words, match.MakeWord(hash, str))
}

// if there are pending phrases, flush them
// ( ex. because a comment is about to be written )
func (pt *PlainText) flushPhrases(tail rt.Assignment) (err error) {
	if ks, ws := pt.phrases, pt.words; len(ks) == 0 && len(ws) == 0 {
		if tail != nil {
			err = fmt.Errorf("expected assignment to be part of a phrase")
		}
	} else {
		pt.phrases, pt.words = nil, nil
		// flush in progress words
		if len(ws) > 0 {
			ks = append(ks, ws)
		}
		// get the paragraph as a solid block of text
		// we've already written newlines and such
		str := pt.str.String()
		pt.str.Reset()
		// write the declare statement
		out := &story.DeclareStatement{
			Text:    &literal.TextValue{Value: str},
			Assign:  tail,
			Matches: ks,
		}
		pt.out = append(pt.out, out)
	}
	return
}

// if there are pending comments, flush them
// ( ex. because a phrase is about to be written )
func (pt *PlainText) flushComment() {
	if cs := pt.comment; len(cs) > 0 {
		pt.comment = nil
		out := &story.Comment{Lines: cs}
		pt.out = append(pt.out, out)
	}
}
