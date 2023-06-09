package mdl

import (
	"database/sql"
	"strconv"

	"github.com/ionous/errutil"
)

// fullpaths start and end with commas;
// for backwards compat this strips the leading comma;
// fix? the trailing comma is also redundant but the db has it for every non-null entry
func trimPath(fullpath string) (ret string) {
	if len(fullpath) > 0 {
		ret = fullpath[1:]
	}
	return
}

// fix? it would probably be better to have a separate table of: domain, aspect, trait
// currently, the runtime expects that aspects are a kind, and its traits are fields.
func updatePath(res sql.Result, path *string) (err error) {
	if i, e := res.LastInsertId(); e != nil {
		err = e
	} else {
		*path = "," + strconv.FormatInt(i, 10) + ","
	}
	return
}

type kindInfo struct {
	id        int64  // unique id of the kind
	name      string // validated name of the kind
	domain    string // validated domain name
	path      string // comma separated ids of ancestors: ,2,1,
	_fullpath string
}

// path starting with the kind's own id. ",id,...,"
func (ki *kindInfo) fullpath() string {
	if ki.id > 0 && len(ki._fullpath) == 0 {
		ki._fullpath = "," + strconv.FormatInt(ki.id, 10) + ki.path
	}
	return ki._fullpath
}

// if specified, must exist.
func (m *Writer) findOptionalKind(domain, kind string) (ret kindInfo, err error) {
	if len(kind) > 0 {
		ret, err = m.findRequiredKind(domain, kind)
	}
	return
}

// if not specified errors, also errors if not found.
func (m *Writer) findRequiredKind(domain, kind string) (ret kindInfo, err error) {
	if out, e := m.findKind(domain, kind); e != nil {
		err = e
	} else if out.id == 0 {
		err = errutil.Fmt("%w kind %q in domain %q", Missing, kind, domain)
	} else {
		ret = out
	}
	return
}

// if not specified errors, makes no assumptions about the results
func (m *Writer) findKind(domain, kind string) (ret kindInfo, err error) {
	if len(kind) == 0 {
		err = errutil.New("empty name for kind")
	} else if e := m.db.QueryRow(`
	select domain, mk.rowid, ',' || mk.path
	from mdl_kind mk
	join domain_tree
		on (uses = domain)
	where base = ?1
	and kind = ?2
	limit 1`, domain, kind).Scan(&ret.domain, &ret.id, &ret.path); e != nil && e != sql.ErrNoRows {
		err = e
	} else {
		ret.name = kind
	}
	return
}
