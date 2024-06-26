package decoder

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
)

// Decoder transforms the raw bytes pulled from a query into in-memory commands.
type Decoder interface {
	DecodeField(a affine.Affinity, b []byte, fieldType string) (literal.LiteralValue, error)
	DecodeAssignment(affine.Affinity, []byte) (rt.Assignment, error)
	DecodeProg([]byte) ([]rt.Execute, error)
	DecodeValue(typeinfo.Instance, []byte) error
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

func (d DecodeNone) DecodeField(affine.Affinity, []byte, string) (_ literal.LiteralValue, err error) {
	err = NotImplemented(d)
	return
}

func (d DecodeNone) DecodeAssignment(affine.Affinity, []byte) (_ rt.Assignment, err error) {
	err = NotImplemented(d)
	return
}

func (d DecodeNone) DecodeProg([]byte) (_ []rt.Execute, err error) {
	err = NotImplemented(d)
	return
}

func (d DecodeNone) DecodeValue(typeinfo.Instance, []byte) error {
	return NotImplemented(d)
}
