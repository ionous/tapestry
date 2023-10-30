package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/charmed"
	"github.com/ionous/errutil"
)

// a sequence of array values are specified with:
// a dash, inline whitespace, the value, and a tail ( whitespace with newlines. )
// assuming the tail is valid, loops back to itself to handle the next dash.
//
// like yaml, whitespace can be any combination of ascii space or ascii tab:
// https://yaml.org/spec/1.2.2/#white-space-characters
type SeqParser struct {
	expectedIndent int
	values         []Value
}

func (p *SeqParser) StateName() string {
	return "sequence"
}

func (p *SeqParser) GetValue() (ret []Value, err error) {
	// if p.lineCount == 0 {
	// 	err = errutil.New("expected a new line, none found")
	// } else {
	// 	retDepth, retLines = p.indent, p.lineCount
	// }
	ret = p.values
	return
}

func (p *SeqParser) NewRune(r rune) (ret charm.State) {
	if r == SequenceDash {
		ret = charm.Step(charmed.RequiredSpaces, charm.Statement("seq inline spaces", func(r rune) (ret charm.State) {
			var value ValueParser
			return charm.RunStep(r, &value, charm.Statement("seq value", func(r rune) (ret charm.State) {
				if val, e := value.GetValue(); e != nil {
					ret = charm.Error(e)
				} else {
					var tail TailParser
					ret = charm.RunStep(r, &tail, charm.Statement("seq tail", func(r rune) (ret charm.State) {
						if depth, lines := tail.GetTail(); lines > 0 && depth != p.expectedIndent {
							e := badIndent{depth}
							ret = charm.Error(e)
						} else if lines == 0 && r != charmed.Eof {
							e := errutil.New("expected a newline after a sequence value")
							ret = charm.Error(e)
						} else {
							// since the tail didn't parse, we're at the rune *after* the tail
							// for example, the next dash.... so loop.
							p.values = append(p.values, val)
							ret = p.NewRune(r)
						}
						return
					}))
				}
				return
			}))
		}))
	}
	return
}
