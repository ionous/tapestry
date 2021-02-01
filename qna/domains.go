package qna

import (
	"database/sql"

	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

func ActivateDomain(db *sql.DB, domain string, active bool) (err error) {
	if _, e := db.Exec(run_domain, domain, active); e != nil {
		err = errutil.New("ActivateDomain", domain, e)
	}
	return
}

// inserts a newly active domain name, or sets an existing domain's status
var run_domain = tables.InsertWith("run_domain",
	"on conflict(domain) do update set active=excluded.active;",
	"domain", "active")
