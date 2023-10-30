package charmed

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

type FloatMode int

//go:generate stringer -type=FloatMode
const (
	Pending FloatMode = iota
	Int10
	Int16
	Float64
)

// implements OperandState.
type NumParser struct {
	runes  charm.Runes
	mode   FloatMode
	negate bool
}

func (*NumParser) StateName() string {
	return "Numbers"
}

// returns int64 or float64
// func (p *NumParser) GetNumber() (ret any, err error) {
// 	switch s := p.runes.String(); p.mode {
// 	case Int10:
// 		ret = fromInt(s, p.negate)
// 	case Int16:
// 		ret = fromHex(s)
// 	case Float64:
// 		ret = fromFloat(s, p.negate)
// 	default:
// 		err = errutil.Fmt("unknown number: '%v' is %v.", s, p.mode)
// 	}
// 	return
// }

func (p *NumParser) GetFloat() (ret float64, err error) {
	switch s := p.runes.String(); p.mode {
	case Int10:
		ret = float64(fromInt(s, p.negate))
	case Int16:
		ret = float64(fromHex(s))
	case Float64:
		ret = fromFloat(s, p.negate)
	default:
		err = errutil.Fmt("unknown number: '%v' is %v.", s, p.mode)
	}
	return
}

// initial state of digit parsing.
// note: this doesn't support leading with just a "."
func (p *NumParser) NewRune(r rune) (ret charm.State) {
	switch {
	// in golang, leading +/- are unary operators;
	// here, they are considered optional parts decimal numbers.
	// note: strconv's base 10 parser doesnt handle leading signs.
	// we therefore leave them out of our result, and just flag the negative ones.
	case r == '-':
		p.negate = true
		fallthrough
	case r == '+':
		ret = charm.Statement("after lead plus", func(r rune) (ret charm.State) {
			if IsNumber(r) {
				p.mode = Int10
				ret = p.runes.Accept(r, charm.Statement("num plus", p.leadingDigit))
			}
			return
		})
	case r == '0':
		// 0 can standalone; but, it might be followed by a hex qualifier.
		p.mode = Int10
		ret = p.runes.Accept(r, charm.Statement("hex check", func(r rune) (ret charm.State) {
			// https://golang.org/ref/spec#hex_literal
			switch {
			case r == 'x' || r == 'X':
				p.mode = Pending
				ret = p.runes.Accept(r, charm.Statement("hex parse", func(r rune) (ret charm.State) {
					if IsHex(r) {
						p.mode = Int16
						ret = p.runes.Accept(r, charm.Statement("num hex", p.hexDigits))
					}
					return
				}))
			default:
				// delegate to number and dot checking...
				// in a statecharmed, it would be a super-state, and
				// x (above) would jump to a sibling of that super-state.
				ret = p.leadingDigit(r)
			}
			return
		}))
	case IsNumber(r):
		// https://golang.org/ref/spec#float_lit
		p.mode = Int10
		ret = p.runes.Accept(r, charm.Statement("num digits", p.leadingDigit))
	}
	return
}

// a string of numbers, possibly followed by a decimal or exponent separator.
// note: golang numbers can end in a pure ".", this does not allow that.
func (p *NumParser) leadingDigit(r rune) (ret charm.State) {
	switch {
	case IsNumber(r):
		ret = p.runes.Accept(r, charm.Statement("leading dig", p.leadingDigit))
	case r == '.':
		p.mode = Pending
		ret = p.runes.Accept(r, charm.Statement("decimal", func(r rune) (ret charm.State) {
			if IsNumber(r) {
				p.mode = Float64
				ret = p.runes.Accept(r, charm.Statement("decimal digits", p.leadingDigit))
			} else {
				ret = p.tryExponent(r) // delegate to exponent checking,,,
			}
			return
		}))
	default:
		ret = p.tryExponent(r) // delegate to exponent checking,,,
	}
	return
}

// https://golang.org/ref/spec#exponent
// exponent  = ( "e" | "E" ) [ "+" | "-" ] decimals
func (p *NumParser) tryExponent(r rune) (ret charm.State) {
	switch {
	case r == 'e' || r == 'E':
		p.mode = Pending
		ret = p.runes.Accept(r, charm.Statement("exp", func(r rune) (ret charm.State) {
			switch {
			case IsNumber(r):
				p.mode = Float64
				ret = p.runes.Accept(r, charm.Statement("exp decimal", p.decimals))
			case r == '+' || r == '-':
				ret = p.runes.Accept(r, charm.Statement("exp power", func(r rune) (ret charm.State) {
					if IsNumber(r) {
						p.mode = Float64
						ret = p.runes.Accept(r, charm.Statement("exp num", p.decimals))
					}
					return
				}))
			}
			return
		}))
	}
	return
}

// a chain of decimal digits 0-9
func (p *NumParser) decimals(r rune) (ret charm.State) {
	if IsNumber(r) {
		ret = p.runes.Accept(r, charm.Statement("decimals", p.decimals))
	}
	return
}

// a chain of hex digits 0-9, a-f
func (p *NumParser) hexDigits(r rune) (ret charm.State) {
	if IsHex(r) {
		ret = p.runes.Accept(r, charm.Statement("hexDigits", p.hexDigits))
	}
	return
}

func fromInt(s string, negate bool) (ret int64) {
	if i, e := strconv.ParseInt(s, 10, 64); e != nil {
		panic(e)
	} else if negate {
		ret = -i
	} else {
		ret = i
	}
	return
}

func fromHex(s string) (ret int64) {
	// hex string - chops out the 0x qualifier
	if i, e := strconv.ParseInt(s[2:], 16, 64); e != nil {
		panic(e)
	} else {
		ret = i // no negative for hex.
	}
	return
}

func fromFloat(s string, negate bool) (ret float64) {
	if f, e := strconv.ParseFloat(s, 64); e != nil {
		panic(e)
	} else if negate {
		ret = -f
	} else {
		ret = f
	}
	return
}
