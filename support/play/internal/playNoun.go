package internal

import (
	"log"

	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

type Noun struct {
	run *Playtime
	id  string
}

func MakeNoun(run *Playtime, name string) *Noun {
	return &Noun{run, lang.Normalize(name)}
}

// Id for the noun. Returned via ResultList.Objects() on a successful match.
func (n *Noun) String() string {
	return n.id
}

func (n *Noun) HasPlural(s string) bool {
	return n.HasClass(s)
}

func (n *Noun) HasName(s string) bool {
	return n.run.HasName(n.String(), s)
}

func (n *Noun) HasClass(s string) (ret bool) {
	if ok, e := safe.IsKindOf(n.run, n.String(), lang.Normalize(s)); e != nil {
		log.Println("parser error", e)
	} else {
		ret = ok
	}
	return
}

// does the noun have the passed trait?
func (n *Noun) HasAttribute(s string) (ret bool) {
	if p, e := n.run.GetField(n.String(), s); e != nil {
		log.Println("parser error", e)
	} else {
		ret = p.Bool()
	}
	return
}