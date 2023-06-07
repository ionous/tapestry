package mdl

import (
	"database/sql"

	"github.com/ionous/errutil"
)

type nounInfo struct {
	id       int    // unique id of the noun
	name     string // validated name of the noun
	domain   string // validated domain name
	fullpath string // full path of the kind
}

func (m *Writer) findNoun(domain, noun string) (retDomain string, retKind int, err error) {
	if len(noun) == 0 {
		err = errutil.New("missing a name for a noun")
	} else if e := m.db.QueryRow(`
	select domain, mn.rowid
	from mdl_noun mn
	join domain_tree
		on (uses = domain)
	where base = ?1
	and noun = ?2
	limit 1`, domain, noun).Scan(&retDomain, &retKind); e == sql.ErrNoRows {
		err = errutil.Fmt("%w noun %q in domain %q", Unknown, noun, domain)
	} else {
		err = e
	}
	return
}

// optional in two respects: the passed name can be empty;
// and the passed name may not exist in the database.
func (m *Writer) findOptionalNoun(domain, noun string) (ret nounInfo, err error) {
	if len(noun) > 0 {
		if e := m.db.QueryRow(`
	select mn.domain, mn.rowid, ',' || mk.rowid || ',' || mk.path
	from mdl_noun mn
	join mdl_kind mk 
		on (mn.kind = mk.rowid)
	join domain_tree dt
		on (dt.uses = mn.domain)
	where base = ?1
	and noun = ?2
	limit 1`, domain, noun).Scan(&ret.domain, &ret.id, &ret.fullpath); e != nil && e != sql.ErrNoRows {
			err = e
		} else {
			ret.name = noun
		}
	}
	return
}
