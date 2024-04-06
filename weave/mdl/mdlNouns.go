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
join domain_tree dt
	on (dt.uses = mp.domain)
join mdl_kind mk
  on (mk.rowid = mp.relKind)
join mdl_noun one
  on (one.rowid = mp.oneNoun)
join mdl_noun other
  on (other.rowid = mp.otherNoun)
where base = ?1
and relKind = ?2
and ((?4 and oneName = ?3) or (not ?4 and otherName=?3))`,
		pen.domain, relation, noun, primary); e != nil {
		err = e
	} else {
		ret, err = tables.ScanStrings(rows)
	}
	return
}

// return a specific field of a specific noun.
// this is a more limited form of the runtime version;
// it doesn't attempt to unpack records.
func (pen *Pen) GetNounValue(noun, field string) (ret []byte, err error) {
	if e := pen.db.QueryRow(`
		select mv.value
		from mdl_noun mn
		join domain_tree dt
			on(dt.uses = mn.domain)
		join mdl_value mv   
			on (mv.noun == mn.rowid)
		join mdl_field mf
			on (mv.field = mf.rowid)
		where base = ?1 
		and mn.noun = ?2
		and mf.field = ?3
		and dot is null`,
		pen.domain, noun, field).Scan(&ret); e != nil {
		if e != sql.ErrNoRows {
			err = e
		} else {
			err = errutil.Fmt("%w noun %q value %q", Missing, noun, field)
		}
	}
	return
}

// prefer runtime meta.ObjectId
func (pen *Pen) GetClosestNoun(name string) (ret MatchedNoun, err error) {
	var noun, kind string
	if e := pen.db.QueryRow(`
	select mn.noun, mk.kind
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
		Scan(&noun, &kind); e != nil && e != sql.ErrNoRows {
		err = e
	} else if e != nil {
		err = errutil.Fmt("%w closest noun %q in domain %q", Missing, name, pen.domain)
	} else {
		ret = MatchedNoun{Name: noun, Kind: kind, Match: name}
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
