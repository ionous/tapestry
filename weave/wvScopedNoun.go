package weave

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang"
	"github.com/ionous/errutil"
)

type ScopedNoun struct {
	name        string
	domain      *Domain
	localRecord localRecord // cache for a literal.RecordValue
}

func (n *ScopedNoun) Domain() (ret string) {
	return n.domain.name
}

func (n *ScopedNoun) Name() (ret string) {
	return n.name
}

func (n *ScopedNoun) FindKind() (ret string, err error) {
	d := n.domain
	err = d.cat.db.QueryRow(`
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

func (n *ScopedNoun) AddAliases(at string, aliases []string) (err error) {
	pen := n.domain.cat.Pin(n.domain.name, at)
	for _, a := range aliases {
		a := lang.Normalize(a)
		if e := pen.AddName(n.name, a, -1); e != nil {
			err = e
			break
		}
	}
	return
}

func (n *ScopedNoun) WriteValue(at, field string, path []string, value literal.LiteralValue) (err error) {
	if value == nil {
		err = errutil.New("null value", n.name, field)
	} else if rv, e := n.recordValues(at); e != nil {
		err = e
	} else {
		err = rv.writeValue(n.name, at, field, path, value)
	}
	return
}

//
func (n *ScopedNoun) recordValues(at string) (ret localRecord, err error) {
	if n.localRecord.isValid() {
		ret = n.localRecord
	} else if kind, e := n.FindKind(); e != nil {
		err = e
	} else {
		k := kindCat{domain: n.domain, kind: kind}
		rv := localRecord{k, new(literal.RecordValue), at}
		ret, n.localRecord = rv, rv
	}
	return
}
