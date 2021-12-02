package eph

import (
	"strings"

	"git.sr.ht/~ionous/iffy/lang"
	"github.com/ionous/errutil"
)

func (c *Catalog) WriteNouns(w Writer) (err error) {
	return forEachNoun(c, func(d *Domain, k *ScopedKind, n *ScopedNoun) (err error) {
		return w.Write(mdl_noun, d.name, n.name, k.name, n.at)
	})
}

func (c *Catalog) WriteNames(w Writer) (err error) {
	return forEachNoun(c, func(d *Domain, k *ScopedKind, n *ScopedNoun) (err error) {
		split := strings.FieldsFunc(n.name, lang.IsBreak)
		spaces := strings.Join(split, " ")

		// the ranked 0 name is used for default display when printing nouns
		// (ex. "toy boat")
		var ofs int
		breaks := n.name
		if e := w.Write(mdl_name, d.name, n.name, spaces, ofs, n.at); e != nil {
			err = e
		} else if cnt := len(split); cnt > 1 {
			// if there is more than one word...
			// these should never match... but that's how the old code was so why not...
			// ( ex. "toy_boat" )
			if spaces != breaks {
				ofs++
				if e := w.Write(mdl_name, d.name, n.name, breaks, ofs, n.at); e != nil {
					err = errutil.Append(err, e)
				}
			}
			// write individual words in increasing rank ( ex. "boat", then "toy" )
			// note: trailing words are considered "stronger"
			// because adjectives in noun names tend to be first ( ie. "toy boat" )
			for i := len(split) - 1; i >= 0; i-- {
				word, rank := split[i], cnt+ofs-i
				if e := w.Write(mdl_name, d.name, n.name, word, rank, n.at); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
		return
	})
}

func forEachNoun(c *Catalog, it func(*Domain, *ScopedKind, *ScopedNoun) error) (err error) {
	if ns, e := c.ResolveNouns(); e != nil {
		err = e
	} else {
		for _, ndep := range ns {
			n := ndep.Leaf().(*ScopedNoun)
			if k, e := n.Kind(); e != nil {
				err = errutil.Append(err, e)
			} else if e := it(n.domain, k, n); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

type NounError struct {
	Noun string
	Err  error
}

func (n NounError) Error() string {
	return errutil.Sprintf("%v for noun %q", n.Err, n.Noun)
}
func (n NounError) Unwrap() error {
	return n.Err
}

// name of a noun to assembly info
func (el *EphNouns) Phase() Phase { return NounPhase }

// noun, kind
func (el *EphNouns) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if name, ok := UniformString(el.Noun); !ok {
		err = InvalidString(el.Noun)
	} else if k, ok := UniformString(el.Kind); !ok {
		err = InvalidString(el.Kind)
	} else {
		noun := d.EnsureNoun(name, at)
		// we can only add requirements to the noun in the same domain that it was declared
		// if in a different domain: the nouns have to match up
		if noun.domain == d {
			noun.AddRequirement(k)
		} else if ok, e := noun.HasAncestor(k); e != nil {
			err = e
		} else if !ok {
			err = NounError{name, errutil.Fmt("can't redefine parent as %q", k)}
		} else {
			e := errutil.New("duplicate noun definition at", at)
			LogWarning(NounError{name, e})
		}
	}
	return
}
