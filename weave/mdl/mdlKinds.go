package mdl

import (
	"database/sql"
	"strconv"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang"
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
func updatePath(res sql.Result, parent string, path *string) (err error) {
	if i, e := res.LastInsertId(); e != nil {
		err = e
	} else {
		part := "," + strconv.FormatInt(i, 10)
		if len(parent) > 0 {
			part += parent
		} else {
			part += ","
		}
		*path = part
	}
	return
}

type KindInfo struct {
	id        int64  // unique id of the kind
	name      string // validated name of the kind
	domain    string // validated domain name
	path      string // comma separated ids of ancestors: ,2,1,
	exact     bool
	_fullpath string
}

func (ki *KindInfo) numAncestors() int {
	// ,   = no ancestors
	// ,2, = 1 ancestor
	return strings.Count(ki.path, ",") - 1
}

// path starting with the kind's own id. ",id,...,"
func (ki *KindInfo) fullpath() string {
	if ki.id > 0 && len(ki._fullpath) == 0 {
		ki._fullpath = "," + strconv.FormatInt(ki.id, 10) + ki.path
	}
	return ki._fullpath
}

// if specified, must exist.
func (m *Pen) findOptionalKind(kind string) (ret KindInfo, err error) {
	if len(kind) > 0 {
		ret, err = m.findRequiredKind(kind)
	}
	return
}

// if not specified errors, also errors if not found.
func (m *Pen) findRequiredKind(kind string) (ret KindInfo, err error) {
	if out, e := m.findKind(kind); e != nil {
		err = e
	} else if out.id == 0 {
		err = errutil.Fmt("%w kind %q in domain %q", Missing, kind, m.domain)
	} else {
		ret = out
	}
	return
}

// if not specified errors, makes no assumptions about the results
func (m *Pen) findKind(kind string) (ret KindInfo, err error) {
	if len(kind) == 0 {
		err = errutil.New("empty name for kind")
	} else if singular, e := m.singularize(kind); e != nil {
		err = e
	} else {
		var rank int
		e := m.db.QueryRow(`
	select domain, 
		mk.rowid, 
		mk.kind,
		',' || mk.path, 
		case when ?2 = kind then 1 
		     when ?3 = kind then 2 
		     when ?2 = singular then 3
		     when ?3 = singular then 4 
		else 0 
		end as rank
	from mdl_kind mk
	join domain_tree
		on (uses = domain)
	where base = ?1
	and rank > 0
	order by rank
	limit 1`, // order by rank means the lowest number is first
			m.domain, kind, singular).Scan(
			&ret.domain, &ret.id, &ret.name, &ret.path, &rank)
		switch e {
		case nil:
			ret.exact = rank == 1
		case sql.ErrNoRows:
			// nothing found? still set the name for easier logging;
			// the empty id can disambiguate success from not found
			ret.name = kind
		default:
			err = e
		}
	}
	return
}

func (m *Pen) singularize(kind string) (ret string, err error) {
	if len(kind) < 2 {
		ret = kind //
	} else if e := m.db.QueryRow(`
	select one
	from mdl_plural
	join domain_tree
		on (uses = domain)
	where base = ?1 and many = ?2
	limit 1`, m.domain, kind).Scan(&ret); e == sql.ErrNoRows {
		ret = lang.Singularize(kind)
	} else {
		err = e // other error or nil.
	}
	return
}
