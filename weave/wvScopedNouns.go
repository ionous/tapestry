package weave

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

// find the noun with the closest name in this scope
// skips aliases for the sake of backwards compatibility:
// there should be a difference between "a noun is known as"
// and "understand this word by the player as" -- and currently there's not.
func (d *Domain) GetClosestNoun(name string) (ret struct{ name, domain string }, err error) {
	if e := d.catalog.db.QueryRow(`
	select mn.noun, mn.domain  
	from mdl_name my 
	join mdl_noun mn
		on (mn.rowid = my.noun)
	join domain_tree dt
		on (dt.uses = my.domain)
	where base = ?1
	and my.name = ?2
	and my.rank >= 0
	order by my.rank, my.rowid asc
	limit 1`, d.name, name).Scan(&ret.name, &ret.domain); e == sql.ErrNoRows {
		err = errutil.Fmt("%w couldn't find a noun named %s", mdl.Missing, name)
	} else if e != nil {
		err = errutil.New("database error", e)
	}
	return
}
