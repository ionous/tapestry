package mdl

import (
	"database/sql"

	"github.com/ionous/errutil"
)

// findDomain validates that the named domain exists
// the returned name is the same as the passed name.
func (m *Writer) findDomain(domain string) (ret string, err error) {
	if e := m.db.QueryRow(`
	select domain 
	from mdl_domain 
	where domain = ?1`, domain).Scan(&ret); e == sql.ErrNoRows {
		err = errutil.Fmt("domain not found %q", domain)
	} else {
		err = e
	}
	return
}
