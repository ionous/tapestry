package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
)

// Decoder transforms the raw bytes pulled from a query into in-memory commands.
type Decoder interface {
	DecodeField(b []byte, a affine.Affinity, fieldType string) (rt.Assignment, error)
	DecodeAssignment(b []byte, a affine.Affinity) (rt.Assignment, error)
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

func (d DecodeNone) DecodeField(b []byte, a affine.Affinity, fieldType string) (ret rt.Assignment, err error) {
	err = NotImplemented(d)
	return
}

func (d DecodeNone) DecodeAssignment(b []byte, a affine.Affinity) (ret rt.Assignment, err error) {
	err = NotImplemented(d)
	return
}

func (d DecodeNone) DecodeFilter(b []byte) (ret rt.BoolEval, err error) {
	err = NotImplemented(d)
	return
}

func (d DecodeNone) DecodeProg(b []byte) (ret rt.Execute_Slice, err error) {
	err = NotImplemented(d)
	return
}
