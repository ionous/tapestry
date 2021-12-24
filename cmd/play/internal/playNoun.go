package internal

import (
	"log"

	"git.sr.ht/~ionous/iffy/ident"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt/meta"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

type Noun struct {
	run *Playtime
	id  ident.Id
}

func MakeNoun(run *Playtime, name string) *Noun {
	return &Noun{run, ident.IdOf(name)}
}

// Id for the noun. Returned via ResultList.Objects() on a successful match.
func (n *Noun) Id() ident.Id {
	return n.id
}
func (n *Noun) String() string {
	return n.id.String()
}

func (n *Noun) HasPlural(s string) bool {
	return n.HasClass(s)
}

func (n *Noun) HasName(s string) bool {
	return n.run.HasName(n.String(), s)
}

func (n *Noun) HasClass(s string) (ret bool) {
	if ok, e := safe.IsKindOf(n.run, n.String(), lang.Underscore(s)); e != nil {
		log.Println("parser error", e)
	} else {
		ret = ok
	}
	return
}

// does the noun have the passed trait?
func (n *Noun) HasAttribute(s string) (ret bool) {
	if obj, e := n.run.GetField(meta.ObjectValue, n.String()); e != nil {
		log.Println("parser error", e)
	} else if p, e := obj.FieldByName(s); e != nil {
		log.Println("parser error", e)
	} else {
		ret = p.Bool()
	}
	return
}
