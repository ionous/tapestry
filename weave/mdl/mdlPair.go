package mdl

import (
	"database/sql"

	"github.com/ionous/errutil"
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
	if e := pen.db.QueryRow(q, domain, rel.id, search.id).Scan(&prevId, &prevString); e != nil && e != sql.ErrNoRows {
		err = e
	} else if prevId.Valid {
		if prevId.Int64 == match.id {
			err = errutil.Fmt("%w relation %q duplicated %q to %q in domain %q",
				Duplicate, rel.name, one.name, other.name, domain)
		} else if !multi {
			err = errutil.Fmt("%w new relation %q of %q to %q in domain %q; was %q to %q",
				Conflict, rel.name, one.name, other.name, domain,
				one.name, prevString.String)
		}
	}
	return
}

func (pen *Pen) addPair(kind kindInfo, one, other nounInfo) (err error) {
	domain, at := pen.domain, pen.at
	_, err = pen.db.Exec(mdl_pair, domain, kind.id, one.id, other.id, at)
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