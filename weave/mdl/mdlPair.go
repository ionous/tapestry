package mdl

import (
	"database/sql"

	"fmt"
)

func (pen *Pen) checkPair(rel kindInfo, one, other nounInfo, reverse, multi bool) (err error) {
	var prevId sql.NullInt64
	var prevString sql.NullString
	var search, match nounInfo
	var q string
	if !reverse {
		q = forwardPairs
		search, match = other, one
	} else {
		q = reversePairs
		search, match = one, other
	}
	domain := pen.domain
	if e := pen.db.QueryRow(q, domain, rel.row, search.id).Scan(&prevId, &prevString); e != nil && e != sql.ErrNoRows {
		err = e
	} else if prevId.Valid {
		if prevId.Int64 == match.id {
			err = fmt.Errorf("%w relation %q duplicated %q to %q in domain %q",
				ErrDuplicate, rel.name, one.name, other.name, domain)
		} else if !multi {
			err = fmt.Errorf("%w new relation %q of %q to %q in domain %q; was %q to %q",
				ErrConflict, rel.name, one.name, other.name, domain,
				one.name, prevString.String)
		}
	}
	return
}

// for a given rhs, there can be only one lhs
var forwardPairs = `
	select mn.rowid, mn.noun
	from mdl_pair mp
	join mdl_noun mn
		on(mn.rowid = mp.oneNoun)
	where mp.domain = ?1 
	and mp.relKind = ?2 
	and mp.otherNoun = ?3`

// for a given lhs, there can be only one rhs
var reversePairs = `
	select mn.rowid, mn.noun
	from mdl_pair mp
	join mdl_noun mn
		on(mn.rowid = mp.otherNoun)
	where mp.domain = ?1 
	and mp.relKind = ?2 
	and mp.oneNoun = ?3`
