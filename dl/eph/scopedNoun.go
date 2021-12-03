package eph

import (
	"math"
	"strings"

	"git.sr.ht/~ionous/iffy/lang"
	"github.com/ionous/errutil"
)

type ScopedNoun struct {
	Requires // kinds ( when resolved, can have one direct parent )
	domain   *Domain
	names    []string
	aliases  []string
	aliasat  []string // origin of each alias
}

func (n *ScopedNoun) Resolve() (ret Dependencies, err error) {
	if len(n.at) == 0 {
		err = NounError{n.name, errutil.New("never defined")}
	} else if ks, e := n.resolve(n, (*kindFinder)(n.domain)); e != nil {
		err = NounError{n.name, e}
	} else {
		ret = ks
	}
	return
}

func (n *ScopedNoun) Kind() (ret *ScopedKind, err error) {
	if dep, e := n.GetDependencies(); e != nil {
		err = e
	} else if ks := dep.Parents(); len(ks) != 1 {
		err = errutil.Fmt("noun %q has unexpected %d parents", n.name, len(ks))
	} else {
		ret = ks[0].(*ScopedKind)
	}
	return
}

func (n *ScopedNoun) AddAlias(a, at string) {
	n.aliases = append(n.aliases, a)
	n.aliasat = append(n.aliasat, at)
}

func (n *ScopedNoun) Names() []string {
	if len(n.names) == 0 {
		n.names = n.makeNames()
	}
	return n.names
}

const UnknownRank = math.MaxInt

func (n *ScopedNoun) FindName(name string, rank int) {
	rank = UnknownRank
	for i, el := range n.Names() {
		if el == name && i < rank {
			rank = i
			break
		}
	}
	return
}

func (n *ScopedNoun) makeNames() []string {
	var out []string
	split := strings.FieldsFunc(n.name, lang.IsBreak)
	spaces := strings.Join(split, " ")

	// the ranked 0 name is used for default display when printing nouns
	// (ex. "toy boat")
	breaks := n.name
	out = append(out, spaces)
	if cnt := len(split); cnt > 1 {
		// if there is more than one word...
		// these should never match... but that's how the old code was so why not...
		// ( ex. "toy_boat" )
		if spaces != breaks {
			out = append(out, breaks)
		}
		// write individual words in increasing rank ( ex. "boat", then "toy" )
		// note: trailing words are considered "stronger"
		// because adjectives in noun names tend to be first ( ie. "toy boat" )
		for i := len(split) - 1; i >= 0; i-- {
			word := split[i]
			out = append(out, word)
		}
	}
	return out
}
