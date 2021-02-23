package internal

import (
	"log"

	"git.sr.ht/~ionous/iffy/ident"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/object"
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
	if objectPath, e := n.run.GetField(object.Kinds, n.String()); e != nil {
		log.Println("parser error", e)
	} else {
		// Contains reports whether second is within first.
		kind := lang.Breakcase(s)
		cp, ck := objectPath.String()+",", kind+","
		ret = cp == ck
	}
	return
}

// does the noun have the passed name attribute?
func (n *Noun) HasAttribute(s string) (ret bool) {
	if p, e := n.run.GetField(n.String(), s); e != nil {
		log.Println("parser error", e)
	} else {
		ret = p.Bool()
	}
	return
}
