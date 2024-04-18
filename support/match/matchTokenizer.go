package match

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ionous/tell/charm"
	"github.com/ionous/tell/charmed"
	"github.com/ionous/tell/runes"
)

// callback when a new token exists
// tbd: maybe a channel instead?
type Notifier interface {
	Decoded(TokenValue) error
}

// read pieces of plain text documents
type Tokenizer struct {
	Notifier
	curr, start Pos
	follows     bool
}

// special runes
const (
	runeComma  = ','
	runeOpener = '('
	runeCloser = ')'
)

func Tokenize(str string) (ret []TokenValue, err error) {
	var at Collector
	if e := at.Collect(str); e != nil {
		err = e
	} else {
		ret = at.Tokens
	}
	return
}

type Collector struct {
	Tokens       []TokenValue
	Lines        [][]TokenValue
	KeepComments bool
	BreakLines   bool
}

func (c *Collector) Collect(str string) (err error) {
	t := Tokenizer{Notifier: c}
	return charm.ParseEof(str, t.Decode())
}

func (at *Collector) Decoded(tv TokenValue) error {
	if at.KeepComments || tv.Token != Comment {
		if at.BreakLines && tv.Token == Stop {
			at.Lines = append(at.Lines, at.Tokens)
			at.Tokens = nil
		} else {
			at.Tokens = append(at.Tokens, tv)
		}
	}
	return nil
}

func NewTokenizer(n Notifier) charm.State {
	t := Tokenizer{Notifier: n}
	return t.Decode()
}

// return a state to parse a stream of runes and notify as they are detected.
func (n *Tokenizer) Decode() charm.State {
	return charm.Parallel("decode",
		n.readNext(),
		charmed.DecodePos(&n.curr.Y, &n.curr.X),
	)
}

// start searching for the next token
// starting with the next rune.
func (n *Tokenizer) readNext() charm.State {
	return charm.Step(eatWhitespace(), n.tokenize())
}

// start searching for the next token,
// starting with the passed rune.
func (n *Tokenizer) sendNext(q rune) charm.State {
	next := n.readNext()
	return next.NewRune(q)
}

func (n *Tokenizer) tokenize() charm.State {
	return charm.Statement("tokenize", func(q rune) (ret charm.State) {
		n.start = n.curr
		switch q {
		// fix? matchHash has/had a filter on bad punctuation
		// ex.  r != '-' && unicode.IsPunct(r)
		case runes.HTab:
			ret = charm.Error(errors.New("tabs are invalid whitespace"))

		case runes.Hash:
			ret = n.readComments() // unlike .tell; eats the hash

		case runes.InterpretQuote: // doublequote
			ret = n.readQuotes(false)

		case runes.RawQuote: // backtick
			ret = n.readQuotes(true)

		case runes.Colon:
			includeComments := true
			ret = DecodeDoc(includeComments, func(q rune, content any) (ret charm.State) {
				// after the async document has finished:
				// notifier the reader using a token
				if e := content.(error); e != nil {
					ret = charm.Error(e)
				} else if e := n.notifyToken(Tell, content); e != nil {
					ret = charm.Error(e)
				} else {
					// then start processing the rune which the document couldnt handle
					ret = n.sendNext(q)
				}
				return
			})

		case runeComma:
			if e := n.notifyToken(Comma, ","); e != nil {
				ret = charm.Error(e)
			} else {
				ret = n.readNext()
			}

		case runeOpener:
			ret = n.readParens()

		default:
			switch {
			case unicode.In(q, unicode.Terminal_Punctuation):
				if e := n.notifyToken(Stop, string(q)); e != nil {
					ret = charm.Error(e)
				} else {
					ret = n.readNext()
				}

			case unicode.IsPrint(q):
				next := n.wordDecoder()
				ret = next.NewRune(q)
			}
		}
		return
	})
}

