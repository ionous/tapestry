package internal

import (
	"log"
	"strings"

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
	if objectPath, e := n.run.GetField(meta.Kinds, n.String()); e != nil {
		log.Println("parser error", e)
	} else {
		// Contains reports whether second is within first.
		kind := lang.Underscore(s)
		cp, ck := objectPath.String()+",", kind+","
		ret = strings.Contains(cp, ck)
	}
	return
}

// does the noun have the passed trait?
func (n *Noun) HasAttribute(s string) (ret bool) {
	if obj, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	}
	if p, e := n.run.GetField(n.String(), s); e != nil {
		log.Println("parser error", e)
	} else {
		ret = p.Bool()
	}
	return
}
