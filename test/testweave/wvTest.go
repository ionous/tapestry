package testweave

import (
	"database/sql"
	"math/rand"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/test/eph"
	"git.sr.ht/~ionous/tapestry/test/testdb"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func NewWeaver(name string) *TestWeave {
	return NewWeaverOptions(name, true)
}

func NewWeaverOptions(name string, shuffle bool) *TestWeave {
	db := testdb.Create(name)
	return &TestWeave{
		name:      name,
		db:        db,
		cat:       weave.NewCatalogWithWarnings(db, nil, mdl.LogWarning),
		noShuffle: !shuffle,
	}
}

func NewWeaverCatalog(name string, db *sql.DB, cat *weave.Catalog, shuffle bool) *TestWeave {
	return &TestWeave{
		name:      name,
		db:        db,
		cat:       cat,
		noShuffle: !shuffle,
	}
}

func OkayError(t *testing.T, e error, prefix string) (okay bool, err error) {
	if okay = e != nil && strings.HasPrefix(e.Error(), prefix); okay {
		t.Log("ok:", e)
	} else {
		err = e
	}
	return
}

type TestWeave struct {
	name string
	// queues the commands so we can MakeDomain without worrying about error handling
	// also allows shuffling the declarations (within a single domain)
	queue     []eph.Ephemera
	noShuffle bool
	db        *sql.DB
	cat       *weave.Catalog
}

func (dt *TestWeave) Close() {
	if dt.db != nil {
		dt.db.Close()
		dt.db = nil
	}
}

func (dt *TestWeave) MakeDomain(names []string, add ...eph.Ephemera) {
	n, req := names[0], names[1:]
	if !dt.noShuffle {
		// shuffle the incoming ephemera
		rand.Shuffle(len(add), func(i, j int) { add[i], add[j] = add[j], add[i] })
		// shuffle the order of domain dependencies
		rand.Shuffle(len(req), func(i, j int) { req[i], req[j] = req[j], req[i] })
	}
	dt.queue = append(dt.queue, &eph.BeginDomain{
		Name:     n,
		Requires: req,
	})
	dt.queue = append(dt.queue, add...)
	dt.queue = append(dt.queue, &eph.EndDomain{
		Name: n,
	})
}

func (dt *TestWeave) Assemble() (ret *weave.Catalog, err error) {
	for _, el := range dt.queue {
		if e := el.Assert(dt.cat); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		err = dt.cat.AssembleCatalog()
	}
	return dt.cat, err
}

func (dt *TestWeave) ReadAssignments() ([]string, error) {
	return mdl.ReadAssignments(dt.db)
}

// an array of dependencies
func (dt *TestWeave) ReadDomain(n string) (ret []string, err error) {
	return mdl.ReadDomain(dt.db, n)
}

func (dt *TestWeave) ReadDomains() ([]string, error) {
	return mdl.ReadDomains(dt.db)
}

// domain, kind, field name, affinity, subtype
// sorted by kind and within each kind, the field name
func (dt *TestWeave) ReadFields() ([]string, error) {
	return mdl.ReadFields(dt.db)
}

// domain, input, serialized program
func (dt *TestWeave) ReadGrammar() ([]string, error) {
	return mdl.ReadGrammar(dt.db)
}

// domain, kind, expanded materialized path
// ordered by domain, length of path, and name
// ( that erases their natural, dependency order --
//   but independent siblings dont otherwise have a consistent order for testing )
func (dt *TestWeave) ReadKinds() (ret []string, err error) {
	return mdl.ReadKinds(dt.db)
}

// domain, kind, name, serialized initialization
func (dt *TestWeave) ReadLocals() ([]string, error) {
	return mdl.ReadLocals(dt.db)
}

// domain, noun, name, rank
func (dt *TestWeave) ReadNames() ([]string, error) {
	return mdl.ReadNames(dt.db)
}

// domain, noun, kind
func (dt *TestWeave) ReadNouns() ([]string, error) {
	return mdl.ReadNouns(dt.db)
}

// domain, oneWord, otherWord
func (dt *TestWeave) ReadOpposites() ([]string, error) {
	return mdl.ReadOpposites(dt.db)
}

// domain, relation, noun, other noun
// original order was domain, and alpha relation name
func (dt *TestWeave) ReadPairs() ([]string, error) {
	return mdl.ReadPairs(dt.db)
}

// domain, pattern, labels, result field
func (dt *TestWeave) ReadPatterns() ([]string, error) {
	return mdl.ReadPatterns(dt.db)
}

// domain, many, one
func (dt *TestWeave) ReadPlurals() ([]string, error) {
	return mdl.ReadPlurals(dt.db)
}

// domain, relation, one kind, other kind, cardinality
// ordered by name of the relation for test consistency.
func (dt *TestWeave) ReadRelations() ([]string, error) {
	return mdl.ReadRelations(dt.db)
}

// domain, pattern, target, phase, filter, prog
func (dt *TestWeave) ReadRules() ([]string, error) {
	return mdl.ReadRules(dt.db)
}

// domain, noun, field, value
func (dt *TestWeave) ReadValues() ([]string, error) {
	return mdl.ReadValues(dt.db)
}
