package decoder

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
)

// Decoder transforms the raw bytes pulled from a query into in-memory commands.
type Decoder interface {
	DecodeField(b []byte, a affine.Affinity, fieldType string) (literal.LiteralValue, error)
	DecodeAssignment(b []byte) (rt.Assignment, error)
	DecodeFilter(b []byte) (rt.BoolEval, error)
	DecodeProg(b []byte) (rt.Execute_Slice, error)
}

// NotImplemented - generic error used returned by QueryNone
type NotImplemented string

func (e NotImplemented) Error() string {
	return string(e)
}

// DecodeNone returns error for every method of Decoder.
type DecodeNone string

// verify that the decoder implements every method
var _ = Decoder(DecodeNone(""))

func (d DecodeNone) DecodeField(b []byte, a affine.Affinity, fieldType string) (_ literal.LiteralValue, err error) {
	err = NotImplemented(d)
	return
}

func (d DecodeNone) DecodeAssignment(b []byte) (_ rt.Assignment, err error) {
	err = NotImplemented(d)
	return
}

func (d DecodeNone) DecodeFilter(b []byte) (_ rt.BoolEval, err error) {
	err = NotImplemented(d)
	return
}

func (d DecodeNone) DecodeProg(b []byte) (_ rt.Execute_Slice, err error) {
	err = NotImplemented(d)
	return
}