// a single word -- roughly, letters and numbers ending with space (or eof)
// but also, commas, fullstops, and maybe a few other things.
func (n *Tokenizer) wordDecoder() charm.State {
	// fix: hash the string as its read, and send the pair of hash and string
	var b strings.Builder
	return charm.Self("roxanne", func(self charm.State, q rune) (ret charm.State) {
		// things that end words. words roxxane, words.
		fini := slices.Contains([]rune{
			runes.Newline,
			runes.Eof,
		}, q) || unicode.In(q, unicode.Space, unicode.Terminal_Punctuation)
		if fini {
			if e := n.notifyToken(String, b.String()); e != nil {
				ret = charm.Error(e)
			} else {
				ret = n.sendNext(q)
			}
		} else if !unicode.IsPrint(q) {
			ret = charm.Error(fmt.Errorf("unknown rune %c", q))
		} else {
			b.WriteRune(q)
			ret = self
		}
		return
	})
}

// read comments but strip leading hash
func (n *Tokenizer) readComments() charm.State {
	var b strings.Builder
	var gotSpace bool
	return charm.Self("comments", func(self charm.State, q rune) (ret charm.State) {
		switch q {
		default:
			if gotSpace {
				b.WriteRune(q)
				ret = self
			} else if q == runes.Space {
				gotSpace = true
				ret = self
			} else {
				ret = charm.Error(fmt.Errorf("expected a space after comment hash, not %c", q))
			}
		case runes.Newline, runes.Eof:
			if e := n.notifyToken(Comment, b.String()); e != nil {
				ret = charm.Error(e)
			} else {
				ret = n.sendNext(q)
			}
		}
		return
	})
}

// quoted string and heredoc decoder
// fix: remember the quoted string types ( for more/accurate reconstruction );
// flag terminal in the process ( ex. return a quoted struct with the relevant info )
func (n *Tokenizer) readQuotes(raw bool) charm.State {
	var d charmed.QuoteDecoder
	start := d.Interpret
	if raw {
		start = d.Record
	}
	return charm.Step(start(), charm.Statement("readQuotes", func(q rune) (ret charm.State) {
		if next, e := n.notifyQuotes(q, d.String()); e != nil {
			ret = charm.Error(e)
		} else {
			ret = next
		}
		return
	}))
}

func (n *Tokenizer) readParens() charm.State {
	var str strings.Builder
	var tailSpaces int
	return charm.Self("readParens", func(self charm.State, q rune) (ret charm.State) {
		switch q {
		case runeCloser:
			out := str.String()
			if e := n.notifyToken(Parenthetical, out[:len(out)-tailSpaces]); e != nil {
				ret = charm.Error(e)
			} else {
				ret = n.readNext()
			}

		case runeOpener, runes.Newline, runes.Eof, runes.HTab:
			e := charm.InvalidRune(q)
			ret = charm.Error(e)

		default:
			// skip initial spaces
			if space := q == runes.Space; !space || str.Len() > 0 {
				str.WriteRune(q)
				// count trailing spaces
				if !space {
					tailSpaces = 0
				} else {
					tailSpaces++
				}
			}
			ret = self // loop...
		}
		return
	})
}

// ---
// tell Notifier of the new token/value pair
func (n *Tokenizer) notifyToken(t Token, v any) error {
	tv := TokenValue{Token: t, Pos: n.start, Value: v, First: !n.follows}
	// when t is a Comment, make no decision about the start of a sentence.
	if t != Comment {
		// when t is Stop, the next token will be considered the start of a sentence.
		n.follows = t != Stop
	}
	return n.Notifier.Decoded(tv)
}

// send a quoted string, and
// if the string ends with a sentence terminal send that too.
func (n *Tokenizer) notifyQuotes(q rune, str string) (ret charm.State, err error) {
	if e := n.notifyToken(Quoted, str); e != nil {
		err = e
	} else {
		if w, at := utf8.DecodeLastRuneInString(str); at > 0 && unicode.Is(unicode.Sentence_Terminal, w) {
			err = n.notifyToken(Stop, string(w))
		}
		if err == nil {
			ret = n.sendNext(q)
		}
	}
	return
}
