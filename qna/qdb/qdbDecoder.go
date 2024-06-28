package qdb

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
)

// CommandDecoder transforms the raw bytes pulled from a query into in-memory commands.
type CommandDecoder interface {
	DecodeField(a affine.Affinity, b []byte, fieldType string) (literal.LiteralValue, error)
	DecodeAssignment(affine.Affinity, []byte) (rt.Assignment, error)
	DecodeProg([]byte) ([]rt.Execute, error)
	DecodeValue(typeinfo.Instance, []byte) error
}

// DecodeNone returns error for every method of Decoder.
// used for testing, and for a simplified
type DecodeNone string

// verify that the decoder implements every method
var _ CommandDecoder = DecodeNone("")

func (d DecodeNone) DecodeField(affine.Affinity, []byte, string) (_ literal.LiteralValue, err error) {
	err = errors.New("not implemented")
	return
}

func (d DecodeNone) DecodeAssignment(affine.Affinity, []byte) (_ rt.Assignment, err error) {
	err = errors.New("not implemented")
	return
}

func (d DecodeNone) DecodeProg([]byte) (_ []rt.Execute, err error) {
	err = errors.New("not implemented")
	return
}

func (d DecodeNone) DecodeValue(typeinfo.Instance, []byte) error {
	return errors.New("not implemented")
}
