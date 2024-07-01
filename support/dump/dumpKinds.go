package dump

import (
	"database/sql"
	"slices"

	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/tables"
)

func QueryKinds(db *sql.DB, dec qdb.CommandDecoder, scene string) (ret []rt.Kind, err error) {
	if ks, e := tables.QueryStrings(db, must("kinds"), scene); e != nil {
		err = e
	} else {
		kb := kindBuilder{db, dec, make(fieldCache)}
		ret = make([]rt.Kind, len(ks))
		for i, k := range ks {
			if rk, e := qdb.BuildKind(kb, k); e != nil {
				err = e
				break
			} else {
				ret[i] = rk
			}
		}
	}
	return
}

type kindBuilder struct {
	db    *sql.DB
	dec   qdb.CommandDecoder
	cache fieldCache
}

type fieldCache map[string][]rt.Field

func (kb kindBuilder) KindOfAncestors(k string) (ret []string, err error) {
	return tables.QueryStrings(kb.db, must("ancestors"), k)
}

func (kb kindBuilder) FieldsOf(k string) (ret []rt.Field, err error) {
	if fs, ok := kb.cache[k]; ok {
		ret = fs
	} else if rows, e := kb.db.Query(must("fields"), k); e != nil {
		err = e
	} else if fs, e := qdb.ScanFields(rows, kb.dec); e != nil {
		err = e
	} else {
		kb.cache[k] = fs
		ret = fs
	}
	return
}

func findField(k *rt.Kind, name string, last int) int {
	return slices.IndexFunc(k.Fields[last:], func(n rt.Field) bool {
		return n.Name == name
	})
}
