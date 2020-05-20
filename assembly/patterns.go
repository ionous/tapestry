package assembly

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

func AssemblePatterns(m *Modeler, db *sql.DB) (err error) {
	if e := checkPatternSetup(db); e != nil {
		err = e
	} else if e := checkRuleSetup(db); e != nil {
		err = e
	} else if e := copyRules(db); e != nil {
		err = e
	} else if e := copyPatterns(db); e != nil {
		err = e
	}
	return
}

// fix: this probably needs work to get parameter ordering sensible
func copyPatterns(db *sql.DB) error {
	// problem: assumes
	_, e := db.Exec(
		`insert into mdl_pat 
		select pattern, param, type, 
		(case param  
		when pattern then 0
		else (select 1+count() 
				from asm_pattern_decl cp 
				where ap.pattern= cp.pattern 
				and cp.pattern != cp.param 
				and ap.ogid > cp.ogid)
		end) idx
		from asm_pattern_decl ap
		order by pattern, idx`)
	return e
}

func checkPatternSetup(db *sql.DB) (err error) {
	var now, last patternInfo
	var declaredReturn string

	// find where variable names and pattern names conflict
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
		&now.pat, &now.arg, &now.typ, &now.decl); e != nil {
		err = e
	} else {
		// search for other conflicts
		if e := tables.QueryAll(db,
			`select distinct pattern, param, type, decl from asm_pattern
			order by pattern, param, type, decl desc`,
			func() error {
				e := now.compare(&last, &declaredReturn)
				last = now
				return e
			},
			&now.pat, &now.arg, &now.typ, &now.decl); e != nil {
			err = e
		} else if e := last.flush(&declaredReturn); e != nil {
			err = e
		}
	}
	return
}

type patternInfo struct {
	pat, arg, typ string
	decl          bool
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
		if change := (now.pat != was.pat || now.arg != was.arg); change && !now.decl {
			// decl(s) come first, so if there's a change... it should only happen with a decl.
			err = errutil.Fmt("Pattern %q's %q missing declaration", now.pat, now.arg)
		} else if !change && (now.typ != was.typ) {
			// regardless -- types should be consistent.
			err = errutil.New("Pattern %q's %q type conflict, was %q now %q", now.pat, now.arg, was.typ, now.typ)
		} else if now.decl && now.pat == now.arg {
			// assuming everything's ok, a decl where pat and arg match means the type of the pattern itself.
			*pret = now.typ
		}
	}
	return
}