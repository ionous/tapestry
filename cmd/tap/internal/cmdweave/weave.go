package cmdweave

import (
	"database/sql"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"log"

	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

// write the catalog to the passed database
func Weave(cat *eph.Catalog, db *sql.DB) (err error) {
	var queue writeCache
	if e := BuildCatalog(cat, &queue); e != nil {
		err = e
	} else {
		log.Println("writing", len(queue), "entries")
		if e := tables.CreateModel(db); e != nil {
			err = errutil.New("couldnt create model", e)
		} else if tx, e := db.Begin(); e != nil {
			err = errutil.New("couldnt create transaction", e)
		} else {
			w := mdl.Writer(func(q string, args ...interface{}) (err error) {
				if _, e := tx.Exec(q, args...); e != nil {
					err = e
				}
				return
			})
			for _, q := range queue {
				if e := w.Write(q.tgt, q.args...); e != nil {
					tx.Rollback()
					err = errutil.New("couldnt write to", q.tgt, e)
					break
				}
			}
			if err == nil {
				if e := tx.Commit(); e != nil {
					err = errutil.New("couldnt commit", e)
				}
			}
		}
	}
	return
}

func BuildCatalog(cat *eph.Catalog, w eph.Writer) (err error) {
	// go process all of the ephemera
	if e := cat.AssembleCatalog(eph.PhaseActions{
		eph.AncestryPhase: eph.AncestryActions,
		eph.FieldPhase:    eph.FieldActions,
		eph.NounPhase:     eph.NounActions,
	}); e != nil {
		err = e
	} else if e := cat.WritePlurals(w); e != nil {
		err = e
	} else if e := cat.WriteOpposites(w); e != nil {
		err = e
	} else if e := cat.WriteDomains(w); e != nil {
		err = e
	} else if e := cat.WriteKinds(w); e != nil {
		err = e
	} else if e := cat.WriteFields(w); e != nil {
		err = e
	} else if e := cat.WriteNouns(w); e != nil {
		err = e
	} else if e := cat.WriteNames(w); e != nil {
		err = e
	} else if e := cat.WritePatterns(w); e != nil {
		err = e
	} else if e := cat.WriteLocals(w); e != nil {
		err = e
	} else if e := cat.WriteDirectives(w); e != nil {
		err = e
	} else if e := cat.WriteRelations(w); e != nil {
		err = e
	} else if e := cat.WritePairs(w); e != nil {
		err = e
	} else if e := cat.WriteRules(w); e != nil {
		err = e
	} else if e := cat.WriteValues(w); e != nil {
		err = e
	} else if e := cat.WriteChecks(w); e != nil {
		err = e
	}
	return
}

// a terrible way to optimize database writes
type cachedWrite struct {
	tgt  string
	args []interface{}
}
type writeCache []cachedWrite

func (q *writeCache) Write(tgt string, args ...interface{}) (err error) {
	(*q) = append(*q, cachedWrite{tgt, args})
	return
}
