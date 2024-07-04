package mdl

import (
	"database/sql"

	"fmt"
)

// findScene validates that the named domain exists
// the returned name is the same as the passed name.
func (pen *Pen) findScene() (ret string, err error) {
	domain := pen.domain
	if e := pen.db.QueryRow(`
	select domain 
	from mdl_domain 
	where domain = ?1`, domain).Scan(&ret); e == sql.ErrNoRows {
		err = fmt.Errorf("scene not found %q", domain)
	} else {
		err = e
	}
	return
}
