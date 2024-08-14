package query

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
)

type NounInfo struct {
	Domain, Noun, Kind string // noun is unique identifier within the domain.
}

type RuleSet struct {
	Rules     []rt.Rule
	UpdateAll bool // true if any of the rules have update ste
}

type Query interface {
	IsDomainActive(name string) (bool, error)
	ActivateDomains(name string) (prev, next []string, err error)
	//
	GetKindByName(rawName string) (*rt.Kind, error)
	// given a plural or singular kind
	// return all ancestors starting with the kind itself.
	KindOfAncestors(kindOrKinds string) ([]string, error)
	// search using a short name
	NounInfo(shortname string) (NounInfo, error)
	// return the friendly name of the exact named noun
	NounName(fullname string) (string, error)
	// find the parser aliases for this noun
	// warning: the  parser expects these to be in alphabetical order.
	NounNames(fullname string) ([]string, error)
	// a single field can contain a set of recursive spare values;
	// so this returns "pairs" of paths and values.
	NounValue(fullname, field string) (rt.Assignment, error)
	// all nouns of the indicated kind
	NounsWithAncestor(kind string) ([]string, error)
	// the empty string if not found
	PluralToSingular(plural string) (string, error)
	// the empty string if not found
	PluralFromSingular(singular string) (string, error)
	// includes the parameters, followed by the result
	// the result can be a blank string for execute statements
	PatternLabels(pat string) ([]string, error)
	RulesFor(pat string) (RuleSet, error)
	ReciprocalsOf(rel, id string) ([]string, error)
	RelativesOf(rel, id string) ([]string, error)
	// relations can be cleared by passing a blank string on the opposite side
	// but -- fix -- there is no way to clear many-many relations.
	// errors if nothing changed.
	// doesnt check to see if the relation is valid;
	// the caller should do that.
	Relate(rel, noun, otherNoun string) error
	//
	LoadGame(path string) (CacheMap, error)
	SaveGame(path string, dynamicValues CacheMap) error
	Random(inclusiveMin, exclusiveMax int) int
	// release all resource
	Close()
}

func (n *NounInfo) IsValid() bool {
	return len(n.Noun) != 0
}

func (n *NounInfo) String() (ret string) {
	if !n.IsValid() {
		ret = "<unknown object>"
	} else {
		ret = strings.Join([]string{n.Domain, n.Noun}, "::")
	}
	return
}
