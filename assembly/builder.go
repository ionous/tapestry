package assembly

import (
	"bytes"
	"database/sql"
	"encoding/gob"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/dl/term"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

type BuildRule struct {
	Query        string
	NewContainer func(name string) interface{}
	NewEl        func(c interface{}) interface{}
}

// map name to pattern interface
type patternEntry struct {
	patternName string          // name of the pattern
	patternType string          // "return" type of the pattern
	prologue    []term.Preparer // list of all parameters sent to the pattern
	locals      []term.Preparer // ...
	returns     term.Preparer
}

func (pat *patternEntry) AddParam(cat string, param term.Preparer) (err error) {
	switch cat {
	case tables.NAMED_PARAMETER:
		pat.prologue = append(pat.prologue, param)
	case tables.NAMED_LOCAL:
		pat.locals = append(pat.locals, param)
	case tables.NAMED_RETURN:
		// fix? check for multiple / different sets?
		pat.returns = param
	default:
		err = errutil.New("unknown category", cat)
	}
	return
}

type patternCache map[string]*patternEntry

func (cache patternCache) init(name, patternType string) (ret *pattern.Pattern, okay bool) {
	if c, ok := (cache)[name]; ok && c.patternType == patternType {
		ret = &pattern.Pattern{
			Name:    name,
			Params:  c.prologue,
			Locals:  c.locals,
			Returns: c.returns,
		}
		okay = true
	}
	return
}

// read pattern declarations from the ephemera db
func buildPatternCache(db *sql.DB) (ret patternCache, err error) {
	// build the pattern cache
	out := make(patternCache)
	var patternName, paramName, category, typeName string
	var kind, affinity sql.NullString
	var prog []byte
	var last *patternEntry
	if e := tables.QueryAll(db,
		`select ap.pattern, ap.param, ap.cat, ap.type, ap.affinity, ap.kind, ep.prog
		from asm_pattern_decl ap
		left join eph_prog ep
		on (ep.rowid = ap.idProg)`,
		func() (err error) {
			// fix: need to handle conflicting prog definitions
			// fix: should watch for locals which shadow parameter names ( i think, ideally merge them )
			if last == nil || last.patternName != patternName {
				if patternName != paramName {
					err = errutil.New("expected the first param should be the pattern return type", patternName, paramName, typeName)
				} else {
					last = &patternEntry{patternName: patternName, patternType: typeName}
					out[patternName] = last
				}
			}
			if err == nil && paramName != patternName {
				// fix: these should probably be tables.PRIM_ names
				// ie. "text" not "text_eval" -- tests and other things have to be adjusted
				// it also seems a bad time to be camelizing things.
				paramName := lang.Breakcase(paramName)

				// locals have simple type names, parameters are still using _eval.
				var p term.Preparer
				var dig digInit
				switch typeName {
				case "text_eval", "text":
					prep := &term.Text{Name: paramName}
					p, dig = prep, func(in core.Assignment) (okay bool) {
						if src, ok := in.(*core.FromText); ok {
							prep.Init, okay = src.Val, true
						}
						return
					}
				case "number_eval", "number":
					prep := &term.Number{Name: paramName}
					p, dig = prep, func(in core.Assignment) (okay bool) {
						if src, ok := in.(*core.FromNum); ok {
							prep.Init, okay = src.Val, true
						}
						return
					}
				case "bool_eval", "bool":
					prep := &term.Bool{Name: paramName}
					p, dig = prep, func(in core.Assignment) (okay bool) {
						if src, ok := in.(*core.FromBool); ok {
							prep.Init, okay = src.Val, true
						}
						return
					}
				case "text_list", "text_list_eval":
					prep := &term.TextList{Name: paramName}
					p, dig = prep, func(in core.Assignment) (okay bool) {
						if src, ok := in.(*core.FromTexts); ok {
							prep.Init, okay = src.Vals, true
						}
						return
					}
				case "num_list", "num_list_eval":
					prep := &term.NumList{Name: paramName}
					p, dig = prep, func(in core.Assignment) (okay bool) {
						if src, ok := in.(*core.FromNumbers); ok {
							prep.Init, okay = src.Vals, true
						}
						return
					}
				default:
					// the type might be some sort of kind...
					if kind := kind.String; len(kind) > 0 {
						switch aff := affinity.String; aff {
						case string(affine.Object):
							prep := &term.Object{Name: paramName, Kind: kind}
							p, dig = prep, func(in core.Assignment) (okay bool) {
								if src, ok := in.(*core.FromText); ok {
									prep.Init, okay = src.Val, true
								}
								return
							}
						case string(affine.Record):
							prep := &term.Record{Name: paramName, Kind: kind}
							p, dig = prep, func(in core.Assignment) (okay bool) {
								if src, ok := in.(*core.FromRecord); ok {
									prep.Init, okay = src.Val, true
								}
								return
							}
						case string(affine.RecordList):
							prep := &term.RecordList{Name: paramName, Kind: kind}
							p, dig = prep, func(in core.Assignment) (okay bool) {
								if src, ok := in.(*core.FromRecords); ok {
									prep.Init, okay = src.Vals, true
								}
								return
							}
						}
					}
				}
				if e := decode(prog, dig); e != nil {
					err = errutil.New("couldnt decode", patternName, paramName, e)
				} else if p != nil {
					err = last.AddParam(category, p)
				} else {
					err = errutil.Fmt("pattern %q parameter %q has unknown type %q(%s)",
						patternName, paramName, typeName, affinity.String)
				}
			}
			return
		},
		&patternName, &paramName, &category, &typeName, &affinity, &kind, &prog); e != nil {
		err = e
	} else {
		ret = out
	}
	return
}

type digInit func(in core.Assignment) (okay bool)

// the author specified a "local init"
// it has a Value assignment, we want to dig out that assignment and assign it to the term prep.
// im not convinced that terms and assigments should be different beasts...
// if they were the same thing.. this would look different.
func decode(prog []byte, dig digInit) (err error) {
	if haveProg := len(prog) > 0; haveProg {
		var local story.LocalInit
		dec := gob.NewDecoder(bytes.NewBuffer(prog))
		if e := dec.Decode(&local); e != nil {
			err = e
		} else if !dig(local.Value) {
			err = errutil.New("couldnt convert from %T", local.Value)
		}
	}
	return
}

// collect the rules of all the various patterns and write them into the assembly
func buildPatterns(asm *Assembler) (err error) {
	if patterns, e := buildPatternCache(asm.cache.DB()); e != nil {
		err = e
	} else {
		err = buildPatternRules(asm, patterns)
	}
	return
}

func buildPatternRules(asm *Assembler, patterns patternCache) error {
	var name string
	var prog []byte
	rule := BuildRule{
		Query: `select pattern, prog from asm_rule where type='rule'`,
		NewContainer: func(name string) (ret interface{}) {
			if c, ok := patterns.init(name, "execute"); ok {
				ret = c
			}
			return
		},
		NewEl: func(c interface{}) interface{} {
			pat := c.(*pattern.Pattern)
			pat.Rules = append(pat.Rules, &pattern.Rule{})
			return pat.Rules[len(pat.Rules)-1]
		},
	}
	return rule.buildFromRule(asm, &name, &prog)
}
