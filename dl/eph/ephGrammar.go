package eph

import (
	"github.com/ionous/errutil"
)

func (el *EphGrammarAlias) Phase() Phase { return GrammarPhase }

func (el *EphGrammarAlias) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if short, ok := UniformString(el.ShortName); !ok {
		err = InvalidString(el.ShortName)
	} else if noun, ok := d.GetClosestNoun(short); !ok {
		err = errutil.New("unknown noun", el.ShortName)
	} else {
		for _, a := range el.Aliases {
			if a, ok := UniformString(a); !ok {
				err = errutil.Append(err, InvalidString(a))
			} else {
				noun.aliases = append(noun.aliases, a)
			}
		}
	}
	return
}

// func (el *EphGrammarDirective) Phase() Phase { return GrammarPhase }

// func (el *EphGrammarDirective) Assemble(c *Catalog, d *Domain, at string) (err error) {
// 	var prog grammar.Directive
// 	if e := decodeGrammar(&prog, el.Prog); el != nil {
// 		err = e
// 	} else {
// 		// for now, just let runtime make the scanners..
// 		// because, if we store the scanner, then we have to register the scanners.
// 		// which means we have both the parser and the maker operations registered
// 		// while using only one of them
// 		key := strings.Join(prog.Lede, "/")
// 		if d.AddDefinition(key, at, "Directive:"+prog); e != nil {
// 			err = e
// 		}
// 	}
// 	return
// }

// func decodeGrammar(op grammar.GrammarMaker, prog string) error {
// 	return cin.Decode(op, prog, iffy.AllSignatures)
// }
