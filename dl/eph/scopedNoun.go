package eph

import (
	"math"
	"strings"

	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/lang"
	"github.com/ionous/errutil"
)

type ScopedNoun struct {
	Requires // kinds ( when resolved, can have one direct parent )
	domain   *Domain
	names    []string
	aliases  []string
	aliasat  []string // origin of each alias
	values   []NounValue
}

type NounValue struct {
	field string
	value literal.LiteralValue
	at    string
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

func (n *ScopedNoun) AddLiteralValue(field string, value literal.LiteralValue, at string) (err error) {
	if k, e := n.Kind(); e != nil {
		err = e
	} else if name, e := k.findCompatibleValue(field, value); e != nil {
		err = e
	} else {
		// the field was a trait, the returned name was an aspect
		if name != field {
			// redo the value we are setting as the trait of the aspect
			value = &literal.TextValue{
				Text: field,
			}
		}
		err = n.addLiteral(field, value, at)
	}
	return
}

// assumes the value is known to be compatible, and the field is a field... not a trait.
func (n *ScopedNoun) addLiteral(field string, value literal.LiteralValue, at string) (err error) {
	// verify we havent already stored a field of this value
	for _, q := range n.values {
		if q.field == field {
			why, was, wants := Redefined, q.field, field
			type stringer interface{ String() string }
			if try, ok := value.(stringer); ok {
				if curr, ok := q.value.(stringer); ok {
					if try, curr := try.String(), curr.String(); try == curr {
						was, wants, why = curr, try, Duplicated
					}
				}
			}
			err = &Conflict{
				Reason: why,
				Was:    Definition{q.at, was},
				Value:  wants,
			}
			break
		}
	}
	if err == nil {
		n.values = append(n.values, NounValue{field, value, at})
	}
	return
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
