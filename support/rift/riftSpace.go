package rift

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

type AfterWhitespace func(newIndent int) charm.State

// func OptionalSpaces(name string, indent int, after AfterWhitespace) charm.State {
// 	return eatSpaces(name, indent, wsOptional, after)
// }

func CommentSpaces(name string, indent int, after AfterWhitespace) charm.State {
	// var b strings.Builder
	// b.WriteRune(VTab) // this should be optional
	return eatSpaces(name, indent, wsRequriesSpaces, after)
}

func CommentLines(name string, indent int, after AfterWhitespace) charm.State {
	// var b strings.Builder
	// b.WriteRune(HTab) // this should be optional
	return eatSpaces(name, indent, wsRequiresLines, after)
}

func eatSpaces(name string, indent int, flags wsFlags, after AfterWhitespace) charm.State {
	ws := whitespace{indent: indent, flags: flags}
	return charm.Step(&ws, charm.Statement(name, func(r rune) charm.State {
		// r is the fist non-whitespace character
		return charm.RunState(r, after(ws.indent))
	}))
}

type wsFlags int

const (
	wsOptional wsFlags = iota
	wsRequriesSpaces
	wsRequiresLines
)

// eats ascii whitespace, tracking indent
// in yaml, spaces are indents and after indentation tabs are allowed.
// that seems mildly interesting for end of line alignment of comments
// but sticking to no tabs at all seems even better.
type whitespace struct {
	indent   int     // number of spaces since starting to read the current line.
	lines    int     // number of newlines encountered
	flags    wsFlags // errors unless if there's no whitespace ( eof is also okay )
	comments strings.Builder
}

func (ws *whitespace) IsEmpty() bool {
	return ws.indent+ws.lines == 0
}

// first character of the signature must be a letter
func (ws *whitespace) NewRune(r rune) (ret charm.State) {
	switch r {
	case Newline:
		ws.lines++
		ws.indent = 0
		ret = ws
	case Space:
		ws.indent++
		ret = ws
	default:
		if ws.flags == wsRequiresLines && ws.lines == 0 {
			e := errutil.New("expected a new line")
			ret = charm.Error(e)
		} else if ws.flags == wsRequriesSpaces && ws.IsEmpty() {
			e := errutil.New("expected whitespace")
			ret = charm.Error(e)
		}
	}
	return
}
