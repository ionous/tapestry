package assembly

import (
	"strings"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/grammar"
	"git.sr.ht/~ionous/iffy/jsn/cin"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

func AssembleGrammar(asm *Assembler) (err error) {
	var directives []*grammar.Directive
	var aliases []*grammar.Alias
	var names []string
	var prog []byte
	if e := tables.QueryAll(asm.cache.DB(),
		`select prog 
		from eph_prog where progType='grammar'
		order by rowid`,
		func() (err error) {
			var decl grammar.Grammar
			var name string // try to make a name for ourselves
			if e := cin.Decode(&decl, prog, iffy.AllSignatures); e != nil {
				err = e
			} else {
				switch d := decl.Grammar.(type) {
				case *grammar.Alias:
					aliases = append(aliases, d)

				case *grammar.Directive:
					// for now, just let runtime make the scanners..
					// because, if we store the scanner, then we have to register the scanners.
					// which means we have both the parser and the maker operations registered
					// while using only one of them
					name = strings.Join(d.Lede, "/")
					names = append(names, name)
					directives = append(directives, d)
				default:
					err = errutil.Fmt("unhandled grammar %T", d)
				}
			}
			return
		}, &prog); e != nil {
		err = e
	} else if e := assembleAliases(asm, aliases); e != nil {
		err = e
	} else if e := assembleDirectives(asm, names, directives); e != nil {
		err = e
	}
	return
}

func assembleDirectives(asm *Assembler, names []string, directives []*grammar.Directive) (err error) {
	// write directives
	for i, d := range directives {
		if e := asm.WriteProgram(names[i], "Directive", d); e != nil {
			err = e
			break
		}
	}
	return
}

func assembleAliases(asm *Assembler, aliases []*grammar.Alias) (err error) {
	for _, d := range aliases {
		shortName := d.AsNoun
		for _, alias := range d.Names {
			var fullName string
			if e := asm.cache.QueryRow(
				`select noun
				from mdl_name
				join mdl_noun
					using (noun)
				where UPPER(name)=UPPER(?)
				order by rank
				limit 1`, shortName).Scan(&fullName); e != nil {
				err = e
			} else {
				asm.WriteName(fullName, alias, -1)
			}
		}
	}

	return
}
