package internal

import (
	r "reflect"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/composer"
)

// groups are the eventual .proto files
type Groups struct {
	specs       map[string][]composer.Composer
	groupBySlat map[string]string
}

func (gs *Groups) GroupBySlat(n string) string {
	return gs.groupBySlat[n]
}

func (gs *Groups) add(c composer.Composer) string {
	gn := c.Compose().Group
	if len(gn) == 0 {
		gn = "misc"
	}
	els := gs.specs[gn]
	els = append(els, c)
	gs.specs[gn] = els
	//
	gs.groupBySlat[r.TypeOf(c).Elem().Name()] = gn
	return gn
}

func makeAll() Groups {
	groups := Groups{
		specs:       make(map[string][]composer.Composer),
		groupBySlat: make(map[string]string),
	}
	for _, slats := range iffy.AllSlats {
		for _, slat := range slats {
			groups.add(slat)
		}
	}
	return groups
}

var AllGroups Groups = makeAll()
