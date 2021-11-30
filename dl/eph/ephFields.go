package eph

import (
	"strings"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"github.com/ionous/errutil"
)

func (el *EphFields) Phase() Phase { return FieldPhase }

// add some fields to a kind.
// see also: EphAspects which generates traits and adds them to a custom aspect kind.
func (el *EphFields) Assemble(c *Catalog, d *Domain, at string) (err error) {
	// hooray parameter validation
	// note: the kinds must exist ( and are resolved if they do ) already ( re: phased processing )
	if singleKind, e := c.Singularize(d.name, strings.TrimSpace(el.Kinds)); e != nil {
		err = e
	} else if newKind, ok := UniformString(singleKind); !ok {
		err = InvalidString(el.Kinds)
	} else if kind, ok := d.GetKind(newKind); !ok {
		err = errutil.New("unknown kind", newKind)
	} else if name, ok := UniformString(el.Name); !ok {
		err = InvalidString(el.Name)
	} else if aff, ok := composer.FindChoice(&el.Affinity, el.Affinity.Str); !ok && len(el.Affinity.Str) > 0 {
		err = errutil.New("unknown affinity", aff)
	} else if class, ok := UniformString(el.Class); !ok && len(el.Class) > 0 {
		err = InvalidString(el.Class)
	} else if _, ok := d.GetKind(class); !ok && len(class) > 0 {
		err = errutil.New("unknown field class", class)
	} else {
		// checks for conflicts, allows duplicates.
		err = kind.AddFields(&fieldDef{
			name: name, affinity: aff, class: class, at: at,
		})
	}
	return
}
