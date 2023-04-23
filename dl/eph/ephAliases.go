package eph

import (
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"github.com/ionous/errutil"
)

// fix? right now we're adding aliases in grammar --
// so its only going to affect player input
// but... we could add this after noun declaration to allow fields, etc. to use aliased names
func (op *EphAliases) Phase() assert.Phase { return assert.AliasPhase }

func (op *EphAliases) Weave(k assert.Assertions) (err error) {
	return k.AssertAlias(op.ShortName, op.Aliases...)
}

func (op *EphAliases) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if noun, e := getClosestNoun(d, op.ShortName); e != nil {
		err = e
	} else {
		for _, a := range op.Aliases {
			if a, ok := UniformString(a); !ok {
				err = errutil.Append(err, InvalidString(a))
			} else {
				if !noun.AddAlias(a, at) {
					LogWarning(errutil.Fmt("duplicate alias %q for %q at %s", a, noun.name, at))
				}
			}
		}
	}
	return
}

func getClosestNoun(d *Domain, rawName string) (ret *ScopedNoun, err error) {
	if short, ok := UniformString(rawName); !ok {
		err = InvalidString(rawName)
	} else if noun, ok := d.GetClosestNoun(short); !ok {
		err = errutil.New("unknown noun", rawName)
	} else {
		ret = noun
	}
	return
}
