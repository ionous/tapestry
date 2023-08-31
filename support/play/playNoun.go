package play

import (
	"log"
	"sort"

	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// these are created by the survey
type Noun struct {
	run rt.Runtime
	id  string
}

func MakeNoun(run rt.Runtime, name string) *Noun {
	return &Noun{run, lang.Normalize(name)}
}

// Id for the noun. Returned via ResultList.Objects() on a successful match.
func (n *Noun) String() string {
	return n.id
}

func (n *Noun) HasPlural(s string) bool {
	return n.HasClass(s)
}

func (n *Noun) HasName(name string) (okay bool) {
	if name == n.id {
		okay = true
	} else if ugh, e := n.run.GetField(meta.ObjectAliases, n.id); e == nil {
		ick := ugh.Strings()
		if i := sort.SearchStrings(ick, name); i < len(ick) && ick[i] == name {
			okay = true
		}
	}
	return
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
