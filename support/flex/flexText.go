package flex

import (
	"fmt"
	"io"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"github.com/ionous/tell/charm"
)

// consumes all text until eof ( and eats the eof error )
// fix: allow line number offset
func ReadText(runes io.RuneReader) (ret []story.StoryStatement, err error) {
	var pt PlainText
	run := match.NewTokenizer(&pt)
	if e := charm.Read(runes, run); e != nil {
		err = e
	} else {
		ret, err = pt.Finalize()
	}
	return
}

// translate a plain text section to paragraphs of commands
// containing comments and jess Declare statements.
type PlainText struct {
	// declare statements, or comments
	out []story.StoryStatement
	// accumulator for declare, and comments
	phrases [][]match.TokenValue
	comment []string
	// accumulator of a phrases
	str   strings.Builder
	words []match.TokenValue
}

func (pt *PlainText) Finalize() (ret []story.StoryStatement, err error) {
	if e := pt.flushAll(); e != nil {
		err = e
	} else {
		ret, pt.out = pt.out, nil
	}
	return
}

// handler for match tokenizer;
func (pt *PlainText) Decoded(tv match.TokenValue) (err error) {
	switch tv.Token {
	default:
		if str, ok := tv.Value.(string); !ok {
			panic(tv.Token.String())
		} else {
			pt.writeToken(str, tv)
		}

	case match.Quoted:
		// fix: to build the proper declare, we'd need quote type
		// and/or original string; and would need to handle terminals.
		// currently we can't distinguish between "word." and "word".
		if str, ok := tv.Value.(string); !ok {
			panic(tv.Token.String())
		} else {
			pt.writeToken(`"`, tv)
			pt.str.WriteString(str)
			pt.str.WriteRune('"')
		}

	case match.Parenthetical:
		if str, ok := tv.Value.(string); !ok {
			panic(tv.Token.String())
		} else {
			pt.writeToken("( ", tv)
			pt.str.WriteString(str)
			pt.str.WriteString(") ")
		}

	case match.Comment:
		// fix? the tokenizer eats newlines;
		// to accurately reconstruct the block we'd want them.
		// and inline comments would be added to the declare
		// not jumped into a new comment op.
		if e := pt.flushPhrases(nil); e != nil {
			err = e
		} else {
			str := tv.Value.(string)
			pt.comment = append(pt.comment, str)
		}

	case match.Tell:
		// rewrite match.Tell tokens as assignments
		switch msg := tv.Value.(type) {
		case []any:
			if exe, e := decodeExecute(msg); e != nil {
				err = e
			} else {
				err = pt.writeAssignment(tv.Pos, &call.FromExe{
					Exe: exe,
				})
			}
		case map[string]any:
			// for now requires, From*:
			if a, e := decodeAssignment(msg); e != nil {
				err = e
			} else {
				err = pt.writeAssignment(tv.Pos, a)
			}
		default:
			err = fmt.Errorf("can only handle tell sequences and maps, not %T", msg)
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

func (pt *PlainText) writeAssignment(pos match.Pos, a rt.Assignment) error {
	pt.writeToken(":", match.TokenValue{
		Token: match.Tell,
		Pos:   pos,
		Value: a,
	})
	return pt.flushPhrases(a)
}

func (pt *PlainText) writeToken(str string, tv match.TokenValue) {
	pt.flushComment()
	// because we write words ( and other such things )
	// new text should have a space before;
	// terminals wouldn't but they dont come round here no more.
	if pt.str.Len() > 0 && tv.Token != match.Stop {
		pt.str.WriteRune(' ')
	}
	pt.str.WriteString(str)
	if tv.Token == match.Stop {
		pt.endSentence()
	} else {
		pt.words = append(pt.words, tv)
	}
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
		out := story.MakeDeclaration(str, tail, ks)
		pt.out = append(pt.out, out)
	}
	return
}

// if there are pending comments, flush them
// ( ex. because a phrase is about to be written )
func (pt *PlainText) flushComment() {
	if cs := pt.comment; len(cs) > 0 {
		pt.comment = nil
		out := &story.StoryNote{Text: strings.Join(cs, "\n")}
		pt.out = append(pt.out, out)
	}
}
