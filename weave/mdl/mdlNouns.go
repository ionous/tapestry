package mdl

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

type nounInfo struct {
	id     int64  // unique id of the noun
	name   string // validated name of the noun
	domain string // validated domain name
	// optional:
	kid      int64 // kind's id
	kind     string
	fullpath string // full path of kind
}

func (n *nounInfo) class() classInfo {
	return classInfo{
		id:       n.kid,
		name:     n.kind,
		fullpath: n.fullpath,
	}
}

func (pen *Pen) GetRelativeNouns(noun, relation string, primary bool) (ret []string, err error) {
	if rows, e := pen.db.Query(`
	select one.noun as oneName, other.noun as otherName
from mdl_pair mp
join mdl_kind mk
  on (mk.rowid = mp.relKind)
join mdl_noun one
  on (one.rowid = mp.oneNoun)
join mdl_noun other
  on (other.rowid = mp.otherNoun)
where (relKind = ?1)
and ((?3 and oneName = ?2) or (not ?3 and otherName=?2))`,
		pen.domain, relation, noun, primary); e != nil {
		err = e
	} else {
		ret, err = tables.ScanStrings(rows)
	}
	return
}

// find the noun with the closest name in this scope;
// assumes the name is lower cased, with spaces normalized.
// skips aliases for the sake of backwards compatibility:
// there should be a difference between "a noun is known as"
// and "understand this word typed by the player as" -- and currently there's not.
func (pen *Pen) GetClosestNoun(name string) (ret string, err error) {
	if noun, e := pen.getClosestNoun(name); e != nil {
		err = e
	} else {
		ret = noun.name
	}
	return
}

func (pen *Pen) getClosestNoun(name string) (ret nounInfo, err error) {
	var out nounInfo
	if e := pen.db.QueryRow(`
	select mn.domain, mn.rowid, mn.noun, mk.rowid, ',' || mk.rowid || ',' || mk.path
	from mdl_name my
	join mdl_noun mn
		on (mn.rowid = my.noun)
	join mdl_kind mk 
		on (mn.kind = mk.rowid)
	join domain_tree dt
		on (dt.uses = my.domain)
	where base = ?1
	and my.name = ?2
	and my.rank >= 0
	order by my.rank, my.rowid asc
	limit 1`, pen.domain, name).
		Scan(&out.domain, &out.id, &out.name, &out.kind, &out.fullpath); e != nil && e != sql.ErrNoRows {
		err = e
	} else if out.id == 0 {
		err = errutil.Fmt("%w closest noun %q in domain %q", Missing, name, pen.domain)
	} else {
		ret = out
	}
	return
}

// if not specified errors, also errors if not found.
func (pen *Pen) findRequiredNoun(noun string, q nounFinder) (ret nounInfo, err error) {
	if out, e := pen.findNoun(noun, q); e != nil && e != sql.ErrNoRows {
		err = e
	} else if out.id == 0 {
		err = errutil.Fmt("%w required noun %q in domain %q", Missing, noun, pen.domain)
	} else {
		ret = out
	}
	return
}

// if not specified errors, makes no assumptions about the results
func (pen *Pen) findNoun(noun string, q nounFinder) (ret nounInfo, err error) {
	if len(noun) == 0 {
		err = errutil.New("empty name for noun")
	} else if out, e := q(pen.db, pen.domain, noun); e != nil && e != sql.ErrNoRows {
		err = e
	} else {
		out.name = noun
		ret = out
	}
	return
}

type nounFinder func(db *tables.Cache, domain, noun string) (ret nounInfo, err error)

func nounWithKind(db *tables.Cache, domain, noun string) (ret nounInfo, err error) {
	err = db.QueryRow(`
	select mn.domain, mn.rowid, mn.noun, mk.rowid, mk.kind, ',' || mk.rowid || ',' || mk.path
	from mdl_noun mn
	join mdl_kind mk 
		on (mn.kind = mk.rowid)
	join domain_tree dt
		on (dt.uses = mn.domain)
	where base = ?1
	and noun = ?2
	limit 1`, domain, noun).Scan(&ret.domain, &ret.id, &ret.name, &ret.kid, &ret.kind, &ret.fullpath)
	return
}

func nounSansKind(db *tables.Cache, domain, noun string) (ret nounInfo, err error) {
	err = db.QueryRow(`
	select domain, mn.rowid, mn.noun
	from mdl_noun mn
	join domain_tree
		on (uses = domain)
	where base = ?1
	and noun = ?2
	limit 1`, domain, noun).Scan(&ret.domain, &ret.id, &ret.name)
	return
}
