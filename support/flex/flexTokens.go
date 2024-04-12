package flex

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/ionous/tell/charm"
	"github.com/ionous/tell/charmed"
	"github.com/ionous/tell/runes"
)

// position of a token
type Pos struct{ X, Y int }

//go:generate stringer -type=Type
type Type int

// types of tokens
const (
	Invalid Type = iota // placeholder, not generated by the tokenizer
	Word                // loosely defined here
	Comma               // a literal comma
	Stop                // full stop or other terminal
	Comment             // # something
	Quoted              // quoted string
	TellBlock

	// Number? but need to handle trailing decimals: "the age is 45."
)

// special runes
const (
	runeComma = ','
)

// callback when a new token exists
// tbd: maybe a channel instead?
type Notifier interface {
	Decoded(Pos, Type, any) error
}

// read pieces of plain text documents
type Tokenizer struct {
	Notifier
	curr, start Pos
}

func NewTokenizer(n Notifier) charm.State {
	t := Tokenizer{Notifier: n}
	return t.Decode()
}

// return a state to parse a stream of runes and notify as they are detected.
func (n Tokenizer) Decode() charm.State {
	return charm.Parallel("decode", n.decode(false), charmed.DecodePos(&n.curr.Y, &n.curr.X))
}

func (n *Tokenizer) decode(afterIndent bool) charm.State {
	return charm.Step(n.whitespace(afterIndent), n.tokenize())
}

// tell Notifier of the new token/value pair
// and then process the next rune (q)
// ( combining the two simplifies error handling in some cases )
func (n *Tokenizer) notifyRune(q rune, t Type, v any) (ret charm.State) {
	if e := n.Notifier.Decoded(n.start, t, v); e != nil {
		ret = charm.Error(e)
	} else {
		ret = send(n.decode(true), q)
	}
	return
}

// tell Notifier of the new token/value pair
func (n *Tokenizer) notifyLoop(t Type, v any) (ret charm.State) {
	if e := n.Notifier.Decoded(n.start, t, v); e != nil {
		ret = charm.Error(e)
	} else {
		ret = n.decode(true)
	}
	return
}

// eat whitespace between tokens;
// previously, would error if it didnt detect whitespace between tokens
// however that doesnt work well for arrays. ex: `5,`
func (n *Tokenizer) whitespace(afterIndent bool) charm.State {
	var spaces int
	return charm.Self("whitespace", func(self charm.State, q rune) (ret charm.State) {
		switch q {
		case runes.Space:
			spaces++
			ret = self
		case runes.Newline:
			spaces = 0
			afterIndent = false
			ret = self
		case runes.Eof:
			ret = charm.Error(nil)
		}
		return
	})
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
			ret = n.commentDecoder() // unlike .tell; eats the hash

		case runes.InterpretQuote: // doublequote
			ret = n.interpretDecoding()

		case runes.RawQuote: // backtick... maybe.
			ret = n.rawDecoding()

		case runes.Colon:
			ret = charm.Error(errors.New("FIX: handle tell doc"))

		case runeComma:
			ret = n.notifyLoop(Comma, q)

		default:
			switch {
			case unicode.In(q, unicode.Terminal_Punctuation):
				ret = n.notifyLoop(Stop, q)

			case unicode.IsPrint(q):
				next := n.wordDecoder()
				ret = send(next, q)
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
			ret = n.notifyRune(q, Word, b.String())
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
func (n *Tokenizer) commentDecoder() charm.State {
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
			ret = n.notifyRune(q, Comment, b.String())
		}
		return
	})
}

// quoted string and heredoc decoder
// fix: remember the quoted string types ( for more/accurate reconstruction );
// flag terminal in the process ( ex. return a quoted struct with the relevant info )
func (n *Tokenizer) interpretDecoding() charm.State {
	var d charmed.QuoteDecoder
	return charm.Step(d.Interpret(), charm.Statement("interpreted", func(q rune) charm.State {
		return n.notifyRune(q, Quoted, d.String())
	}))
}

func (n *Tokenizer) rawDecoding() charm.State {
	var d charmed.QuoteDecoder
	return charm.Step(d.Record(), charm.Statement("recorded", func(q rune) charm.State {
		return n.notifyRune(q, Quoted, d.String())
	}))
}

// send one or more runes to the next state
func send(next charm.State, qs ...rune) charm.State {
	for _, q := range qs {
		if next = next.NewRune(q); next == nil {
			break
		}
	}
	return next
}
