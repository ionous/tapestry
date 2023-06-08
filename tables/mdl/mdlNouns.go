package mdl

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

type nounInfo struct {
	id     int    // unique id of the noun
	name   string // validated name of the noun
	domain string // validated domain name

	kind     int    // kind's id
	fullpath string // full path
}

// if specified, must exist.
func (m *Writer) findOptionalNoun(domain, noun string, q nounFinder) (ret nounInfo, err error) {
	if len(noun) > 0 {
		ret, err = m.findRequiredNoun(domain, noun, q)
	}
	return
}

// if not specified errors, also errors if not found.
func (m *Writer) findRequiredNoun(domain, noun string, q nounFinder) (ret nounInfo, err error) {
	if out, e := m.findNoun(domain, noun, q); e != nil {
		err = e
	} else if out.id == 0 {
		err = errutil.Fmt("%w noun %q in domain %q", Missing, noun, domain)
	} else {
		ret = out
	}
	return
}

// if not specified errors, makes no assumptions about the results
func (m *Writer) findNoun(domain, noun string, q nounFinder) (ret nounInfo, err error) {
	if len(noun) == 0 {
		err = errutil.New("empty name for noun")
	} else if out, e := q(m.db, domain, noun); e != nil && e != sql.ErrNoRows {
		err = e
	} else {
		out.name = noun
		ret = out
	}
	return
}

type nounFinder func(db tables.QueryRow, domain, noun string) (ret nounInfo, err error)

func nounWithKind(db tables.QueryRow, domain, noun string) (ret nounInfo, err error) {
	err = db.QueryRow(`
	select mn.domain, mn.rowid, mk.rowid, ',' || mk.rowid || ',' || mk.path
	from mdl_noun mn
	join mdl_kind mk 
		on (mn.kind = mk.rowid)
	join domain_tree dt
		on (dt.uses = mn.domain)
	where base = ?1
	and noun = ?2
	limit 1`, domain, noun).Scan(&ret.domain, &ret.id, &ret.kind, &ret.fullpath)
	return
}

func nounSansKind(db tables.QueryRow, domain, noun string) (ret nounInfo, err error) {
	err = db.QueryRow(`
	select domain, mn.rowid
	from mdl_noun mn
	join domain_tree
		on (uses = domain)
	where base = ?1
	and noun = ?2
	limit 1`, domain, noun).Scan(&ret.domain, &ret.id)
	return
}
