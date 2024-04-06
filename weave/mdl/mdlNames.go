package mdl

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/support/inflect"
)

type NounMaker interface {
	AddNounKind(name, kind string) error
	AddNounName(noun, name string, rank int) error
}

// generate names for a noun;
// the standard rules use these for printing when no specific printed names exist;
// and for finding nouns that were specified by authors in a story.
func MakeNames(source string) (out []string) {
	out = append(out, source)
	if base := inflect.Normalize(source); base != source {
		out = append(out, base)
	}
	// generate additional names by splitting the name into parts
	split := inflect.Fields(source)
	if cnt := len(split); cnt > 1 {
		// in case the name was reduced due to multiple separators
		if breaks := strings.Join(split, " "); breaks != source {
			out = append(out, breaks)
		}
		// write individual words in increasing rank ( ex. "boat", then "toy" )
		// note: trailing words are considered "stronger"
		// because adjectives in noun names tend to be first ( ie. "toy boat" )
		for i := len(split) - 1; i >= 0; i-- {
			word := strings.ToLower(split[i])
			out = append(out, word)
		}
	}
	return out
}

// given an author specified name generate a new noun and its names
func AddNamedNoun(pen NounMaker, longName, kind string) (ret string, err error) {
	noun := inflect.Normalize(longName)
	if e := pen.AddNounKind(noun, kind); e == nil {
		names := MakeNames(longName)
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
