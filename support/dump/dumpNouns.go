package dump

import (
	"database/sql"
	"fmt"

	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/qna/raw"
	"git.sr.ht/~ionous/tapestry/tables"
)

func QueryNouns(db *sql.DB, kd qdb.KindDecoder, scene string) (ret []raw.NounData, err error) {
	if ns, e := QueryInnerNouns(db, scene); e != nil {
		err = fmt.Errorf("%w while querying ids", e)
	} else if e := QueryAliases(db, ns); e != nil {
		err = fmt.Errorf("%w while querying aliases", e)
	} else if e := QueryValues(db, kd, ns); e != nil {
		err = fmt.Errorf("%w while querying values", e)
	} else {
		ret = ns
	}
	return
}

func QueryInnerNouns(db *sql.DB, scene string) (ret []raw.NounData, err error) {
	var n raw.NounData
	if rows, e := db.Query(must("nouns"), scene); e != nil {
		err = e
	} else {
		err = tables.ScanAll(rows, func() (err error) {
			ret = append(ret, n)
			return
		}, &n.Id, &n.Domain, &n.Noun, &n.Kind)
	}
	return
}

func QueryValues(db *sql.DB, kd qdb.KindDecoder, ns []raw.NounData) (err error) {
	vals, recs := must("values"), must("records")
	for i, n := range ns {
		if k, e := kd.GetKindByName(n.Kind); e != nil {
			err = e
			break
		} else {
			read := makeValueReader(&ns[i], k, kd)
			if rows, e := db.Query(vals, n.Id); e != nil {
				err = e
				break
			} else if e := read.values(rows); e != nil {
				err = e
				break
			} else if rows, e := db.Query(recs, n.Id); e != nil {
				err = e
				break
			} else if e := read.records(rows); e != nil {
				err = e
				break
			}
		}
	}
	return
}

func QueryAliases(db *sql.DB, ns []raw.NounData) (err error) {
	q := must("aliases")
	for i, n := range ns {
		if as, e := tables.QueryStrings(db, q, n.Id); e != nil {
			err = e
		} else {
			if cn := as[0]; cn != n.Noun {
				ns[i].CommonName = cn
			}
			ns[i].Aliases = as[1:]
		}
	}
	return
}
