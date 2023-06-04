package weave

import (
	"github.com/ionous/errutil"
)

type ScopedKind struct {
	Requires // references to ancestors ( at most it can have one direct parent )
	domain   *Domain
}

func (k *ScopedKind) Resolve() (ret Dependencies, err error) {
	if len(k.at) == 0 {
		err = errutil.Fmt("kind %q never defined", k.name)
	} else if ks, e := k.resolve(k, (*kindFinder)(k.domain)); e != nil {
		err = errutil.Fmt("%s for kind %q ", e, k.name)
	} else {
		ret = ks
	}
	return
}

// private helper to make the catalog compatible with the DependencyFinder ( for domains )
type kindFinder Domain

// look upwards through the domains to find the named kind
func (kf *kindFinder) FindDependency(name string) (ret Dependency, okay bool) {
	d := (*Domain)(kf)

	// tries to use the passed name as a plural, if that fails, tries as a singular word
	if a, ok := d.getKind(name); ok {
		ret, okay = a, true
	} else if p := d.catalog.run.PluralOf(name); p != name {
		ret, okay = d.getKind(p)
	}
	return
}
