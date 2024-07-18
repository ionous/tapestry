package mdl

import (
	"database/sql"

	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
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

func (pen *Pen) WriteSceneStart(rank int, exe []rt.Execute) (err error) {
	// fix: current domain changed looks for the pattern "... begins"
	eventName := inflect.Normalize(pen.domain + " begins")
	pb := NewPatternBuilder(eventName)
	pb.AppendRule(0, rt.Rule{
		Name: "scene " + pen.domain, // arbitrary name
		Exe:  exe,
	})
	return pen.AddPattern(pb.Pattern)
}
