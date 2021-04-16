package assembly

import (
	"strconv"

	"git.sr.ht/~ionous/iffy/dl/grammar"
	"git.sr.ht/~ionous/iffy/tables"
)

func AssembleGrammar(asm *Assembler) (err error) {
	var grams []grammar.Grammar
	var names []string
	var prog []byte
	if e := tables.QueryAll(asm.cache.DB(),
		`select prog 
		from eph_prog where progType='grammar'
		order by rowid`,
		func() (err error) {
			var gram grammar.Grammar
			var name string // try to make a name for ourselves
			if e := tables.DecodeGob(prog, &gram); e != nil {
				err = e
			} else {
				// if we store the scanner, then we have to register the scanners
				// which means we have both the parser and the maker operations
				// in gob -- but only use the maker operations here.
				// for now, just let runtime make the scanners.. (too).
				gram.Scanner.MakeScanner()
				if allOf, ok := gram.Scanner.(*grammar.AllOf); ok {
					if len(allOf.Series) > 0 {
						if words, ok := allOf.Series[0].(*grammar.Words); ok {
							name = words.Words
						}
					}
				}
				if len(name) == 0 {
					name = "scanner" + strconv.Itoa(len(grams)+1)
				}
				grams = append(grams, gram)
				names = append(names, name)
			}
			return
		}, &prog); e != nil {
		err = e
	} else {
		for i, gram := range grams {
			if e := asm.WriteGob(names[i], &gram); e != nil {
				err = e
				break
			}
		}
	}
	return
}
