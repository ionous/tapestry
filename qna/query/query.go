package query

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
)

type CheckData struct {
	Name   string
	Domain string
	Aff    affine.Affinity
	Prog   []byte
	Value  []byte
}

type FieldData struct {
	Name     string
	Affinity affine.Affinity
	Class    string
	Init     []byte
}

type NounInfo struct {
	Domain, Id, Kind string // id is the string identifier for the noun, unique within the domain.
}

type Rules struct {
	Id                  string // really an id, but we'll let the driver convert
	Phase               int
	Updates, Terminates bool
	Filter, Prog        []byte
}

type Query interface {
	IsDomainActive(name string) (bool, error)
	ActivateDomain(name string) (string, error)
	ReadChecks(actuallyJustThisOne string) ([]CheckData, error)
	FieldsOf(kind string) ([]FieldData, error)
	KindOfAncestors(kind string) ([]string, error)
	NounInfo(name string) (NounInfo, error)
	NounIsNamed(id, name string) (bool, error)
	NounName(id string) (string, error)
	// a single field can contain a set of recursive spare values;
	// so this returns pairs of path, value.
	NounValues(id, field string) ([]string, error)
	NounsByKind(kind string) ([]string, error)
	PluralToSingular(plural string) (string, error)
	PluralFromSingular(singular string) (string, error)
	OppositeOf(word string) (string, error)
	// includes the parameters, followed by the result
	// the result can be a blank string for execute statements
	PatternLabels(pat string) ([]string, error)
	RulesFor(pat, target string) ([]Rules, error)
	ReciprocalsOf(rel, id string) ([]string, error)
	RelativesOf(rel, id string) ([]string, error)
	Relate(rel, noun, otherNoun string) error
}

func (n *NounInfo) IsValid() bool {
	return len(n.Id) != 0
}

func (n *NounInfo) String() (ret string) {
	if !n.IsValid() {
		ret = "<unknown object>"
	} else {
		ret = strings.Join([]string{n.Domain, n.Id}, "::")
	}
	return
}
