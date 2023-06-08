package weave

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
)

type ScopedNoun struct {
	// Requires    // kinds ( when resolved, can have one direct parent )
	name        string
	domain      *Domain
	localRecord localRecord // store the values of the noun as a record.
}

func (n *ScopedNoun) Kind() (ret string, err error) {
	d := n.domain
	err = d.catalog.db.QueryRow(`
		select mk.kind
	from mdl_kind mk 
	join mdl_noun mn
		on (mn.kind = mk.rowid)
	join domain_tree dt
		on (dt.uses = mn.domain)
	where base = ?1
	and noun = ?2
	limit 1`, d.name, n.name).Scan(&ret)
	return
}

// stores literal values because they are serializable
// ( as opposed to generic values which aren't. )
func (n *ScopedNoun) recordValues(at string) (ret localRecord, err error) {
	if n.localRecord.isValid() {
		ret = n.localRecord
	} else if kind, e := n.Kind(); e != nil {
		err = e
	} else {
		k := kindCat{domain: n.domain, kind: kind}
		rv := localRecord{k, new(literal.RecordValue), at}
		ret, n.localRecord = rv, rv
	}
	return
}
