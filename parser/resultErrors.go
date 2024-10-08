package parser

import "github.com/ionous/errutil"

type ErrorDepth interface {
	ErrorDepth() int
}

type Depth int

func DepthOf(e interface{}) (ret int) {
	if d, ok := e.(ErrorDepth); ok {
		ret = d.ErrorDepth()
	}
	return
}

type AmbiguousObject struct {
	Nouns []NounInstance
	Depth
}

type MismatchedWord struct {
	Want, Have string
	Depth
}

type MissingObject struct {
	Depth
}

// NoSuchObjects after asking for multiple items, and finding none.
type NoSuchObjects struct {
	Depth
}

// Overflow when we expect to be done, but input tokens remain.
type Overflow struct {
	Depth
}

// Underflow when we expect a word, but the input is empty
type Underflow struct {
	Depth
}

type UnknownObject struct {
	Depth
}

func (d Depth) ErrorDepth() int {
	return int(d)
}

func (a AmbiguousObject) Error() string {
	return errutil.Sprint("couldnt determine object", a.Nouns)
}
func (a MismatchedWord) Error() string {
	return errutil.Sprintf("mismatched word %s != %s at %d", a.Have, a.Want, a.Depth)
}
func (a MissingObject) Error() string {
	return errutil.Sprintf("missing an object at %d", a.Depth)
}
func (NoSuchObjects) Error() string {
	return "You can't see any such things."
}
func (Overflow) Error() string {
	return "too many words"
}
func (Underflow) Error() string {
	return "too few words"
}
func (UnknownObject) Error() string {
	return "You can't see any such thing."
}
