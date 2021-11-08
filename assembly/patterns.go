package assembly

import (
	"database/sql"
	"strings"

	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

func AssemblePatterns(asm *Assembler) (err error) {
	db := asm.cache.DB()
	if e := checkPatternSetup(db); e != nil {
		err = e
	} else if cache, e := buildPatternCache(asm.cache.DB()); e != nil {
		err = errutil.New(e, "reading patterns")
	} else if e := buildPatternRules(asm, cache); e != nil {
		err = errutil.New(e, "building rules")
	} else if e := cache.WriteFragments(asm, "patterns"); e != nil {
		err = errutil.New(e, "building patterns")
	} else if e := cache.WriteFragments(asm, "actions"); e != nil {
		err = errutil.New(e, "building actions")
	} else if e := cache.WriteFragments(asm, "events"); e != nil {
		err = errutil.New(e, "building events")
	}
	return
}

func checkPatternSetup(db *sql.DB) (err error) {
	var now, last patternInfo
	var declaredReturn string
	// find where variable names and pattern names conflict
	// FIX: this had been pretending to query multiple parameters...
	// but it only actually looks at the pattern name....
	if e := tables.QueryAll(db,
		`select distinct pn.name
		from eph_pattern ep
		left join eph_named pn
			on (ep.idNamedPattern = pn.rowid)
		left join eph_named kn
			on (ep.idNamedParam = kn.rowid)
		where ep.idNamedPattern != ep.idNamedParam
		and pn.name = kn.name`,
		func() error {
			e := now.compare(&last, &declaredReturn)
			last = now
			return e
		},
		&now.pat /*, &now.arg, &now.typ, &now.decl */); e != nil {
		err = e
	} /* else {
		// -- disabled for now,
		// there are various issues including
		// 	. inconsistent types ( _eval vs. prim )
		//  . name casing
		//  . mixing of params and locals ( which will hurt $1 parameter indexing )

		// search for other conflicts
		// note: these arent complete b/c we'd need to identify the types of vars and fields and carry those out
		// ( similar to NewPatternRef )
		if e := tables.QueryAll(db,
			`select distinct pattern, param, type, affinity, decl from asm_pattern
			order by pattern, param, type, affinity, decl desc`,
			func() error {
				e := now.compare(&last, &declaredReturn)
				last = now
				return e
			},
			&now.pat, &now.arg, &now.typ, &now.aff, &now.decl); e != nil {
			err = e
		} else if e := last.flush(&declaredReturn); e != nil {
			err = e
		}
	} */
	return
}

type patternInfo struct {
	pat, arg, typ, aff string
	decl               bool
}

func (now *patternInfo) flush(pret *string) (err error) {
	if len(now.pat) > 0 && len(*pret) == 0 {
		err = errutil.Fmt("Pattern %q never declared a return type", now.pat)
	}
	(*pret) = ""
	return
}

func (now *patternInfo) compare(was *patternInfo, pret *string) (err error) {
	// new pattern detected....
	if now.pat != was.pat {
		if e := was.flush(pret); e != nil {
			err = e
		}
	}
	if err == nil {
		// fix: hack off eval names
		// ex. pattern	param	type	affinity	decl
		// objectGroupingTest	objects	text_list		1
		// objectGroupingTest	objects	text_list_eval		0
		now.typ, was.typ = strings.TrimSuffix(now.typ, "_eval"), strings.TrimSuffix(was.typ, "_eval")
		if change := (now.pat != was.pat || now.arg != was.arg); change && !now.decl {
			// decl(s) come first, so if there's a change... it should only happen with a decl.
			err = errutil.Fmt("Pattern %q's %q missing declaration", now.pat, now.arg)
		} else if !change && (now.typ != was.typ) {
			// regardless -- types should be consistent.
			err = errutil.Fmt("Pattern %q's %q type conflict, was %q now %q", now.pat, now.arg, was.typ, now.typ)
		} else if !change && (now.aff != was.aff) && (len(now.aff) > 0 || len(was.aff) > 0) {
			// regardless -- types should be consistent.
			err = errutil.Fmt("Pattern %q's %q affinities conflict, was %q now %q", now.pat, now.arg, was.aff, now.aff)
		} else if now.decl && now.pat == now.arg {
			// assuming everything's ok, a decl where pat and arg match means the type of the pattern itself.
			*pret = now.typ
		}
	}
	return
}
