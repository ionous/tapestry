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
	Id           string // really an id, but we'll let the driver convert
	Phase        int
	Filter, Prog []byte
}

type Query interface {
	IsDomainActive(name string) (okay bool, err error)
	ActivateDomain(name string) (ret string, err error)
	ReadChecks(actuallyJustThisOne string) (ret []CheckData, err error)
	FieldsOf(kind string) (ret []FieldData, err error)
	KindOfAncestors(kind string) ([]string, error)
	NounInfo(name string) (ret NounInfo, err error)
	NounIsNamed(id, name string) (ret bool, err error)
	NounName(id string) (ret string, err error)
	NounValue(id, field string) (retVal []byte, err error)
	NounsByKind(kind string) ([]string, error)
	PluralToSingular(plural string) (ret string, err error)
	PluralFromSingular(singular string) (ret string, err error)
	OppositeOf(word string) (ret string, err error)
	PatternLabels(pat string) (ret []string, err error)
	RulesFor(pat, target string) (ret []Rules, err error)
	ReciprocalsOf(rel, id string) ([]string, error)
	RelativesOf(rel, id string) ([]string, error)
	Relate(rel, noun, otherNoun string) (err error)
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
