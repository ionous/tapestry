package mdl

import (
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

var mdl_opposite = `insert into mdl_rev(domain, oneWord, otherWord, at) 
				values(?1, ?2, ?3, ?4), (?1, ?3, ?2, ?4)`

// domain captures the scope in which the pairing was defined.
// within that scope: the noun, relation, and otherNoun are all unique names --
// even if they are not unique globally, and even if they a broader/different scope than the pair's domain.
func (pen *Pen) AddOpposite(a, b string) (err error) {
	domain, at := pen.domain, pen.at
	if d, e := pen.findDomain(); e != nil {
		err = e
	} else if rows, e := pen.db.Query(
		`select oneWord, otherWord, domain
			from mdl_rev 
			join domain_tree
				on(uses=domain)
			where base = ?1`, d); e != nil {
		err = errutil.New("database error", e)
	} else {
		var x, y, from string
		if e := tables.ScanAll(rows, func() (err error) {
			// the testing is a bit weird so we handle it all app side
			if (x == a && y == b) || (x == b && y == a) {
				err = errutil.Fmt("%w opposite %q <=> %q in %q and %q", Duplicate, a, b, from, domain)
			} else if x == a || y == a || x == b || y == b {
				err = errutil.Fmt(
					"%w %q <=> %q defined as opposites in %q now %q <=> %q in %q",
					Conflict, x, y, from, a, b, domain)
			}
			return
		}, &x, &y, &from); e != nil {
			err = eatDuplicates(pen.warn, e)
		} else {
			// writes the opposite paring as well
			_, err = pen.db.Exec(mdl_opposite, d, a, b, at)
		}
	}
	return
}
func (pen *Pen) GetOpposite(word string) (ret string, err error) {
	if e := pen.db.QueryRow(`
	select otherWord 
from mdl_rev mv
join domain_tree dt
	on(dt.uses = mv.domain)
where base = ?1
and oneWord = ?2
limit 1`,
		pen.domain, word).Scan(&ret); e != nil {
		err = e
	}
	return
}
