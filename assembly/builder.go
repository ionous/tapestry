package assembly

import (
	"database/sql"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

// map name to pattern interface
type patternEntry struct {
	patternName             string // name of the pattern
	patternType             string // "return" type of the pattern
	params, locals, returns []fieldInit
}

type fieldInit struct {
	Name     string
	Affinity affine.Affinity
	Type     string // ex. record name, "aspect", "trait", "float64", ...
	Init     rt.Assignment
}

func (fi *fieldInit) Field() g.Field {
	return g.Field{fi.Name, fi.Affinity, fi.Type}
}

func (pat *patternEntry) AddField(cat string, fi fieldInit) (err error) {
	switch cat {
	case tables.NAMED_PARAMETER:
		pat.params = append(pat.params, fi)
	case tables.NAMED_LOCAL:
		pat.locals = append(pat.locals, fi)
	case tables.NAMED_RETURN:
		pat.returns = append(pat.returns, fi)
	default:
		err = errutil.New("unknown category", cat)
	}
	return
}

type patternCache map[string]*patternEntry

// fix: report errors
func (p *patternEntry) newFragment() *PatternFrag {
	pat := PatternFrag{Name: p.patternName}
	//
	if ps := p.params; len(ps) > 0 {
		for _, fi := range ps {
			pat.Fields = append(pat.Fields, fi.Field())
			// eventually labels might be different than parameter names
			//( cause swift makes that seem cool )
			pat.Labels = append(pat.Labels, fi.Name)
		}
	}
	if ps := p.locals; len(ps) > 0 {
		for _, fi := range ps {
			pat.Fields = append(pat.Fields, fi.Field())
			pat.Locals = append(pat.Locals, fi.Init)
		}
	}

	// fix: report if too many returns
	if ps := p.returns; len(ps) > 0 {
		var found bool
		res := ps[0].Field()
		for _, f := range pat.Fields {
			if f.Name == res.Name {
				found = true
				break
			}
		}
		if !found {
			pat.Fields = append(pat.Fields, res)
		}
		pat.Return = res.Name
	}
	//
	return &pat
}

// read pattern declarations from the ephemera db
func buildPatternCache(db *sql.DB) (ret patternCache, err error) {
	// build the pattern cache
	out := make(patternCache)
	var inPat, inParam, inCat, inType string
	var inKind, inAff sql.NullString
	var inProg []byte
	var last *patternEntry
	if e := tables.QueryAll(db,
		// fix: these are grouped by pattern, param, cat --
		// so there are conflicts in names and types we wont see
		// this needs much better handling of conflicting and redundant info
		`select ap.pattern, ap.param, ap.cat, ap.type, ap.affinity, ap.kind, ep.prog
		from asm_pattern_decl ap
		left join eph_prog ep
		on (ep.rowid = ap.idProg)`,
		func() (err error) {
			// fix: need to handle conflicting prog definitions
			// fix: should watch for locals which shadow parameter names ( i think, ideally merge them )
			if last == nil || last.patternName != inPat {
				if inPat != inParam {
					err = errutil.New("expected the first param should be the pattern return type", inPat, inProg, inType)
				} else {
					last = &patternEntry{patternName: inPat, patternType: inType}
					out[inPat] = last
				}
			}
			if err == nil && inParam != inPat {
				// fix: these should probably be tables.PRIM_ names
				// ie. "text" not "text_eval" -- tests and other things have to be adjusted
				// it also seems a bad time to be camelizing things.
				paramName := lang.Breakcase(inParam)
				if aff, typeName := convertType(inType, inKind.String, inAff.String); len(aff) == 0 {
					err = errutil.New("unknown type", inType, inKind, inAff)
				} else if i, e := decodeProg(inProg, aff); e != nil {
					err = errutil.New("couldnt decode", inPat, paramName, e)
				} else {
					err = last.AddField(inCat, fieldInit{paramName, aff, typeName, i})
				}
			}
			return
		},
		&inPat, &inParam, &inCat, &inType, &inAff, &inKind, &inProg); e != nil {
		err = e
	} else {
		ret = out
	}
	return
}

func convertType(inType, inKind, inAff string) (retAff affine.Affinity, retType string) {
	// locals have simple type names, parameters are still using _eval.
	switch inType {
	case "text_eval", "text":
		retAff = affine.Text
	case "number_eval", "number":
		retAff = affine.Number
	case "bool_eval", "bool":
		retAff = affine.Bool
	case "text_list", "text_list_eval":
		retAff = affine.TextList
	case "num_list", "num_list_eval":
		retAff = affine.NumList
	default:
		// the type might be some sort of kind...
		if len(inKind) > 0 {
			switch aff := affine.Affinity(inAff); aff {
			case affine.Object:
				retAff, retType = affine.Text, "object="+inKind
			case affine.Record, affine.RecordList:
				retAff, retType = aff, inKind
			}
		}
	}
	return
}

// the author specified a "local init"
// it has a Value assignment, we want to dig out that assignment and assign it to the term prep.
// im not convinced that terms and assigments should be different beasts...
// if they were the same thing.. this would look different.
func decodeProg(prog []byte, aff affine.Affinity) (ret rt.Assignment, err error) {
	if haveProg := len(prog) > 0; haveProg {
		var local story.LocalInit
		if e := tables.DecodeGob(prog, &local); e != nil {
			err = e
		} else if a := local.Value.Affinity(); len(a) > 0 && a != aff {
			// note: some expressions (ex. GetAtField) cant determine affinity until runtime
			err = errutil.New("incompatible arguments, wanted", aff, "have expression of", a)
		} else {
			ret = local.Value
		}
	}
	return
}

func (ps patternCache) WriteFragments(asm *Assembler, kind string) (err error) {
	if e := asm.WriteAncestor(kind, ""); e != nil {
		err = e
	} else {
		for _, p := range ps {
			if p.patternType == kind {
				if e := WriteFragment(asm, kind, p.newFragment()); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
	}
	return
}

//
func WriteFragment(asm *Assembler, kind string, pat *PatternFrag) (err error) {
	if e := asm.WriteAncestor(pat.Name, kind); e != nil {
		err = e
	} else if e := asm.WritePattern(pat.Name, pat.Return, pat.Labels); e != nil {
		err = e
	} else {
		localOfs := len(pat.Labels)
		// write fields exactly as if the pattern was a kind
		for _, prop := range pat.Fields {
			if e := asm.WriteField(pat.Name, prop.Name, prop.Type, prop.Affinity.String()); e != nil {
				err = errutil.Append(err, e)
			}
		}
		// write initialization to mdl_start
		for i, val := range pat.Locals {
			if val != nil { // not every local has an initial value
				f := pat.Fields[localOfs+i]
				// fix: write encoded value... which we just recently decoded... :/
				meta := struct{ Init rt.Assignment }{val}
				if prog, e := tables.EncodeGob(&meta); e != nil {
					err = errutil.Append(err, e)
				} else if e := asm.WriteStart(pat.Name, f.Name, prog); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
	}
	return
}
