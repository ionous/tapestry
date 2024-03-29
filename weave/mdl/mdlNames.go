package mdl

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

type NounMaker interface {
	AddNounKind(name, kind string) error
	AddNounName(noun, name string, rank int) error
}

// given an author specified name generate a new noun and its names
func AddNamedNoun(pen NounMaker, longName, kind string) (ret string, err error) {
	noun := inflect.Normalize(longName)
	if e := pen.AddNounKind(noun, kind); e == nil {
		names := weaver.MakeNames(longName)
		err = AddNounNames(pen, noun, names)
	} else if !errors.Is(e, Duplicate) {
		err = e
	}
	if err == nil {
		ret = noun
	}
	return
}

func AddNounNames(pen NounMaker, noun string, names []string) (err error) {
	for i, name := range names {
		if e := pen.AddNounName(noun, name, i); e != nil {
			err = e
			break
		}
	}
	return
}
