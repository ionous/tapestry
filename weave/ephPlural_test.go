package weave

import (
	"errors"
	"testing"

	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/test/testdb"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// catalog some plural ephemera from different domain levels
// and verify things wind up in the right place
func TestPluralConflict(t *testing.T) {
	var dt domainTest
	defer dt.Close()
	dt.makeDomain(dd("a"),
		// one singular can have several plurals:
		// exTestAncestryMultipleParents. "person" can be "people" or "persons".
		// but the same plural "persons" cant have multiple singular definitions
		&eph.Plurals{Singular: "raven", Plural: "unkindness"},
		&eph.Plurals{Singular: "witch", Plural: "unkindness"},
	)
	db := testdb.Open(t.Name(), testdb.Memory, "")
	defer db.Close()
	cat := NewCatalog(db)
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else {
		err := cat.AssembleCatalog()
		if e := okDomainConflict("a", Redefined, err); e != nil {
			t.Fatal(e)
		}
	}
}

// catalog some plural ephemera from different domain levels
// and verify things wind up in the right place
func TestPluralAssembly(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()
	// because this test picks out two warnings, one by one...
	// we cant shuffle the statements...
	dt := domainTest{noShuffle: true}
	defer dt.Close()
	// yes, these are collective nouns not plurals... shhh...
	dt.makeDomain(dd("a"),
		&eph.Plurals{Singular: "raven", Plural: "unkindness"},
		// one singular can have several plurals:
		// ex. "person" can be "people" or "persons".
		&eph.Plurals{Singular: "bat", Plural: "cloud"},
		&eph.Plurals{Singular: "bat", Plural: "cauldron"},
	)
	dt.makeDomain(dd("b", "a"),
		// add something new:
		&eph.Plurals{Singular: "fish", Plural: "school"},
		// collapse:
		&eph.Plurals{Singular: "bat", Plural: "cauldron"},
	)
	dt.makeDomain(dd("c", "a"),
		// redefine:
		&eph.Plurals{Singular: "witch", Plural: "unkindness"},
	)
	//
	db := testdb.Open(t.Name(), testdb.Memory, "")
	defer db.Close()
	cat := NewCatalog(db)
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else {
		err := cat.AssembleCatalog()
		if e := okDomainConflict("c", Redefined, err); e != nil {
			t.Fatal(e)
		} else if e := okDomainConflict("b", Duplicated, warnings.shift()); e != nil {
			t.Fatal(e)
		} else {
			if out, e := tables.ScanStrings(db, readPlurals); e != nil {
				t.Fatal(e)
			} else {
				if diff := pretty.Diff(out, []string{
					"a:unkindness:raven",
					"a:cloud:bat",
					"a:cauldron:bat",
					"b:school:fish",
					// plural redefinition is (no longer) allowed.
					// ( wicca good and love the earth: and i'll be over here. )
					// "c:unkindness:witch:x",
					// we dont expect to see our duplicated definition of cauldron of bat(s)
					// c is dependent on a: so the definition would be redundant.
					// "c:cauldron:bat:x",
				}); len(diff) > 0 {
					t.Log("got", len(out), out)
					t.Fatal(diff)
				}
			}
		}
	}
}

func okDomainConflict(d string, y ReasonForConflict, e error) (err error) {
	var de domainError
	var conflict *Conflict
	if !errors.As(e, &de) || de.Domain != d ||
		!errors.As(de.Err, &conflict) || conflict.Reason != y {
		err = errutil.New("unexpected conflict in", de.Domain, e)
	}
	return
}

var readPlurals = `
select md.domain ||':'|| mp.many ||':'|| mp.one
from mdl_plural mp 
join mdl_domain md 
where md.rowid == mp.domain`
